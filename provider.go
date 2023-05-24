package kms_wallet_provider

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ether_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/patrickmn/go-cache"
	"math/big"
	"time"
)

const (
	publicKeyCacheKey = "kms-public-key:%s"
)

var (
	secp256k1N           = crypto.S256().Params().N
	secp256k1HalfN       = new(big.Int).Div(secp256k1N, big.NewInt(2))
	walletAddressTagKey  = "walletAddress"
	defaultCacheDuration = time.Hour * 24 * 365
)

type asn1EcPublicKey struct {
	EcPublicKeyInfo asn1EcPublicKeyInfo
	PublicKey       asn1.BitString
}

type asn1EcPublicKeyInfo struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

type asn1EcSig struct {
	R asn1.RawValue
	S asn1.RawValue
}

type KMSClient interface {
	CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
	CreateAlias(ctx context.Context, params *kms.CreateAliasInput, optFns ...func(*kms.Options)) (*kms.CreateAliasOutput, error)
	TagResource(ctx context.Context, params *kms.TagResourceInput, optFns ...func(*kms.Options)) (*kms.TagResourceOutput, error)
	DescribeKey(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error)
	GetPublicKey(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error)
	Sign(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error)
	EnableKey(ctx context.Context, params *kms.EnableKeyInput, optFns ...func(*kms.Options)) (*kms.EnableKeyOutput, error)
	DisableKey(ctx context.Context, params *kms.DisableKeyInput, optFns ...func(*kms.Options)) (*kms.DisableKeyOutput, error)
}

type KMSWallet struct {
	Address string
	KeyId   string
}

type CreateWalletInput struct {
	Alias                           *string
	IgnoreDefaultWalletAddressAlias bool
	AddWalletAddressTag             bool
	BypassPolicyLockoutSafetyCheck  bool
	CustomKeyStoreId                *string
	Description                     *string
	MultiRegion                     *bool
	Origin                          types.OriginType
	Policy                          *string
	Tags                            map[string]string
	XksKeyId                        *string
}

type Provider interface {
	CreateWallet(ctx context.Context, input CreateWalletInput) (wallet KMSWallet, err error)
	GetWallet(ctx context.Context, keyId string) (wallet KMSWallet, err error)
	GetWalletTransactor(ctx context.Context, keyId string, chainId *big.Int) (*bind.TransactOpts, error)
	GetWalletCaller(ctx context.Context, keyId string, chainId *big.Int) (*bind.CallOpts, error)
	SignMessage(ctx context.Context, keyId string, message []byte) ([]byte, error)
	EnableWallet(ctx context.Context, keyId string) (*kms.EnableKeyOutput, error)
	DisableWallet(ctx context.Context, keyId string) (*kms.DisableKeyOutput, error)

	GetWalletByAlias(ctx context.Context, alias string) (wallet KMSWallet, err error)
	GetWalletTransactorByAlias(ctx context.Context, alias string, chainId *big.Int) (*bind.TransactOpts, error)
	GetWalletCallerByAlias(ctx context.Context, alias string, chainId *big.Int) (*bind.CallOpts, error)
	SignMessageByAlias(ctx context.Context, alias string, message []byte) ([]byte, error)
	EnableWalletByAlias(ctx context.Context, alias string) (*kms.EnableKeyOutput, error)
	DisableWalletByAlias(ctx context.Context, alias string) (*kms.DisableKeyOutput, error)
	GetKeyIdByAlias(ctx context.Context, alias string) (keyId string, err error)
}
type provider struct {
	client KMSClient
	cache  *cache.Cache
}

func New(client KMSClient, cacheExpiration *time.Duration) Provider {
	if cacheExpiration == nil {
		cacheExpiration = &defaultCacheDuration
	}

	return &provider{
		client: client,
		cache:  cache.New(*cacheExpiration, time.Hour),
	}
}

func (c *provider) CreateWallet(ctx context.Context, input CreateWalletInput) (wallet KMSWallet, err error) {
	var tags []types.Tag
	for key, value := range input.Tags {
		tagKey := key
		tagValue := value
		tags = append(tags, types.Tag{
			TagKey:   &tagKey,
			TagValue: &tagValue,
		})
	}

	output, err := c.client.CreateKey(ctx, &kms.CreateKeyInput{
		BypassPolicyLockoutSafetyCheck: input.BypassPolicyLockoutSafetyCheck,
		CustomKeyStoreId:               input.CustomKeyStoreId,
		Description:                    input.Description,
		KeySpec:                        types.KeySpecEccSecgP256k1,
		KeyUsage:                       types.KeyUsageTypeSignVerify,
		MultiRegion:                    input.MultiRegion,
		Origin:                         input.Origin,
		Policy:                         input.Policy,
		Tags:                           tags,
		XksKeyId:                       input.XksKeyId,
	})

	if err != nil {
		return wallet, err
	}

	wallet, err = c.GetWallet(ctx, *output.KeyMetadata.KeyId)
	if err != nil {
		return wallet, err
	}

	alias := input.Alias
	if alias == nil && !input.IgnoreDefaultWalletAddressAlias {
		alias = &wallet.Address
	}

	if alias != nil {
		prefixedAlias := getPrefixedAlias(*alias)
		_, err = c.client.CreateAlias(ctx, &kms.CreateAliasInput{
			AliasName:   &prefixedAlias,
			TargetKeyId: output.KeyMetadata.KeyId,
		})

		if err != nil {
			return wallet, err
		}
	}

	if input.AddWalletAddressTag {
		_, err = c.client.TagResource(ctx, &kms.TagResourceInput{
			KeyId: output.KeyMetadata.KeyId,
			Tags: []types.Tag{
				{
					TagKey:   &walletAddressTagKey,
					TagValue: &wallet.Address,
				},
			},
		})

		if err != nil {
			return wallet, err
		}
	}

	return wallet, err
}

func (c *provider) GetWallet(ctx context.Context, keyId string) (wallet KMSWallet, err error) {
	publicKey, err := c.getPublicKey(ctx, keyId)
	if err != nil {
		return wallet, err
	}

	publicKeyAddress := crypto.PubkeyToAddress(*publicKey)
	return KMSWallet{
		Address: publicKeyAddress.String(),
		KeyId:   keyId,
	}, err
}

func (c *provider) GetWalletByAlias(ctx context.Context, alias string) (wallet KMSWallet, err error) {
	keyId, err := c.GetKeyIdByAlias(ctx, alias)
	if err != nil {
		return wallet, err
	}

	return c.GetWallet(ctx, keyId)
}

func (c *provider) DisableWallet(ctx context.Context, keyId string) (*kms.DisableKeyOutput, error) {
	return c.client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: &keyId})
}

func (c *provider) DisableWalletByAlias(ctx context.Context, alias string) (*kms.DisableKeyOutput, error) {
	keyId, err := c.GetKeyIdByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	return c.client.DisableKey(ctx, &kms.DisableKeyInput{KeyId: &keyId})
}

func (c *provider) EnableWallet(ctx context.Context, keyId string) (*kms.EnableKeyOutput, error) {
	return c.client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: &keyId})
}

func (c *provider) EnableWalletByAlias(ctx context.Context, alias string) (*kms.EnableKeyOutput, error) {
	keyId, err := c.GetKeyIdByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	return c.client.EnableKey(ctx, &kms.EnableKeyInput{KeyId: &keyId})
}

func (c *provider) GetWalletTransactor(ctx context.Context, keyId string, chainId *big.Int) (*bind.TransactOpts, error) {
	publicKey, err := c.getPublicKey(ctx, keyId)
	if err != nil {
		return nil, err
	}

	publicKeyBytes := secp256k1.S256().Marshal(publicKey.X, publicKey.Y)
	publicKeyAddress := crypto.PubkeyToAddress(*publicKey)
	if chainId == nil {
		return nil, bind.ErrNoChainID
	}

	signer := ether_types.LatestSignerForChainID(chainId)
	signerFn := func(address common.Address, tx *ether_types.Transaction) (*ether_types.Transaction, error) {
		if address != publicKeyAddress {
			return nil, bind.ErrNotAuthorized
		}

		txHashBytes := signer.Hash(tx).Bytes()

		rBytes, sBytes, err := c.getSignatureFromKms(ctx, keyId, txHashBytes)
		if err != nil {
			return nil, err
		}

		// Adjust S value from signature according to Ethereum standard
		sBigInt := new(big.Int).SetBytes(sBytes)
		if sBigInt.Cmp(secp256k1HalfN) > 0 {
			sBytes = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
		}

		signature, err := c.getEthereumSignature(publicKeyBytes, txHashBytes, rBytes, sBytes)
		if err != nil {
			return nil, err
		}

		return tx.WithSignature(signer, signature)
	}

	return &bind.TransactOpts{
		From:   publicKeyAddress,
		Signer: signerFn,
	}, nil
}

func (c *provider) GetWalletTransactorByAlias(ctx context.Context, alias string, chainId *big.Int) (*bind.TransactOpts, error) {
	keyId, err := c.GetKeyIdByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	return c.GetWalletTransactor(ctx, keyId, chainId)
}

func (c *provider) GetWalletCaller(ctx context.Context, keyId string, chainId *big.Int) (*bind.CallOpts, error) {
	publicKey, err := c.getPublicKey(ctx, keyId)
	if err != nil {
		return nil, err
	}

	publicKeyAddress := crypto.PubkeyToAddress(*publicKey)
	if chainId == nil {
		return nil, bind.ErrNoChainID
	}

	return &bind.CallOpts{
		From: publicKeyAddress,
	}, nil
}

func (c *provider) GetWalletCallerByAlias(ctx context.Context, alias string, chainId *big.Int) (*bind.CallOpts, error) {
	keyId, err := c.GetKeyIdByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	return c.GetWalletCaller(ctx, keyId, chainId)
}

func (c *provider) SignMessage(ctx context.Context, keyId string, message []byte) ([]byte, error) {
	hashedMessage := toEthSignedMessageHash(message)

	rBytes, sBytes, err := c.getSignatureFromKms(ctx, keyId, hashedMessage)
	if err != nil {
		return nil, err
	}

	sBigInt := new(big.Int).SetBytes(sBytes)
	if sBigInt.Cmp(secp256k1HalfN) > 0 {
		sBytes = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
	}

	publicKey, err := c.getPublicKey(ctx, keyId)
	if err != nil {
		return nil, err
	}

	publicKeyBytes := secp256k1.S256().Marshal(publicKey.X, publicKey.Y)
	signature, err := c.getEthereumSignature(publicKeyBytes, hashedMessage, rBytes, sBytes)
	if err != nil {
		return nil, err
	}

	signature[64] += 27
	return signature, nil
}

func (c *provider) SignMessageByAlias(ctx context.Context, alias string, message []byte) ([]byte, error) {
	keyId, err := c.GetKeyIdByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	return c.SignMessage(ctx, keyId, message)
}

func (c *provider) GetKeyIdByAlias(ctx context.Context, alias string) (keyId string, err error) {
	prefixedAlias := getPrefixedAlias(alias)
	output, err := c.client.DescribeKey(ctx, &kms.DescribeKeyInput{
		KeyId: &prefixedAlias,
	})

	if err != nil {
		return keyId, fmt.Errorf("can not get public key from KMS for alias: %s, err: %+v", alias, err)
	}

	return *output.KeyMetadata.KeyId, err
}

func (c *provider) getSignatureFromKms(
	ctx context.Context, keyId string, txHashBytes []byte,
) ([]byte, []byte, error) {
	signInput := &kms.SignInput{
		KeyId:            aws.String(keyId),
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
		MessageType:      types.MessageTypeDigest,
		Message:          txHashBytes,
	}

	signOutput, err := c.client.Sign(ctx, signInput)
	if err != nil {
		return nil, nil, err
	}

	var sigAsn1 asn1EcSig
	_, err = asn1.Unmarshal(signOutput.Signature, &sigAsn1)
	if err != nil {
		return nil, nil, err
	}

	return sigAsn1.R.Bytes, sigAsn1.S.Bytes, nil
}

func (c *provider) getEthereumSignature(expectedPublicKeyBytes []byte, txHash []byte, r []byte, s []byte) ([]byte, error) {
	rsSignature := append(adjustSignatureLength(r), adjustSignatureLength(s)...)
	signature := append(rsSignature, []byte{0}...)

	recoveredPublicKeyBytes, err := crypto.Ecrecover(txHash, signature)
	if err != nil {
		return nil, err
	}

	if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
		signature = append(rsSignature, []byte{1}...)
		recoveredPublicKeyBytes, err = crypto.Ecrecover(txHash, signature)
		if err != nil {
			return nil, err
		}

		if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
			return nil, errors.New("can not reconstruct public key from sig")
		}
	}

	return signature, nil
}

func (c *provider) getPublicKey(ctx context.Context, keyId string) (*ecdsa.PublicKey, error) {
	cacheKey := fmt.Sprintf(publicKeyCacheKey, keyId)
	foundPublicKey, found := c.cache.Get(cacheKey)
	if found {
		publicKey := foundPublicKey.(ecdsa.PublicKey)
		return &publicKey, nil
	}

	publicKeyBytes, err := c.getPublicKeyBytes(ctx, keyId)
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("can not construct secp256k1 public key from key bytes for keyId: %s, err: %+v", keyId, err)
	}

	c.cache.Set(cacheKey, *publicKey, 24*30*time.Hour)
	return publicKey, err
}

func (c *provider) getPublicKeyBytes(ctx context.Context, keyId string) ([]byte, error) {
	getPubKeyOutput, err := c.client.GetPublicKey(ctx, &kms.GetPublicKeyInput{
		KeyId: aws.String(keyId),
	})

	if err != nil {
		return nil, fmt.Errorf("can not get public key from KMS for keyId: %s, err: %+v", keyId, err)
	}

	var asn1pubk asn1EcPublicKey
	_, err = asn1.Unmarshal(getPubKeyOutput.PublicKey, &asn1pubk)
	if err != nil {
		return nil, fmt.Errorf("can not parse asn1 public key for keyId: %s, err: %+v", keyId, err)
	}

	return asn1pubk.PublicKey.Bytes, nil
}

func getPrefixedAlias(alias string) string {
	return fmt.Sprintf("alias/%s", alias)
}

func toEthSignedMessageHash(hash []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256Hash([]byte(msg)).Bytes()
}

func adjustSignatureLength(buffer []byte) []byte {
	buffer = bytes.TrimLeft(buffer, "\x00")
	for len(buffer) < 32 {
		zeroBuf := []byte{0}
		buffer = append(zeroBuf, buffer...)
	}
	return buffer
}

package kms_wallet_provider_test

import (
	"context"
	"encoding/base64"
	"github.com/aliarbak/go-ethereum-aws-kms-wallet-provider"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/big"
	"testing"
)

type mockKMSClient struct {
	mock.Mock
}

func (m *mockKMSClient) CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.CreateKeyOutput), args.Error(1)
}

func (m *mockKMSClient) CreateAlias(ctx context.Context, params *kms.CreateAliasInput, optFns ...func(*kms.Options)) (*kms.CreateAliasOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.CreateAliasOutput), args.Error(1)
}

func (m *mockKMSClient) TagResource(ctx context.Context, params *kms.TagResourceInput, optFns ...func(*kms.Options)) (*kms.TagResourceOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.TagResourceOutput), args.Error(1)
}

func (m *mockKMSClient) DescribeKey(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.DescribeKeyOutput), args.Error(1)
}

func (m *mockKMSClient) Sign(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.SignOutput), args.Error(1)
}

func (m *mockKMSClient) EnableKey(ctx context.Context, params *kms.EnableKeyInput, optFns ...func(*kms.Options)) (*kms.EnableKeyOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.EnableKeyOutput), args.Error(1)
}

func (m *mockKMSClient) DisableKey(ctx context.Context, params *kms.DisableKeyInput, optFns ...func(*kms.Options)) (*kms.DisableKeyOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.DisableKeyOutput), args.Error(1)
}

func (m *mockKMSClient) GetPublicKey(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*kms.GetPublicKeyOutput), args.Error(1)
}

func TestCreateWallet_Should_Create_Wallet_With_Wallet_Address_Tag_When_Add_Wallet_Address_Tag_Is_True(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	input := kms_wallet_provider.CreateWalletInput{
		Tags:                            map[string]string{"tag1": "value1"},
		AddWalletAddressTag:             true,
		IgnoreDefaultWalletAddressAlias: true,
	}

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   "keyId",
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}

	mockClient.On("CreateKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.CreateKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: aws.String("keyId"),
		},
	}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	mockClient.On("TagResource", mock.Anything, mock.Anything, mock.Anything).Return(&kms.TagResourceOutput{}, nil)

	// when
	output, err := provider.CreateWallet(context.Background(), input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	mockClient.AssertNumberOfCalls(t, "CreateKey", 1)
	mockClient.AssertNumberOfCalls(t, "CreateAlias", 0)
	mockClient.AssertNumberOfCalls(t, "TagResource", 1)
}

func TestCreateWallet_Should_Create_Wallet_Without_Wallet_Address_Tag_When_Add_Wallet_Address_Tag_Is_False(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	input := kms_wallet_provider.CreateWalletInput{
		Tags:                            map[string]string{"tag1": "value1"},
		AddWalletAddressTag:             false,
		IgnoreDefaultWalletAddressAlias: true,
	}

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   "keyId",
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}

	mockClient.On("CreateKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.CreateKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: aws.String("keyId"),
		},
	}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.CreateWallet(context.Background(), input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	mockClient.AssertNumberOfCalls(t, "CreateKey", 1)
	mockClient.AssertNumberOfCalls(t, "CreateAlias", 0)
	mockClient.AssertNumberOfCalls(t, "TagResource", 0)
}

func TestCreateWallet_Should_Create_Wallet_With_Alias_When_Alias_Is_Provided(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	alias := "alias"
	input := kms_wallet_provider.CreateWalletInput{
		Alias:                           &alias,
		Tags:                            map[string]string{"tag1": "value1"},
		AddWalletAddressTag:             false,
		IgnoreDefaultWalletAddressAlias: true,
	}

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   "keyId",
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}

	mockClient.On("CreateKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.CreateKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: aws.String("keyId"),
		},
	}, nil)

	mockClient.On("CreateAlias", mock.Anything, mock.Anything, mock.Anything).Return(&kms.CreateAliasOutput{}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.CreateWallet(context.Background(), input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	mockClient.AssertNumberOfCalls(t, "CreateKey", 1)
	mockClient.AssertNumberOfCalls(t, "CreateAlias", 1)
	mockClient.AssertNumberOfCalls(t, "TagResource", 0)
}

func TestCreateWallet_Should_Create_Wallet_With_Alias_When_Alias_Is_Not_Provided_But_Ignore_Default_Wallet_address_Alias_Is_False(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	input := kms_wallet_provider.CreateWalletInput{
		Tags:                            map[string]string{"tag1": "value1"},
		AddWalletAddressTag:             false,
		IgnoreDefaultWalletAddressAlias: false,
	}

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   "keyId",
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}

	mockClient.On("CreateKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.CreateKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: aws.String("keyId"),
		},
	}, nil)

	mockClient.On("CreateAlias", mock.Anything, mock.Anything, mock.Anything).Return(&kms.CreateAliasOutput{}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.CreateWallet(context.Background(), input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	mockClient.AssertNumberOfCalls(t, "CreateKey", 1)
	mockClient.AssertNumberOfCalls(t, "CreateAlias", 1)
	mockClient.AssertNumberOfCalls(t, "TagResource", 0)
}

func TestGetWallet(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   "keyId",
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}
	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.GetWallet(context.Background(), "keyId")

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetWallet_When_Public_Key_Already_Cached(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   "keyId",
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}
	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	firstOutput, err := provider.GetWallet(context.Background(), "keyId")
	secondOutput, err := provider.GetWallet(context.Background(), "keyId")

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, firstOutput)
	assert.Equal(t, firstOutput, secondOutput)
	mockClient.AssertNumberOfCalls(t, "GetPublicKey", 1)
}

func TestGetWalletByAlias(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	keyId := "keyId"
	expectedOutput := kms_wallet_provider.KMSWallet{
		KeyId:   keyId,
		Address: "0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118",
	}

	mockClient.On("DescribeKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DescribeKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: &keyId,
		},
	}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.GetWalletByAlias(context.Background(), "alias")

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetWalletTransactor(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	keyId := "keyId"

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.GetWalletTransactor(context.Background(), keyId, big.NewInt(1))

	// then
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118"), output.From)
}

func TestGetWalletTransactorByAlias(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	keyId := "keyId"

	mockClient.On("DescribeKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DescribeKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: &keyId,
		},
	}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.GetWalletTransactorByAlias(context.Background(), "alias", big.NewInt(1))

	// then
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118"), output.From)
}

func TestGetWalletCaller(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	keyId := "keyId"

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	// when
	output, err := provider.GetWalletCaller(context.Background(), keyId, big.NewInt(1))

	// then
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118"), output.From)
}

func TestGetWalletCallerByAlias(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)

	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	keyId := "keyId"

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	mockClient.On("DescribeKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DescribeKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: &keyId,
		},
	}, nil)

	// when
	output, err := provider.GetWalletCallerByAlias(context.Background(), "alias", big.NewInt(1))

	// then
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x5B1a501FAB5c6D78CBd61F31f3B4B42286Bcf118"), output.From)
}

func TestDisableWallet(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	keyId := "keyId"

	mockClient.On("DisableKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DisableKeyOutput{}, nil)

	// when
	_, err := provider.DisableWallet(context.Background(), keyId)

	// then
	assert.NoError(t, err)
	mockClient.AssertNumberOfCalls(t, "DisableKey", 1)
}

func TestDisableWalletByAlias(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	keyId := "keyId"

	mockClient.On("DescribeKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DescribeKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: &keyId,
		},
	}, nil)

	mockClient.On("DisableKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DisableKeyOutput{}, nil)

	// when
	_, err := provider.DisableWalletByAlias(context.Background(), "alias")

	// then
	assert.NoError(t, err)
	mockClient.AssertNumberOfCalls(t, "DisableKey", 1)
}

func TestEnableWallet(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	keyId := "keyId"

	mockClient.On("EnableKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.EnableKeyOutput{}, nil)

	// when
	_, err := provider.EnableWallet(context.Background(), keyId)

	// then
	assert.NoError(t, err)
	mockClient.AssertNumberOfCalls(t, "EnableKey", 1)
}

func TestEnableWalletByAlias(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	keyId := "keyId"

	mockClient.On("DescribeKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DescribeKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: &keyId,
		},
	}, nil)

	mockClient.On("EnableKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.EnableKeyOutput{}, nil)

	// when
	_, err := provider.EnableWalletByAlias(context.Background(), "alias")

	// then
	assert.NoError(t, err)
	mockClient.AssertNumberOfCalls(t, "EnableKey", 1)
}

func TestSignMessage_When_Successful(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE1TWYYp+bySM6J3m99CxKhHZFgJqdwm5V6kGziuJ4kc7f0xORg5AqGRKbFNTSsmYTXNi2Z/cl298eyKTbmy+8aQ==")
	message := []byte("Hello World!")
	keyId := "keyId"

	signOutput, _ := base64.StdEncoding.DecodeString("MEUCIQD6wEcOzyjl8wr+OR8In54bVKgR5/ZogQQWiHPkqb70NQIgZZnGhIRlfGj9xexxACWQ4WZPB60swZP5DK6OjyEJIIc=")
	expectedOutput, _ := base64.StdEncoding.DecodeString("+sBHDs8o5fMK/jkfCJ+eG1SoEef2aIEEFohz5Km+9DVlmcaEhGV8aP3F7HEAJZDhZk8HrSzBk/kMro6PIQkghxw=")

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	mockClient.On("Sign", mock.Anything, mock.Anything, mock.Anything).Return(&kms.SignOutput{
		Signature: signOutput,
	}, nil)

	// when
	output, err := provider.SignMessage(context.Background(), keyId, message)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
	mockClient.AssertNumberOfCalls(t, "Sign", 1)
}

func TestSignMessage_When_Recover_Failed(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAERtrxsFyn7UzP2OgzzJA6Y89p/2175fOwXeP33ACZgmdD2jJlQdypNM9CCDm3J6uqTrvYrO0hwF8p/k/Tf94DjA==")
	message := []byte("Hello World!")
	keyId := "keyId"

	signOutput, _ := base64.StdEncoding.DecodeString("MEUCIQD6wEcOzyjl8wr+OR8In54bVKgR5/ZogQQWiHPkqb70NQIgZZnGhIRlfGj9xexxACWQ4WZPB60swZP5DK6OjyEJIIc=")

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	mockClient.On("Sign", mock.Anything, mock.Anything, mock.Anything).Return(&kms.SignOutput{
		Signature: signOutput,
	}, nil)

	// when
	_, err := provider.SignMessage(context.Background(), keyId, message)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "can not reconstruct public key from sig")
	mockClient.AssertNumberOfCalls(t, "Sign", 1)
}

func TestSignMessageByAlias(t *testing.T) {
	// given
	mockClient := &mockKMSClient{}
	provider := kms_wallet_provider.New(mockClient, nil)
	publicKey, _ := base64.StdEncoding.DecodeString("MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE1TWYYp+bySM6J3m99CxKhHZFgJqdwm5V6kGziuJ4kc7f0xORg5AqGRKbFNTSsmYTXNi2Z/cl298eyKTbmy+8aQ==")
	message := []byte("Hello World!")
	keyId := "keyId"

	signOutput, _ := base64.StdEncoding.DecodeString("MEUCIQD6wEcOzyjl8wr+OR8In54bVKgR5/ZogQQWiHPkqb70NQIgZZnGhIRlfGj9xexxACWQ4WZPB60swZP5DK6OjyEJIIc=")
	expectedOutput, _ := base64.StdEncoding.DecodeString("+sBHDs8o5fMK/jkfCJ+eG1SoEef2aIEEFohz5Km+9DVlmcaEhGV8aP3F7HEAJZDhZk8HrSzBk/kMro6PIQkghxw=")

	mockClient.On("DescribeKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.DescribeKeyOutput{
		KeyMetadata: &types.KeyMetadata{
			KeyId: &keyId,
		},
	}, nil)

	mockClient.On("GetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(&kms.GetPublicKeyOutput{
		PublicKey: publicKey,
	}, nil)

	mockClient.On("Sign", mock.Anything, mock.Anything, mock.Anything).Return(&kms.SignOutput{
		Signature: signOutput,
	}, nil)

	// when
	output, err := provider.SignMessageByAlias(context.Background(), "alias", message)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
	mockClient.AssertNumberOfCalls(t, "Sign", 1)
}

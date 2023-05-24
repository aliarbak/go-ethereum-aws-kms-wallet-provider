# AWS KMS Ethereum Wallet Provider for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/aliarbak/go-ethereum-aws-kms-wallet-provider.svg)](https://pkg.go.dev/github.com/aliarbak/go-ethereum-aws-kms-wallet-provider) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

The `kms_wallet_provider` is a Go package that enables the creation of Ethereum wallets using the AWS Key Management Service (KMS). It allows you to create  wallets(keys on KMS) and sign transactions or messages with them.

The transaction signing implementations in this package are derived from the [go-ethereum-aws-kms-tx-signer](https://github.com/welthee/go-ethereum-aws-kms-tx-signer), which is licensed under the MIT License.

### Table of Contents
- [Installation](#installation)
- [Functionality and Usage](#functionality-and-usage)
	- [CreateWallet](#createwallet)
	- [GetWallet](#getwallet)
	- [GetWalletTransactor](#getwallettransactor)
	- [GetWalletCaller](#getwalletcaller)
	- [SignMessage](#signmessage)
	- [EnableWallet](#enablewallet)
	- [DisableWallet](#disablewallet)
	- [Additional Functions](#additional-functions)
- [Example Usage](#example-usage)


## Installation

You can install it using the following command:

```bash
go get github.com/aliarbak/go-ethereum-aws-kms-wallet-provider
```

Once installed, you can import the package in your Go code:

```go
import "github.com/aliarbak/go-ethereum-aws-kms-wallet-provider/kms_wallet_provider"
```

To create a provider, call the `kms_wallet_provider.New(client *kms.Client, cacheExpiration *time.Duration) Provider` function. It requires the following parameters:

- `client`: A reference to the `kms.Client` for AWS KMS.
- `cacheExpiration`: The cache expiration duration for public keys to avoid fetching them from KMS every time. If `nil` is provided, the default duration of 1 year will be used.

To create a new kms.Client:
```go
config := aws.Config{
    Region:      "eu-central-1",
    Credentials: credentials.NewStaticCredentialsProvider("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", ""),
}

kmsClient := kms.NewFromConfig(config) // or you can use kms.New(...)
```

## Functionality and Usage

The `kms_wallet_provider` package provides the following functions:

### CreateWallet

```go
func CreateWallet(ctx context.Context, input CreateWalletInput) (wallet KMSWallet, err error)
```

The `CreateWallet` function is used to create a new wallet with the specified input parameters. The `CreateWalletInput` struct is defined as follows:

```go
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
```

- `Alias`: Specifies a custom alias for the key (e.g., userId).
- `IgnoreDefaultWalletAddressAlias`: If the `Alias` value is `nil`, the generated wallet address is assigned as the alias. Set this value to `true` if you want to prevent this and add an alias to the key.
- `AddWalletAddressTag`: If set to `true`, the generated wallet address is added as a tag (`walletAddress`) to the key.

### GetWallet

```go
func GetWallet(ctx context.Context, keyId string) (wallet KMSWallet, err error)
```

The `GetWallet` function retrieves a wallet by the specified `keyId`.

### GetWalletTransactor

```go
func GetWalletTransactor(ctx context.Context, keyId string, chainId *big.Int) (*bind.TransactOpts, error)
```

The `GetWalletTransactor` function returns a transaction signer (`bind.TransactOpts`) for the wallet associated with the given `keyId` and `chainId`.

### GetWalletCaller

```go
func GetWalletCaller(ctx context.Context, keyId string, chainId *big.Int) (*bind.CallOpts, error)
```

The `GetWalletCaller` function returns a contract caller (`bind.CallOpts`) for the wallet associated with the given `keyId` and `chainId`.

### SignMessage

```go
func SignMessage(ctx context.Context, keyId string, message []byte) ([]byte, error)
```

The `SignMessage` function signs the specified `message` using the wallet associated with the given `keyId` and returns the signature.

### EnableWallet

```go
func EnableWallet(ctx context.Context, keyId string) (*kms.EnableKeyOutput, error)
```

The `EnableWallet` function enables the wallet associated with the given `keyId`.

### DisableWallet

```go
func DisableWallet(ctx context.Context, keyId string) (*kms.DisableKeyOutput, error)
```

The `DisableWallet` function disables the wallet associated with the given `keyId`.

### Additional Functions

The package also provides several utility functions to work with aliases:

- `GetWalletByAlias`: Retrieves a wallet by the specified `alias`.
- `GetWalletTransactorByAlias`: Returns a transaction signer for the wallet associated with the given `alias` and `chainId`.
- `GetWalletCallerByAlias`: Returns a contract caller for the wallet associated with the given `alias` and `chainId`.
- `SignMessageByAlias`: Signs the specified `message` using the wallet associated with the given `alias` and returns the signature.
- `EnableWalletByAlias`: Enables the wallet associated with the given `alias`.
- `DisableWalletByAlias`: Disables the wallet associated with the given `alias`.
- `GetKeyIdByAlias`: Retrieves the keyId associated with the given `alias`.

## Example Usage
You can access detailed usage example [from this link](https://github.com/aliarbak/go-ethereum-aws-kms-wallet-provider/example/readme.md).
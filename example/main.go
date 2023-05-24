package main

import (
	"context"
	kmswallet "github.com/aliarbak/go-ethereum-aws-kms-wallet-provider"
	contracts "github.com/aliarbak/go-ethereum-aws-kms-wallet-provider/example/contracts"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	ctx := context.Background()
	chainId := big.NewInt(80001) // means mumbai testnet

	// Initialize the KMS Client
	config := aws.Config{
		Region:      "eu-central-1",
		Credentials: credentials.NewStaticCredentialsProvider("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", ""),
	}

	kmsClient := kms.NewFromConfig(config) // or you can use kms.New(...)

	// Initialize the wallet provider (with default cache duration)
	walletProvider := kmswallet.NewProvider(kmsClient, nil)

	// Create 3 wallets: account holder, authorized wallet, ordinary(non-authorized) wallet
	accountHolderAlias := "michael"
	accountHolderWallet, err := walletProvider.CreateWallet(ctx, kmswallet.CreateWalletInput{
		Alias:               &accountHolderAlias, // the alias of the key in KMS will be "michael"
		AddWalletAddressTag: true,                // it will add the wallet address into the key tags
		Tags:                map[string]string{"role": "account-holder"},
		// you can set other KMS options
	})

	if err != nil {
		log.Fatalf("account holder wallet creation failed, err: %s", err)
	}

	authorizedWallet, err := walletProvider.CreateWallet(ctx, kmswallet.CreateWalletInput{
		Alias:                           nil,   // It won't have a custom alias
		IgnoreDefaultWalletAddressAlias: false, // It will set the generated wallet address as the alias
		Tags:                            map[string]string{"role": "authorized"},
	})

	if err != nil {
		log.Fatalf("authorized wallet creation failed, err: %s", err)
	}

	ordinaryWallet, err := walletProvider.CreateWallet(ctx, kmswallet.CreateWalletInput{
		Alias:                           nil,
		IgnoreDefaultWalletAddressAlias: true, // This key won't have any alias or tags
	})

	if err != nil {
		log.Fatalf("ordinary wallet creation failed, err: %s", err)
	}

	log.Printf("account holder wallet: %s, authorized wallet: %s, ordinary wallet: %s", accountHolderWallet.Address, authorizedWallet.Address, ordinaryWallet.Address)
	// Before proceeding further, you need to transfer MATIC to the created wallets (from faucets)

	// Initialize the eth rpc client (mumbai network on this case)
	client, err := ethclient.Dial("https://rpc.ankr.com/polygon_mumbai")
	if err != nil {
		log.Fatalf("ethClient initialization failed, err: %s", err)
	}

	defer client.Close()

	// Deploy the BankAccount contract
	accountHolderTransactor, err := walletProvider.GetWalletTransactorByAlias(ctx, accountHolderAlias, chainId)
	if err != nil {
		log.Fatalf("account holder transactor creation failed, err: %s", err)
	}

	contractAddress, tx, bankAccountContract, err := contracts.DeployBankAccount(accountHolderTransactor, client)
	if err != nil {
		log.Fatalf("contract deploy failed, err: %s", err)
	}

	deployReceipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("contract deploy failed, err: %s", err)
	}

	if deployReceipt.Status != types.ReceiptStatusSuccessful {
		log.Fatalf("contract deploy transaction failed, status: %d", deployReceipt.Status)
	}

	log.Printf("bank account contract deployed! address: %s\n", contractAddress.Hex())

	// Deposit balance
	ordinaryTransactor, err := walletProvider.GetWalletTransactor(ctx, ordinaryWallet.KeyId, chainId)
	if err != nil {
		log.Fatalf("ordinary wallet transactor creation failed, err: %s", err)
	}

	ordinaryTransactor.Value = big.NewInt(10000000) // ordinary wallet will deposit 10000000 wei (MATIC)
	depositTx, err := bankAccountContract.Deposit(ordinaryTransactor)
	if err != nil {
		log.Fatalf("deposit failed, err: %s", err)
	}

	depositReceipt, err := bind.WaitMined(ctx, client, depositTx)
	if err != nil {
		log.Fatalf("deposit failed, err: %s", err)
	}

	if depositReceipt.Status != types.ReceiptStatusSuccessful {
		log.Fatalf("deposit transaction failed, status: %d", depositReceipt.Status)
	}

	// View the total balance of the bank account with caller
	accountHolderCaller, err := walletProvider.GetWalletCaller(ctx, accountHolderWallet.KeyId, chainId)
	if err != nil {
		log.Fatalf("account holder caller creation failed, err: %s", err)
	}

	balance, err := bankAccountContract.TotalBalance(accountHolderCaller)
	if err != nil {
		log.Fatalf("total balance read failed, err: %s", err)
	}

	log.Printf("deposited balance: %d\n", balance)

	// Give permission for withdraw to authorized wallet
	withdrawAmount := big.NewInt(50000)
	messageToSign := crypto.Keccak256Hash(
		append(contractAddress.Bytes(), common.HexToAddress(authorizedWallet.Address).Bytes()...),
	)

	signature, err := walletProvider.SignMessage(ctx, accountHolderWallet.KeyId, messageToSign.Bytes())
	if err != nil {
		log.Fatalf("sign message failed, err: %s", err)
	}

	// Withdraw
	authorizedTransactor, err := walletProvider.GetWalletTransactor(ctx, authorizedWallet.KeyId, chainId)
	if err != nil {
		log.Fatalf("authorized wallet transactor creation failed, err: %s", err)
	}

	withdrawTx, err := bankAccountContract.Withdraw(authorizedTransactor, withdrawAmount, signature)
	if err != nil {
		log.Fatalf("withdraw failed, err: %s", err)
	}

	withdrawReceipt, err := bind.WaitMined(ctx, client, withdrawTx)
	if err != nil {
		log.Fatalf("withdraw failed, err: %s", err)
	}

	if withdrawReceipt.Status != types.ReceiptStatusSuccessful {
		log.Fatalf("withdraw transaction failed, status: %d", withdrawReceipt.Status)
	}

	log.Println("withdraw completed!")
}

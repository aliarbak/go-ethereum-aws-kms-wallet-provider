# KMS Wallet Provider Package Example

This example demonstrates the operations that can be performed using the KMS Wallet Provider package. It showcases the usage of the package with a simple `BankAccount` smart contract.

## BankAccount Smart Contract

The smart contract implementation can be found in the `sol/contracts/BankAccount.sol` file. Here is the code:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract BankAccount {
    using ECDSA for bytes32;

    address public accountHolder;
    uint256 public totalBalance;

    constructor() {
        accountHolder = msg.sender;
    }

    function deposit() external payable {
        totalBalance += msg.value;
        // We could use address(this).balance for tracking the balance contract owns,
        // but this is just a sample
    }

    function withdraw(uint256 amount, bytes memory signature) external {
        bytes32 message = keccak256(
            abi.encodePacked(address(this), msg.sender)
        );

        address signer = message.toEthSignedMessageHash().recover(signature);
        require(signer == accountHolder, "Invalid signer");
        // We are checking the signature if the account holder(contract owner) allows msg.sender to withdraw.
        // But this is a dummy logic, do not use it anywhere

        require(totalBalance >= amount, "Insufficient balance");
        totalBalance -= amount;

        (bool sent, ) = msg.sender.call{value: amount}("");
        require(sent, "Withdraw failed");
    }
}
```

The `BankAccount` is a simple bank account contract with the following features:
- The wallet that deploys the contract becomes the account holder.
- Anyone can deposit funds into the contract.
- Only wallets authorized by the account holder can perform withdrawals.

In this example, we will create three wallets:
- Account Holder Wallet: The wallet that deploys the contract and authorizes other wallets for withdrawals.
- Authorized Wallet: A wallet authorized by the account holder to perform withdrawals.
- Ordinary (Non-authorized) Wallet: A wallet that will deposit funds into the contract and query the contract's balance.

### Step 1: Creating the KMS client and KMS Wallet Provider

First, initialize the KMS client using your AWS access credentials and desired configuration. In this example, we use the `eu-central-1` region and provide the AWS access key ID and secret access key.

```go
ctx := context.Background()
chainId := big.NewInt(80001) // means Mumbai Testnet

// Initialize the KMS Client
config := aws.Config{
    Region:      "eu-central-1",
    Credentials: credentials.NewStaticCredentialsProvider("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", ""),
}

kmsClient := kms.NewFromConfig(config) // or you can use kms.New(...)

// Initialize the wallet provider (with default cache duration)
walletProvider := kmswallet.NewProvider(kmsClient, nil)
```

### Step 2: Creating the Account Holder Wallet

Create the account holder wallet by specifying a custom alias (in this example, "michael") and adding a role tag.

```go
accountHolderAlias := "michael"
accountHolderWallet, err := walletProvider.CreateWallet(ctx, kmswallet.CreateWalletInput{
    Alias:               &accountHolderAlias, // the alias of the key in KMS will be "michael"
    AddWalletAddressTag: true,                // it will add the wallet address into the key tags
    Tags:                map[string]string{"role": "account-holder"},
    // you can set other KMS options
})

if err != nil {
    log.Fatalf

("account holder wallet creation failed, err: %s", err)
}
```

### Step 3: Creating the Authorized Wallet

Create the authorized wallet without specifying a custom alias, but set the `IgnoreDefaultWalletAddressAlias` flag to `false`. This will assign the generated wallet address as the alias for the wallet.

```go
authorizedWallet, err := walletProvider.CreateWallet(ctx, kmswallet.CreateWalletInput{
    Alias:                           nil,   // It won't have a custom alias
    IgnoreDefaultWalletAddressAlias: false, // It will set the generated wallet address as the alias
    Tags:                            map[string]string{"role": "authorized"},
})

if err != nil {
    log.Fatalf("authorized wallet creation failed, err: %s", err)
}
```

### Step 4: Creating the Ordinary (Non-authorized) Wallet

Create the ordinary wallet without specifying a custom alias and set the `IgnoreDefaultWalletAddressAlias` flag to `true`. This means the wallet will not have any alias.

```go
ordinaryWallet, err := walletProvider.CreateWallet(ctx, kmswallet.CreateWalletInput{
    Alias:                           nil,
    IgnoreDefaultWalletAddressAlias: true, // This key won't have any alias or tags
})

if err != nil {
    log.Fatalf("ordinary wallet creation failed, err: %s", err)
}
```

Note: Before proceeding to the next steps, you need to load ETH/MATIC into the created wallets to cover gas and deposit fees.

### Step 5: Creating the Ethereum RPC Client

Create an Ethereum RPC client to interact with the Ethereum network. In this example, we use the Mumbai Testnet public RPC node.

```go
client, err := ethclient.Dial("https://rpc.ankr.com/polygon_mumbai")
if err != nil {
    log.Fatalf("ethClient initialization failed, err: %s", err)
}

defer client.Close()
```

### Step 6: Creating the Account Holder Wallet Transactor and Deploying the Bank Account Contract

Create the account holder wallet transactor and deploy the `BankAccount` contract.

```go
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
```

### Step 7: Creating the Ordinary Wallet Transactor and Making a Deposit

Create the ordinary wallet transactor and deposit funds into the bank account contract. In this example, we deposit 10,000,000 wei MATIC.

```go
ordinaryTransactor, err := walletProvider.GetWalletTransactor(ctx, ordinaryWallet.KeyId, chainId)
if err != nil {
    log.Fatalf("ordinary wallet transactor creation failed, err: %s", err)
}

ordinaryTransactor.Value = big.NewInt(10000000) // ordinary wallet will deposit 10,000,000 wei (MATIC)
depositTx, err := bankAccountContract.Deposit(ordinaryTransactor)
if err != nil {
    log.Fatalf("deposit failed, err: %s", err)
}

depositReceipt, err := bind.WaitMined(ctx,

 client, depositTx)
if err != nil {
    log.Fatalf("deposit failed, err: %s", err)
}

if depositReceipt.Status != types.ReceiptStatusSuccessful {
    log.Fatalf("deposit transaction failed, status: %d", depositReceipt.Status)
}
```

### Step 8: Creating the Account Holder Wallet Caller and Retrieving the Total Balance

Create the account holder wallet caller and retrieve the total balance of the bank account.

```go 
accountHolderCaller, err := walletProvider.GetWalletCaller(ctx, accountHolderWallet.KeyId, chainId)
if err != nil {
    log.Fatalf("account holder caller creation failed, err: %s", err)
}

balance, err := bankAccountContract.TotalBalance(accountHolderCaller)
if err != nil {
    log.Fatalf("total balance read failed, err: %s", err)
}
```

### Step 9: Signing the Authorization Message for the Authorized Wallet

Sign the authorization message for the authorized wallet using the account holder wallet.

```go
messageToSign := crypto.Keccak256Hash(
    append(contractAddress.Bytes(), common.HexToAddress(authorizedWallet.Address).Bytes()...),
)

signature, err := walletProvider.SignMessage(ctx, accountHolderWallet.KeyId, messageToSign.Bytes())
if err != nil {
    log.Fatalf("sign message failed, err: %s", err)
}
```

### Step 10: Creating the Authorized Wallet Transactor and Withdrawing Funds

Create the authorized wallet transactor and initiate a withdrawal from the bank account. In this example, we withdraw 50,000 wei MATIC.

```go
authorizedTransactor, err := walletProvider.GetWalletTransactor(ctx, authorizedWallet.KeyId, chainId)
if err != nil {
    log.Fatalf("authorized wallet transactor creation failed, err: %s", err)
}

withdrawAmount := big.NewInt(50000)
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
```

You can access the complete code in the `example/main.go` file.

Please note that in order to execute the code, you need to load ETH/MATIC into the created wallets beforehand to cover gas and deposit fees.
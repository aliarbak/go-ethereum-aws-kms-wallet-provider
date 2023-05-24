// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract BankAccount {
    // We have a bank account,
    // anyone can deposit, 
    // but only authorized wallets can withdraw.
    // Withdrawal permission is checked with the signature of the bank account holder.

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

        require(totalBalance >= amount, "Insufficent balance");
        totalBalance -= amount;

        (bool sent, ) = msg.sender.call{value: amount}("");
        require(sent, "Withdraw failed");
    }
}
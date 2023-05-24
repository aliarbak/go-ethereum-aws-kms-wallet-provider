- npm install
- truffle compile --all

- rm -rf build & solc --optimize --base-path '/' --include-path 'node_modules/' --abi contracts/BankAccount.sol -o build
- solc --optimize --base-path '/' --include-path 'node_modules/' --bin contracts/BankAccount.sol -o build 
- abigen --abi=./build/BankAccount.abi --bin=./build/BankAccount.bin --pkg=bank_account --out=../contracts/bank_account.go
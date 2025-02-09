# <h1 align="center"> NBA TopShot on FlowEVM [Initial Draft Version] </h1>

**! This directory currently contains work in progress only !**

## Introduction

The `BridgedTopShotMoments` smart contract facilitates the creation of 1:1 ERC721 references for existing Cadence-native NBA Top Shot moments. By associating these references with the same metadata, it ensures seamless integration and interaction between Cadence and FlowEVM environments. This allows users to enjoy the benefits of both ecosystems while maintaining the integrity and uniqueness of their NBA Top Shot moments.

## Getting Started

Install Foundry:

```sh
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

Compile contracts and run tests:
```sh
forge test --force -vvv
```

Install Flow CLI: [Instructions](https://developers.flow.com/tools/flow-cli/install)

### Deploy & Verify Contracts

Load environment variables after populating address and key details:

```sh
cp .env.example.testnet .env
source .env
```

Run script to deploy and verify contracts (proxy and implementation):

```sh
forge script --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy script/Deploy.s.sol:DeployScript --broadcast --verify --verifier $VERIFIER_PROVIDER --verifier-url $VERIFIER_URL
```

If verification fails for one or both contracts, verify separately:

```sh
forge verify-contract --rpc-url $RPC_URL --verifier $VERIFIER_PROVIDER --verifier-url $VERIFIER_URL <address-of-contract-to-verify>
```

## Run Transactions

### Direct EVM Calls

Set NFT symbol (admin):

```sh
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy "setSymbol(string)" <new-nft-symbol>
```

### EVM Calls From Cadence

Note: Populate arguments in json file before submitting the transactions.

Bridge NFTs to EVM and wrap:

```sh
flow transactions send ./cadence/transactions/bridge_nfts_to_evm_and_wrap.cdc --args-json "$(cat ./cadence/transactions/bridge_nft_to_evm_and_wrap_args.json)" --network <network> --signer <signer> --gas-limit 8000
```

Wrap NFTs (NFTs already bridged to EVM):

```sh
flow transactions send ./cadence/transactions/wrap_nfts.cdc --args-json "$(cat ./cadence/transactions/wrap_nfts_args.json)" --network <network> --signer <signer>
```

Unwrap NFTs and Bridge NFTs from EVM:

```sh
flow transactions send ./cadence/transactions/unwrap_nfts_and_bridge_from_evm.cdc --args-json "$(cat ./cadence/transactions/unwrap_nfts_and_bridge_from_evm_args.json)" --network <network> --signer <signer> --gas-limit 8000
```

Unwrap NFTs:

```sh
flow transactions send ./cadence/transactions/unwrap_nfts.cdc --args-json "$(cat ./cadence/transactions/unwrap_nfts_args.json)" --network <network> --signer <signer>
```



## Execute Queries

### Direct EVM Calls

BalanceOf:
```sh
cast call $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL "balanceOf(address)(uint256)" $DEPLOYER_ADDRESS
```

OwnerOf:
```sh
cast call $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL "ownerOf(uint256)(address)" <nft-id>
```

### EVM Calls From Cadence

```sh
flow scripts execute ./evm-bridging/cadence/scripts/get_underlying_erc721_address.cdc <nft_contract_flow_address> <nft_contract_evm_address> --network testnet
```

## Misc

Fund testnet Flow EVM account:

1. Use Flow Faucet: https://faucet.flow.com/fund-account

2. Transfer FLOW to EVM address:

```sh
flow transactions send ./cadence/transfer_flow_to_evm_address.cdc <evm_address_hex> <ufix64_amount> --network testnet --signer testnet-account
```

## Useful links

- [Flow Developers Doc - Using Foundry with Flow](https://developers.flow.com/evm/guides/foundry)
- [Flow Developers Doc - Interacting with COAs from Cadence](https://developers.flow.com/evm/cadence/interacting-with-coa)
- [evm-testnet.flowscan.io](https://evm-testnet.flowscan.io)
- [Foundry references](https://book.getfoundry.sh/reference)
- [OpenZeppelin Doc - Foundry Upgrades](https://docs.openzeppelin.com/upgrades-plugins/foundry-upgrades)
- [OpenZeppelin Doc - ERC721 Contracts v5](https://docs.openzeppelin.com/contracts/5.x/api/token/erc721)
- [GitHub - OpenZeppelin Upgradeable Contracts](https://github.com/OpenZeppelin/openzeppelin-contracts-upgradeable)
- [GitHub - LimitBreak Creator Token Standards](https://github.com/limitbreakinc/creator-token-standards)
- [OpenSea Doc - Creator Fee Enforcement](https://docs.opensea.io/docs/creator-fee-enforcement)

# <h1 align="center"> NBA TopShot on FlowEVM </h1>

## Introduction

The `BridgedTopShotMoments` smart contract enables NBA Top Shot moments to exist on FlowEVM as ERC721 tokens. Each ERC721 token is a 1:1 reference to a Cadence-native NBA Top Shot moment, maintaining the same metadata and uniqueness while allowing users to leverage both Flow and EVM ecosystems.

### Core Features

1. **ERC721 Implementation**
   - Full ERC721 compliance with enumeration support
   - Metadata support with customizable base URI
   - Burning capability for token destruction

2. **Bridge Integration**
   - Wrapper functionality for ERC721s from bridged-deployed contract
   - Cross-VM compatibility for Flow ↔ EVM bridging (after bridge upgrade allowing custom associations, and after contract is onboarded to the bridge)
     - Fulfillment of ERC721s from Flow to EVM
     - Bridge permissions management
     - Cadence-specific identifiers tracking

3. **Royalty Management**
   - ERC2981 royalty standard implementation
   - Transfer validation for royalty enforcement via ERC721C/Token Creator Standard
   - Configurable royalty rates (in basis points)
   - Updatable royalty receiver address

> **Note**: This contract is under active development. Features and implementations may change.


## Prerequisites

1. Install Foundry:

```sh
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

2. Install Flow CLI: [Instructions](https://developers.flow.com/tools/flow-cli/install)

## Development

1. Compile and test contracts:

```sh
forge test --force -vvv
```

2. Set up environment:


```sh
cp .env.example.testnet .env
# Add your account details to .env and source it
source .env
```

3. Deploy and verify contracts:

```sh
# Deploy both proxy and implementation contracts
forge script --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy script/Deploy.s.sol:DeployScript --broadcast --verify --verifier $VERIFIER_PROVIDER --verifier-url $VERIFIER_URL

# If verification fails, verify individually
forge verify-contract --rpc-url $RPC_URL --verifier $VERIFIER_PROVIDER --verifier-url $VERIFIER_URL <address-of-contract-to-verify>
```

## Usage

### EVM Operations

```sh
# Approve operator for a NFT
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key <private-key> --legacy "approve(address,uint256)" <operator-address> <token-id>

# Approve operator for all NFTs
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key <private-key> --legacy "setApprovalForAll(address,bool)" <operator-address> <true>

# Transfer NFT
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key <private-key> --legacy "safeTransferFrom(address,address,uint256)" <from-address> <to-address> <token-id>

# Query balance
cast call $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL "balanceOf(address)(uint256)" $DEPLOYER_ADDRESS

# Query owner
cast call $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL "ownerOf(uint256)(address)" <nft-id>

# Query token URI
cast call $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL "tokenURI(uint256)(string)" <nft-id>

# Set NFT symbol (admin only)
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy "setSymbol(string)" <new-nft-symbol>

# Set transfer validator (admin only)
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy "setTransferValidator(address)" <validator-address>

# Set royalty info (admin only)
cast send $DEPLOYED_PROXY_CONTRACT_ADDRESS --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy "setRoyaltyInfo((address,uint96))" "(<royalty-receiver-address>,<royalty-basis-points>)"
```

### Cadence Operations

> **Note**: Populate arguments in json file before submitting the transactions.

```sh
# Transfer erc721 NFTs
flow transactions send ./cadence/transactions/transfer_erc721s_to_evm_address.cdc --args-json "$(cat ./cadence/transactions/transfer_erc721s_to_evm_address_args.json)" --network <network> --signer <signer>

# Bridge and wrap NFTs
flow transactions send ./cadence/transactions/bridge_nfts_to_evm_and_wrap.cdc --args-json "$(cat ./cadence/transactions/bridge_nft_to_evm_and_wrap_args.json)" --network <network> --signer <signer> --gas-limit 8000

# Wrap already-bridged NFTs
flow transactions send ./cadence/transactions/wrap_nfts.cdc --args-json "$(cat ./cadence/transactions/wrap_nfts_args.json)" --network <network> --signer <signer>

# Unwrap and bridge back NFTs
flow transactions send ./cadence/transactions/unwrap_nfts_and_bridge_from_evm.cdc --args-json "$(cat ./cadence/transactions/unwrap_nfts_and_bridge_from_evm_args.json)" --network <network> --signer <signer> --gas-limit 8000

# Unwrap NFTs
flow transactions send ./cadence/transactions/unwrap_nfts.cdc --args-json "$(cat ./cadence/transactions/unwrap_nfts_args.json)" --network <network> --signer <signer>

# Query ERC721 address
flow scripts execute ./evm-bridging/cadence/scripts/get_underlying_erc721_address.cdc <nft_contract_flow_address> <nft_contract_evm_address> --network testnet

# Set up royalty management (admin only)
flow transactions send ./cadence/transactions/admin/set_up_royalty_management.cdc --args-json "$(cat ./cadence/transactions/admin/set_up_royalty_management_args.json)" --network <network> --signer <signer>
```

### Testnet Setup

1. Get testnet FLOW from [Flow Faucet](https://faucet.flow.com/fund-account)

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

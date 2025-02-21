# <h1 align="center"> NBA Top Shot on Flow EVM </h1>

## Introduction

The `BridgedTopShotMoments` smart contract enables NBA Top Shot moments to exist on Flow EVM as ERC721 tokens. Each ERC721 token is a 1:1 reference to a Cadence-native NBA Top Shot moment, maintaining the same metadata and uniqueness while allowing users to leverage both Flow and EVM ecosystems.

### Deployments

|Testnet|Mainnet|
|---|---|
|[0x87859e1d295e5065A2c73be242e3abBd56BAa576](https://evm.flowscan.io/address/0x87859e1d295e5065A2c73be242e3abBd56BAa576)||

### Core Features

1. **ERC721 Implementation**
   - Full ERC721 compliance with enumeration and burning capabilities
   - NFT metadata support with customizable base URI
   - Ownable contract for admin operations
   - Upgradeable via UUPS proxy

2. **Bridge Integration**
   - Wrapper functionality for ERC721s from bridge-deployed contract
   - Cross-VM compatibility for Flow ↔ EVM bridging (after [FLIP-318](https://github.com/onflow/flips/pull/319) implementation allowing custom associations, and after contract is onboarded to the bridge)
     - Fulfillment of ERC721s from Flow to EVM
     - Bridge permissions management
     - Cadence-specific identifiers tracking

3. **Royalty Management**
   - ERC2981 royalty standard implementation
   - Transfer validation for royalty enforcement via ERC721C/Token Creator Standard
   - Configurable royalty rates (in basis points)
   - Updatable royalty receiver address

> **Note**: This contract will be integrated with the Flow EVM bridge once [FLIP-318](https://github.com/onflow/flips/pull/319) is implemented. Currently:
>
> - The contract acts as a wrapper for ERC721s from the bridged-deployed contract
> - Bridging transactions rely on the `CrossVMMetadataViews.EVMPointer` implementation in the Cadence `TopShot` contract and internal logic to determine whether wrapping/unwrapping is necessary
> - After bridge onboarding is complete:
>   - All bridging operations to EVM will use the `BridgedTopShotMoments` contract
>   - Bridging from EVM will support both `BridgedTopShotMoments` and the legacy bridge-deployed contract

## Prerequisites

1. Install Foundry:

```sh
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

2. Install Flow CLI: [Instructions](https://developers.flow.com/tools/flow-cli/install)

## Usage

This section provides commands for interacting with NBA Top Shot Moment ERC721s deployed on Flow EVM using either Cadence operations via `flow` CLI or EVM operations via `cast` CLI.

### Cadence Operations

#### Notes

- Ensure all transaction arguments are populated in the corresponding JSON file template before submission
- If you encounter an `insufficient computation` error, increase the gas limit (i.e., `--gas-limit <new-gas-limit>`)

```sh

# Bridge NFTs to EVM (wraps NFTs if applicable)
flow transactions send ./evm-bridging/cadence/transactions/bridge_nfts_to_evm.cdc --args-json "$(cat ./evm-bridging/cadence/transactions/bridge_nfts_to_evm_args.json)" --network <network> --signer <signer> --gas-limit 8000

# Bridge NFTs from EVM (unwraps NFTs if applicable)
flow transactions send ./evm-bridging/cadence/transactions/bridge_nfts_from_evm.cdc --args-json "$(cat ./evm-bridging/cadence/transactions/bridge_nfts_from_evm_args.json)" --network <network> --signer <signer> --gas-limit 8000

# Transfer erc721 NFTs
flow transactions send ./evm-bridging/cadence/transactions/transfer_erc721s_to_evm_address.cdc --args-json "$(cat ./evm-bridging/cadence/transactions/transfer_erc721s_to_evm_address_args.json)" --network <network> --signer <signer>

# Query ERC721 address
flow scripts execute ./evm-bridging/cadence/scripts/get_evm_address_string.cdc <flow_address> --network testnet

# Set up royalty management (admin only)
flow transactions send ./evm-bridging/cadence/transactions/admin/set_up_royalty_management.cdc --args-json "$(cat ./evm-bridging/cadence/transactions/admin/set_up_royalty_management_args.json)" --network <network> --signer <signer>
```

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

## Development

### Tests

Compile and test contracts:

```sh
forge test --force -vvv
```

### Deploy Using Flow

```sh
# If deploying on emulator, start emulator
flow emulator --config-path ./cadence/transactions/admin/deploy/flow.json --transaction-fees

# Use go 1.22.3
go install golang.org/dl/go1.22.3@latest
go1.22.3 download

# Deploy both proxy and implementation contracts
go1.22.3 run main.go <script-type> <network-name> # for example: go1.22.3 run main.go setup emulator

# If getting the error below:
# vendor/github.com/onflow/crypto/blst_include.h:5:10: fatal error: 'consts.h' file not found
# #include "consts.h"
#
# Try running the following:
CGO_ENABLED=0 go1.22.3 run -tags=no_cgo main.go <script-type> <network-name>
```

### Deploy Using EVM (Initial Testing)

1. Set up environment:


```sh
cp .env.flowevm.testnet.example .env
# Add your account details to .env and source it
source .env
```

2. Deploy and verify contracts:

```sh
# Deploy both proxy and implementation contracts
forge clean
forge script --rpc-url $RPC_URL --private-key $DEPLOYER_PRIVATE_KEY --legacy script/InitialTestingDeploy.s.sol:InitialTestingDeployScript --broadcast --verify --verifier $VERIFIER_PROVIDER --verifier-url $VERIFIER_URL

# If verification fails, verify individually
forge verify-contract --rpc-url $RPC_URL --verifier $VERIFIER_PROVIDER --verifier-url $VERIFIER_URL <address-of-contract-to-verify>
```

## Useful links

- [Flow Faucet](https://faucet.flow.com/fund-account)
- [Flow Developers Doc - Using Foundry with Flow](https://developers.flow.com/evm/guides/foundry)
- [Flow Developers Doc - Interacting with COAs from Cadence](https://developers.flow.com/evm/cadence/interacting-with-coa)
- [evm-testnet.flowscan.io](https://evm-testnet.flowscan.io)
- [Foundry references](https://book.getfoundry.sh/reference)
- [OpenZeppelin Doc - Foundry Upgrades](https://docs.openzeppelin.com/upgrades-plugins/foundry-upgrades)
- [OpenZeppelin Doc - ERC721 Contracts v5](https://docs.openzeppelin.com/contracts/5.x/api/token/erc721)
- [GitHub - OpenZeppelin Upgradeable Contracts](https://github.com/OpenZeppelin/openzeppelin-contracts-upgradeable)
- [GitHub - LimitBreak Creator Token Standards](https://github.com/limitbreakinc/creator-token-standards)
- [OpenSea Doc - Creator Fee Enforcement](https://docs.opensea.io/docs/creator-fee-enforcement)

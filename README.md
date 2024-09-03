# NBA Top Shot

## Introduction

This repository contains the smart contracts and transactions that implement
the core functionality of NBA Top Shot.

The smart contracts are written in Cadence, a new resource oriented
smart contract programming language designed for the Flow Blockchain.

### What is NBA Top Shot

NBA Top Shot is the official digital collectibles
game for the National Basketball Association. Players collect and trade
digital collectibles that represent highlights from the best players 
in the world. See more at nbatopshot.com

### What is Flow?

Flow is a new blockchain for open worlds. Read more about it [here](https://www.onflow.org/).

### What is Cadence?

Cadence is a new Resource-oriented programming language 
for developing smart contracts for the Flow Blockchain.
Read more about it [here](https://www.docs.onflow.org)

We recommend that anyone who is reading this should have already
completed the [Cadence Tutorials](https://docs.onflow.org/cadence) 
so they can build a basic understanding of the programming language.

Resource-oriented programming, and by extension Cadence, 
is the perfect programming environment for Non-Fungible Tokens (NFTs), because users are able
to store their NFT objects directly in their accounts and transact
peer-to-peer. Please see the [blog post about resources](https://medium.com/dapperlabs/resource-oriented-programming-bee4d69c8f8e)
to understand why they are perfect for digital assets like NBA Top Shot Moments.

### Contributing

If you see an issue with the code for the contracts, the transactions, scripts,
documentation, or anything else, please do not hesitate to make an issue or
a pull request with your desired changes. This is an open source project
and we welcome all assistance from the community!

## Top Shot Contract Addresses

`TopShot.cdc`: This is the main Top Shot smart contract that defines
the core functionality of the NFT.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | `0x877931736ee77cff` |
| Mainnet | `0x0b2a3299cc857e29` |

`MarketTopShot.cdc`: This is the top shot marketplace contract that allows users
to buy and sell their NFTs.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | `0x547f177b243b4d80` |
| Mainnet | `0xc1e4f4f4c4257510` |

### Non Fungible Token Standard

The NBA Top Shot contracts utilize the [Flow NFT standard](https://github.com/onflow/flow-nft)
which is equivalent to ERC-721 or ERC-1155 on Ethereum. If you want to build an NFT contract,
please familiarize yourself with the Flow NFT standard before starting and make sure you utilize it 
in your project in order to be interoperable with other tokens and contracts that implement the standard.

### Top Shot Marketplace contract

The top shot marketplace contract was designed in the very early days of Cadence, and therefore
uses some language features that are NOT RECOMMENDED to use by newer projects.
For example, the marketplace contract stores the moments that are for sale in the sale collection.
The correct way to manage this in cadence is to give a collection capability to the market collection
so that the nfts do not have to leave the main collection when going up for sale. The sale collection
would use this capability to withdraw moments from the main collection when they are purchased.

This way, any other smart contracts that need to check a user's account for what they own only need to check
the main collection and not all of the sale collections that could possibly be in their account.

See the [kitty items marketplace contract](https://github.com/onflow/kitty-items/blob/master/cadence/contracts/KittyItemsMarket.cdc) for an example of the current best practices when
it comes to marketplace contracts.

## Directory Structure

The directories here are organized into contracts, scripts, and transactions.

Contracts contain the source code for the Top Shot contracts that are deployed to Flow.

Scripts contain read-only transactions to get information about
the state of someone's Collection or about the state of the TopShot contract.

Transactions contain the transactions that various admins and users can use
to perform actions in the smart contract like creating plays and sets,
minting Moments, and transferring Moments.

 - `contracts/` : Where the Top Shot related smart contracts live.
 - `transactions/` : This directory contains all the transactions and scripts
 that are associated with the Top Shot smart contracts.
 - `transactions/scripts/`  : This contains all the read-only Cadence scripts 
 that are used to read information from the smart contract
 or from a resource in account storage.
 - `lib/` : This directory contains packages for specific programming languages
 to be able to read copies of the Top Shot smart contracts, transaction templates,
 and scripts. Also contains automated tests written in those languages. Currently,
 Go is the only language that is supported, but we are hoping to add javascript
 and other languages soon. See the README in `lib/go/` for more information
 about how to use the Go packages.

## Top Shot Contract Overview

Each Top Shot Moment NFT represents a play from a game in the NBA season.
Plays are grouped into sets which usually have some overarching theme,
like rarity or the type of the play. 

A set can have one or more plays in it and the same play can exist in
multiple sets, but the combination of a play and a set, 
otherwise known as an edition, is unique and is what classifies an individual Moment.

Multiple Moments can be minted from the same edition and each receives a 
serial number that indicates where in the edition it was minted.

Each Moment is a resource object 
with roughly the following structure:

```cadence
pub resource Moment {

    // global unique Moment ID
    pub let id: UInt64
    
    // the ID of the Set that the Moment comes from
    pub let setID: UInt32

    // the ID of the Play that the Moment references
    pub let playID: UInt32

    // the place in the edition that this Moment was minted
    // Otherwise known as the serial number
    pub let serialNumber: UInt32
}
```

The other types that are defined in `TopShot` are as follows:

 - `Play`: A struct type that holds most of the metadata for the Moments.
    All plays in Top Shot will be stored and modified in the main contract.
 - `SetData`: A struct that contains constant information about sets in Top Shot
    like the name, the series, the id, and such.
 - `Set`: A resource that contains variable data for sets 
    and the functionality to modify sets,
    like adding and retiring plays, locking the set, and minting Moments from
    the set.
 - `MomentData`: A struct that contains the metadata associated with a Moment.
    instances of it will be stored in each Moment.
 - `NFT`: A resource type that is the NFT that represents the Moment
    highlight a user owns. It stores its unique ID and other metadata. This
    is the collectible object that the users store in their accounts.
 - `Collection`: Similar to the `NFTCollection` resource from the NFT
    example, this resource is a repository for a user's Moments.  Users can
    withdraw and deposit from this collection and get information about the 
    contained Moments.
 - `Admin`: This is a resource type that can be used by admins to perform
    various actions in the smart contract like starting a new series, 
    creating a new play or set, and getting a reference to an existing set.
 - `QuerySetData`: A struct that contains the metadata associated with a set.
    This is currently the only way to access the metadata of a set.
    Can be accessed by calling the public function in the `TopShot` smart contract called `getSetData(setID)`

Metadata structs associated with plays and sets are stored in the main smart contract
and can be queried by anyone. For example, If a player wanted to find out the 
name of the team that the player represented in their Moment plays for, they
would call a public function in the `TopShot` smart contract 
called `getPlayMetaDataByField`, providing, from their owned Moment,
the play and field that they want to query. 
They can do the same with information about sets by calling `getSetData` with the setID.

The power to create new plays, sets, and Moments rests 
with the owner of the `Admin` resource.

Admins create plays and sets which are stored in the main smart contract,
Admins can add plays to sets to create editions, which Moments can be minted from.

Admins also can restrict the abilities of sets and editions to be further expanded.
A set begins as being unlocked, which means plays can be added to it,
but when an admin locks the set, plays can no longer be added to it. 
This cannot be reversed.

The same applies to editions. Editions start out open, and an admin can mint as
many Moments they want from the edition. When an admin retires the edition, 
Moments can no longer be minted from that edition. This cannot be reversed.

These rules are in place to ensure the scarcity of sets and editions
once they are closed.

Once a user owns a Moment object, that Moment is stored directly 
in their account storage via their `Collection` object. The collection object
contains a dictionary that stores the Moments and gives utility functions
to move them in and out and to read data about the collection and its Moments.

## How to Deploy and Test the Top Shot Contract in VSCode

The first step for using any smart contract is deploying it to the blockchain,
or emulator in our case. Do these commands in vscode. 
See the [vscode extension instructions](https://docs.onflow.org/vscode-extension/) 
to learn how to use it.

 1. Start the emulator with the `Run emulator` vscode command.
 2. Open the `NonFungibleToken.cdc` file from the [flow-nft repo](https://github.com/onflow/flow-nft/blob/master/contracts/NonFungibleToken.cdc) and the `TopShot.cdc` file.  Feel free to read as much as you want to familiarize yourself with the contracts.
 3. In `NonFungibleToken.cdc`, click the `deploy contract to account` 
 above the `Dummy` contract at the bottom of the file to deploy it.
 This also deploys the `NonFungbleToken` interface.
 4. In `TopShot.cdc`, make sure it imports `NonFungibleToken` from the account you deployed it to.
 5. Click the `deploy contract to account` button that appears over the 
    `TopShot` contract declaration to deploy it to a new account.

This deploys the contract code. It also runs the contract's
`init` function, which initializes the contract storage variables,
stores the `Collection` and `Admin` resources 
in account storage, and creates links to the `Collection`.

As you can see, whenever we want to call a function, read a field,
or use a type that is defined in a smart contract, we simply import
that contract from the address it is defined in and then use the imported
contract to access those type definitions and fields.

After the contracts have been deployed, you can run the sample transactions
to interact with the contracts. The sample transactions are meant to be used
in an automated context, so they use transaction arguments and string template
fields. These make it easier for a program to use and interact with them.
If you are running these transactions manually in the Flow Playground or
vscode extension, you will need to remove the transaction arguments and
hard code the values that they are used for. 

You also need to replace the `ADDRESS` placeholders with the actual Flow 
addresses that you want to import from.

## How to Run Transactions Against the Top Shot Contract
This repository contains sample transactions that can be executed against the Top Shot contract either via Flow CLI or using VSCode. This section will describe how to create a new Top Shot set on the Flow emulator.

#### Send Transaction with Flow CLI
1. Install the [Flow CLI and emulator](https://docs.onflow.org/flow-cli/install/)
2. Initialize the flow emulator configuration.  
`flow emulator --init`
3. [Configure the contracts & deployment section](https://docs.onflow.org/flow-cli/configuration/) of the initialized flow.json file. 
4. Start the emulator.  
`flow emulator`
5. On TopShot.cdc substitute the placeholder address `import NonFungibleToken from 0xNFTADDRESS` with the address the NonFungibleToken was deployed to. This will be the emulator address found in the accounts object of the initialized flow.json.
6. Deploy the NonFungibleToken & TopShot contracts to the flow emulator.  
`flow project deploy --network=emulator`
7. Use the Flow CLI to execute transactions against the emulator. This transaction creates a new set on the flow emulator called "new set name".   
`flow transactions send ./transactions/admin/create_set.cdc "new set name"`

#### Send Transaction with VSCode
1. [Install and configure](https://docs.onflow.org/vscode-extension/) VSCode extension.
2. Start flow emulator by running the VSCode command.  
`Cadence: Run emulator`
3. On TopShot.cdc substitute the placeholder address `import NonFungibleToken from 0xNFTADDRESS` with the address the NonFungibleToken was deployed to. Typically, this will be the service account address.
4. Above the contract definition `pub contract interface NonFungibleToken` you will see and press text to deploy this contract to the service account.
5. Above the contract definition `pub contract TopShot: NonFungibleToken` you will see and press text to deploy this contract to the service account.
6. Navigate to `transactions/admin/create_set.cdc` Substitute the placeholder address `import TopShot from 0xNFTADDRESS` with the address TopShot.cdc was deployed to.
7. Transactions run in VSCode cannot take arguments. Replace the line `transaction(setName : String)` with `transaction()` and find every instance of setName in the contract and replace with a hard coded value like "new set name".
8. Above the line `transaction()` you will now see and press the text `Send signed by service account`. This will create a set on the flow emulator called "new set name".

## How to run the automated tests for the contracts

See the `lib/go` README for instructions about how to run the automated tests.

## Instructions for creating plays and minting moments

A common order of creating new Moments would be

1. Creating new plays with `transactions/admin/create_play.cdc`.
2. Creating new sets with `transactions/admin/create_set.cdc`.
3. Adding plays to the sets to create editions
   with `transactions/admin/add_plays_to_set.cdc`.
4. Minting Moments from those editions with 
   `transactions/admin/batch_mint_moment.cdc`.

You can also see the scripts in `transactions/scripts` to see how information
can be read from the real Top Shot smart contract deployed on the
Flow Beta Mainnet. 

### Accessing the NBA Top Shot smart contract on Flow Beta Mainnet

The Flow Beta mainnet is still a work in progress and still has
a limited number of accounts that can run nodes and submit transactions.
Anyone can read data from the contract by running any of the scripts in the 
`transactions` directory using one of the public access nodes.

For example, this is how you would query the total supply via the Flow CLI.

`flow scripts execute transactions/scripts/get_totalSupply.cdc --host access.mainnet.nodes.onflow.org:9000`

Make sure that the import address in the script is correct for mainnet.

## NBA Top Shot Events

The smart contract and its various resources will emit certain events
that show when specific actions are taken, like transferring an NFT. This
is a list of events that can be emitted, and what each event means.
You can find definitions for interpreting these events in Go by seeing
the `lib/go/events` package.
    

#### Events for plays
- `pub event PlayCreated(id: UInt32, metadata: {String:String})`
    
    Emitted when a new play has been created and added to the smart contract by an admin.

- `pub event NewSeriesStarted(newCurrentSeries: UInt32)`
    
    Emitted when a new series has been triggered by an admin.

#### Events for set-Related actions

- `pub event SetCreated(setID: UInt32, series: UInt32)`
    
    Emitted when a new set is created.
    
- `pub event PlayAddedToSet(setID: UInt32, playID: UInt32)`
    
    Emitted when a new play is added to a set.
    
- `pub event PlayRetiredFromSet(setID: UInt32, playID: UInt32, numMoments: UInt32)`

    Emitted when a play is retired from a set. Indicates that 
    that play/set combination and cannot be used to mint moments any more.
    
- `pub event SetLocked(setID: UInt32)`

    Emitted when a set is locked, meaning plays cannot be added.
    
- `pub event MomentMinted(momentID: UInt64, playID: UInt32, setID: UInt32, serialNumber: UInt32)`

    Emitted when a Moment is minted from a set. The `momentID` is the global unique identifier that differentiates a Moment from all other Top Shot Moments in existence. The `serialNumber` is the identifier that differentiates the Moment within an Edition. It corresponds to the place in that edition where it was minted. 

#### Events for Collection-related actions
    
- `pub event Withdraw(id: UInt64, from: Address?)`

    Emitted when a Moment is withdrawn from a collection. `id` refers to the global Moment ID. If the collection was in an account's storage when it was withdrawn, `from` will show the address of the account that it was withdrawn from. If the collection was not in storage when the Moment was withdrawn, `from` will be `nil`.

- `pub event Deposit(id: UInt64, to: Address?)`

    Emitted when a Moment is deposited into a collection. `id` refers to the global Moment ID. If the collection was in an account's storage when it was deposited, `to` will show the address of the account that it was deposited to. If the collection was not in storage when the Moment was deposited, `to` will be `nil`.

### Top Shot NFT Metadata

NFT metadata is represented in a flexible and modular way using the [standard proposed in FLIP-0636](https://github.com/onflow/flow/blob/master/flips/20210916-nft-metadata.md). The Top Shot contract implements the [`MetadataViews.Resolver`](https://github.com/onflow/flow-nft/blob/master/contracts/MetadataViews.cdc#L21) interface, which standardizes the display of Top Shot NFT in accordance with FLIP-0636. The Top Shot contract also defines a custom view of moment play data called TopShotMomentMetadataView.

## NBA Top Shot Packs

NBA Top Shot packs are currently off-chain and not managed by the NBA Top Shot smart contract. Moments in a pack are minted on-chain, and assembled into a pack for purchase off-chain on the NBA Top Shot platform. When a collector purchases a pack, the moments within the pack are transferred directly to this collector on-chain. The NBA Top Shot smart contract has no knowledge of packs.

## NBA Top Shot Marketplace

The `contracts/MarketTopShot.cdc` contract allows users to create a sale object
in their account to sell their Moments.

When a user wants to sell their Moment, they create a sale collection
in their account and specify a beneficiary of a cut of the sale if they wish.

A Top Shot Sale Collection contains a capability to the owner's moment collection
that allows the sale to withdraw the moment when it is purchased.

When another user wants to buy the Moment that is for sale, they simply send 
their fungible tokens to the `purchase` function 
and if they sent the correct amount, they get the Moment back.

#### Events for Market-related actions

- `pub event MomentListed(id: UInt64, price: UFix64, seller: Address?)`
   
   Emitted when a user lists a Moment for sale in their SaleCollection.

- `pub event MomentPriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)`

   Emitted when a user changes the price of their Moment.

- `pub event MomentPurchased(id: UInt64, price: UFix64, seller: Address?)`

   Emitted when a user purchases a Moment that is for sale.

- `pub event MomentWithdrawn(id: UInt64, owner: Address?)`

   Emitted when a seller withdraws their Moment from their SaleCollection.

- `pub event CutPercentageChanged(newPercent: UFix64, seller: Address?)`

   Emitted when a seller changes the percentage cut that is taken
   from their sales and sent to a beneficiary.

### Different Versions of the Market Contract

There are two versions of the Top Shot Market Contract.
`TopShotMarket.cdc` is the original version of the contract that was used
for the first set of sales in the p2p marketplace, but we made improvements
to it which are now in `TopShotMarketV3.cdc`.

There is also a V2 version that was deployed to mainnet, but will never be used.

Both versions define a `SaleCollection` resource that users store in their account.
The resource manages the logic of the sale like listings, de-listing, prices, and 
purchases. The first version actually stores the moments that are for sale, but 
we realized that this causes issues if other contracts need to access a user's
main collection to see what they own. We created the second version to simply
store a capability to the owner's moment collection so that the moments 
that are for sale do not need to be removed from the main collection to be
put up for sale. In this version, when a moment is purchased, the sale collection
uses the capability to withdraw the moment from the main collection and 
returns it to the buyer.

The new version of the market contract is currently NOT DEPLOYED to mainnet,
but it will be deployed and utilized in the near future.

## TopShot contract improvement
Some improvements were made to the Topshot contract to reflect some cadence best practices and fix a bug.
In-depth explanation on the changes and why we made them can be found in our [Blog Post](https://blog.nbatopshot.com/posts/nba-top-shot-smart-contract-improvements) 

## TopShot Locking Contract Overview

Contract Name: `TopShotLocking`

TopShot NFTs can be locked for a duration meaning they are unable to be withdrawn, listed for sale, burned, etc. 
In the NBA TopShot product users are rewarded for locking their moments.

An NFT may be unlocked after the lock duration has passed, or the contract admin has marked it eligible for unlocking.

The moment is locked even if expiry has passed until the owner requests it be unlocked.
The address which owns the locked NFT must make an unlocking transaction once it is eligible.

### Available functions

#### lockNFT
`pub fun lockNFT(nft: @NonFungibleToken.NFT, expiryTimestamp: UFix64): @NonFungibleToken.NFT`  
Takes a TopShot.NFT resource and sets it in the lockedNFTs dictionary, the value of the entry is the expiry timestamp  
Params:  
`nft` - a `NonFungibleToken.NFT` resource, but must conform to `TopShot.NFT` asserted at runtime  
`expiryTimestamp` - the unix timestamp in seconds at which this nft can be unlocked

Example:
```cadence
let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

let ONE_YEAR_IN_SECONDS: UFix64 = UFix64(31536000)
collectionRef.lock(id: 1, duration: ONE_YEAR_IN_SECONDS)
```

#### unlockNFT
`pub fun unlockNFT(nft: @NonFungibleToken.NFT): @NonFungibleToken.NFT`  
Takes a `NonFungibleToken.NFT` resource and attempts to remove it from the lockedNFTs dictionary.
This function will panic if the nft lock has not expired or been overridden by an admin.
Params:  
`nft` - a `NonFungibleToken.NFT` resource 

Example:
```cadence
let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

collectionRef.unlock(id: 1)
```

#### isLocked
`pub fun isLocked(nftRef: &NonFungibleToken.NFT): Bool`  
Returns true if the moment is locked

#### getLockExpiry
`pub fun getLockExpiry(nftRef: &NonFungibleToken.NFT): UFix64`  
Returns the unix timestamp when the nft is eligible for unlock

### Admin Functions

#### markNFTUnlockable
`pub fun markNFTUnlockable(nftRef: &NonFungibleToken.NFT)`  
Places the nft id in an unlockableNFTs dictionary. This dictionary is checked in the `unlockNFT` function and bypasses the `expiryTimestamp`
Params:  
`nftRef` - a reference to an `NonFungibleToken.NFT` resource  

Example:
```cadence
let adminRef: &NFTLocking.Admin

prepare(acct: AuthAccount) {
    // Set TopShotLocking admin ref
    self.adminRef = acct.borrow<&NFTLocking.Admin>(from: /storage/TopShotLockingAdmin)!
}

execute {
    // Set Top Shot NFT Owner collection ref
    let owner = getAccount(0x179b6b1cb6755e31)
    let collectionRef = owner.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not reference owner's moment collection")

    let nftRef = collectionRef.borrowNFT(id: 1)
    self.adminRef.markNFTUnlockable(nftRef: nftRef)
}
```
### Contracts Honoring the Lock

- TopShot `withdraw`
- MarketTopShot relies on the NFT being withdrawn first so no additional code is needed
- TopShotMarketV3 `listForSale`

## License 

The works in these folders 
/dapperlabs/nba-smart-contracts/blob/master/contracts/TopShot.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/MarketTopShot.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/MarketTopShotV3.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/TopShotAdminReceiver.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/TopShotShardedCollection.cdc 

are under the Unlicense
https://github.com/onflow/flow-NFT/blob/master/LICENSE












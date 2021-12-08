# NBA Top Shot

## Introduction

This repository contains the smart contracts and transactions that implement
the core functionality of NBA Top Shot.

The smart contracts are written in Cadence, a new resource oriented
smart contract programming language designed for the Flow Blockchain.

### What is NBA Top Shot

NBA Top Shot is the official digital collecitibles
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
the state of someones Collection or about the state of the TopShot contract.

Transactions contain the transactions that various admins and users can use
to perform actions in the smart contract like creating plays and sets,
minting Moments, and transfering Moments.

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
    // Otherwise know as the serial number
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
    various acitions in the smart contract like starting a new series, 
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
See the [vscode extension instructions](https://docs.onflow.org/docs/visual-studio-code-extension) 
to learn how to use it.

 1. Start the emulator with the `Run emulator` vscode command.
 2. Open the `NonFungibleToken.cdc` file from the [flow-nft repo](https://github.com/onflow/flow-nft/blob/master/contracts/NonFungibleToken.cdc) and the `TopShot.cdc` file.  Feel free to read as much as you want to familiarize yourself with the contracts.
 3. In `NonFungibleToken.cdc`, click the `deploy contract to account` 
 above the `Dummy` contract at the bottom of the file to deploy it.
 This also deploys the `NonFungbleToken` interface.
 4. Switch to a different account.
 5. In `TopShot.cdc`, make sure it imports `NonFungibleToken` from the account you deployed it to.
 6. Click the `deploy contract to account` button that appears over the 
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
    
- `pub event ContractInitialized()`
    
    This event is emitted when the `TopShot` contract is created.

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
These improvements are dscribed below

### Make dictionary and array fields private
Following cadence best practices in [Cadence Anti-Patterns](https://docs.onflow.org/cadence/anti-patterns/#array-or-dictionary-fields-should-be-private),
variables `plays, retired, numberMintedPerPlay` in the `Set` resource were changed from `pub` to `access(contract)`.
This makes it impossible for `Admin` to directly modify these fields. Any modification will have to be done
through the methods defined in the `Set` resource.

We now have
```
pub resource Set {
    access(contract) var plays: [UInt32]
    access(contract) var retired: {UInt32: Bool}
    access(contract) var numberMintedPerPlay: {UInt32: UInt32}
}
```

### Unified Set metadata struct
In addition to `SetData` (which records the id, name and series of a set), and the `Set` resource
(which records other information about the set and acts as an authorization resource for the admin
to create editions, mint moments, retire plays, and more), a new struct `QuerySetData` was added.
```
pub struct QuerySetData {
    pub let setID: UInt32
    pub let name: String
    pub let series: UInt32
    access(self) var plays: [UInt32]
    access(self) var retired: {UInt32: Bool}
    pub var locked: Bool
    access(self) var numberMintedPerPlay: {UInt32: UInt32}
}
```
This new struct consolidates all the important information about a set and can be queried using
`TopShot.getSetData(setID: UInt32): QuerySetData?` method. This makes it easier to get a `Set`
information instead of having to call multiple methods and stitching together their responses

### Perform state changing operations in admin resources, not in public structs
state changing operations like incrementing `TopShot.nextPlayID` and emitting `PlayCreated` event were moved to be done by the admin only
For instance
```
TopShot.nextPlayID = TopShot.nextPlayID + UInt32(1)
emit PlayCreated(id: newPlay.playID, metadata: metadata)
```
was moved to `Admin.createPlay()`

### Borrow references to resources instead of loading them from storage
Instead of the inefficient loading of a set from storage to read it's fields and putting it back. We moved to borrowing a reference to the resource.

Example: Instead of

```
if let setToRead <- TopShot.sets.remove(key: setID) {
    // See if the Play is retired from this Set
    let retired = setToRead.retired[playID]
    
    // Put the Set back in the contract storage
    TopShot.sets[setID] <-! setToRead
    
    // Return the retired status
    return retired
}
```
do
```
let set = &TopShot.sets[setID] as! &Set
// See if the Play is retired from this Set
let retired = set.retired[playID]
return retired
``` 
In-depth explanation on these changes and why we made them can be found in our [Blog Post](https://blog.nbatopshot.com/posts/nba-top-shot-smart-contract-improvements) 

## License 

The works in these folders 
/dapperlabs/nba-smart-contracts/blob/master/contracts/TopShot.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/MarketTopShot.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/MarketTopShotV3.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/TopShotAdminReceiver.cdc 
/dapperlabs/nba-smart-contracts/blob/master/contracts/TopShotShardedCollection.cdc 

are under the Unlicense
https://github.com/onflow/flow-NFT/blob/master/LICENSE












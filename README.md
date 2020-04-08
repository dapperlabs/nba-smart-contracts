# NBA Top Shot Smart Contracts

## Table of Contents

## Introduction

This repository contains the  NBA top shot smart contracts. 
The contracts currently implement all the functionality required for the game,
but will continue to go through changes as the design solidifies and 
as Cadence evolves as a smart contract programming language.

Features:

Admins create Plays and Sets which are stored in the main smart contract,
Admins can add plays to Sets to create editions, 
which moments can be minted from. 

Users can own and transfer moments but using the Collection resource.

By following this tutorial, you should be able to get an understanding of
how the topshot smart contracts work.  

Before you read this tutorial, you should be familiar with Flow and the 
Cadence Programming Language.  

 - [Read the Flow Primer](https://www.withflow.org/en/primer)
 - [Complete the Flow Developer Preview to learn the basics of Cadence](docs.onflow.org)


## Contract Overview

All functionality and type definitions are included in the `TopShot.cdc` contract.

The TopShot contract defines  types.

 - `Play`: A struct type that holds most of the metadata for the moments.
    All plays in Top Shot will be stored and modified in the main contract.
 - `SetData`: A struct that contains constant information about sets in topshot
    like the name, the series, the id, and such.
 - `Set`: A resource that contains functionality to modify sets,
    like adding and removing plays, locking the set, and minting moments from
    the set.
 - `MomentData`: A struct that contains the metadata associated with a moment.
    instances of it will be stored in each moment.
 - `NFT`: A resource type that is the NFT that represents the Moment
    highlight a user owns. It stores its unique ID and other metadata.
 - `Collection`: Similar to the `NFTCollection` resource from the NFT
    example, this resource is a repository for a user's moments.  Users can
    withdraw and deposit from this collection and get information about the 
    contained moments.
 - `Admin`: This is a resource type that can be used by admins to perform
    various acitions in the smart contract like starting a new series, 
    creating a new play or set, and getting a reference to an existing set.

## How to Deploy and Test the TopShot Contract

The first step for using any smart contract is deploying it to the blockchain,
or emulator in our case. Do these commands in vscode

 1. Start the emulator with the `Run emulator` vscode command.
 2. Open the `NonFungibleToken.cdc` file and the `TopShot.cdc` file.  Feel free to read as much as you want to
    familiarize yourself with the contract
 3. The Marketplace smart contract implements the fungible token interface in `fungible-token.cdc`, so you need
    to open that first and click the `deploy contract to account 0x01` button 
    that appears above the `FlowToken` contract. This will deploy the interface definition and contract
    to account 1.
 4. Run the `switch account` command from the vscode comman palette.  Switch to account 2.
 5. In `NonFungibleToken.cdc`, click the `deploy contract to account` to deploy to account 2.
 6. Switch to account 3.
 6. In `topshot.cdc`, click the `deploy contract to account` button that appears over the 
    `TopShot` contract declaration to deploy to account 3.

This deploys the contract code to account 3. It also runs the contracts
`init` function, which initializes the contract storage variables,
stores the `Collection` and `Admin` resources 
in account storage, and stores references to `Collection`.

As you can see, whenever we want to call a function, read a field,
or use a type that is defined in a smart contract, we simply import
that contract from the address it is defined in and then use the imported
contract to access those fields.

## TopShot Events

The smart contract and its various resources will emit certain events
that show when specific actions are taken, like transferring an NFT. This
is a list of events that can be emitted, and what each event means.

    
- `pub event ContractInitialized()`
    
    This event is emitted when the TopShot contract is created

#### Events for plays
- `pub event PlayCreated(id: UInt32, metadata: {String:String})`
    
    Emitted when a new play has been created and added to the smart contract by an admin.

- `pub event NewSeriesStarted(newCurrentSeries: UInt32)`
    
    Emitted when a new series has been triggered by an admin

#### Events for Set-Related actions

- `pub event SetCreated(setID: UInt32, series: UInt32)`
    
    Emitted when a new Set is created
- `pub event PlayAddedToSet(setID: UInt32, playID: UInt32)`
    
    Emitted when a new play is added to a set.
    
- `pub event PlayRetiredFromSet(setID: UInt32, playID: UInt32, numMoments: UInt32)`

    Emitted when a play is retired from a set and cannot be used to mint
    
- `pub event SetLocked(setID: UInt32)`

    Emitted when a set is locked, meaning plays cannot be added
    
- `pub event MomentMinted(momentID: UInt64, playID: UInt32, setID: UInt32, serialNumber: UInt32)`

    Emitted when a moment is minted from a set. The `momentID` is the global unique identifier that differentiates a moment from all other TopShot moments in existence. The `serialNumber` is the identifier that differentiates the moment within an Edition. It corresponds to the place in that edition where it was minted. 

#### Events for Collection-related actions
    
- `pub event Withdraw(id: UInt64, from: Address?)`

    Emitted when a moment is withdrawn from a collection. `id` refers to the global moment ID. If the collection was in an account's storage when it was withdrawn, `from` will show the address of the account that it was withdrawn from. If the collection was not in storage when the Moment was withdrawn, `from` will be `nil`

- `pub event Deposit(id: UInt64, to: Address?)`

    Emitted when a moment is deposited into a collection. `id` refers to the global moment ID. If the collection was in an account's storage when it was deposited, `to` will show the address of the account that it was deposited to. If the collection was not in storage when the Moment was deposited, `to` will be `nil`

## Directory Structure

The directories here are organized into scripts and transactions.

Scripts contain read-only transactions to get information about
the state of someones Collection or about the state of the TopShot contract.

Transactions contain the transactions that various admins and users can use
to performa actions in the smart contract like creating plays and sets,
minting moments, and transfering moments.

 - `contracts/` : Where the TopShot related smart contracts live
 - `scripts/`  : This contains all the read-only Cadence scripts 
 that are used to read information from the smart contract
 or from a resource in account storage
   - `collections/`: Used to read information about a user's Moment collection
   - `market/`: Used to read information about the market smart contract and about a user's Moment sales.
   - `plays/`: Used to read information related to play data and IDs.
   - `sets/`: Used to read various information about set names and IDs
 - `test/`  : This contains special transactions and scripts that are used
 to quickly test the functionality of the smart contracts. They are not all 
 working now and will be updated to be extensive when the contracts
 are updated with the new storage interface changes.
 - `transactions/` : This directory contains all the state-changing transactions
 that are associated with the TopShot smart contracts.
   - `collections/`: Transactions for a user to interact with their Moment collection
   - `market/`: Transactions for users to buy and sell moments using
   the market smart contract.
   - `plays/`: Transactions for an Admin to add and modify play data and IDs.
   - `sets/`: Transactions for an Admin to take all 
   the actions associated wit sets, like adding plays, minting moments,
   locking sets, and more.

### Marketplace

The `topshot_market.cdc` contract allows users to create a marketplace object in their account to sell their moments.

1. Make sure you have followed the steps to get topshot set up.
2. Deploy `fungible-token.cdc` to account 1
3. Deploy `NonFungibleToken.cdc` to account 2 and `TopShot.cdc` to account 3.
4. Deploy `MarketTopShot.cdc` to account 4. Feel free to look at the various 
   fields and functions in the smart contract.

There currently aren't many example transactions for the market but they will be added soon once we have a better idea of exactly how it will function.

#### Events for Market-related actions

- `pub event MomentListed(id: UInt64, price: UFix64, seller: Address?)`
   
   Emitted when a user lists a moment for sale in their SaleCollection.

- `pub event PriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)`

   Emitted when a user changes the price of their moment.

- `pub event TokenPurchased(id: UInt64, price: UFix64, seller: Address?)`

   Emitted when a user purchases a moment that is for sale.

- `pub event SaleWithdrawn(id: UInt64, owner: Address?)`

   Emitted when a seller withdraws their moment from their SaleCollection

- `pub event CutPercentageChanged(newPercent: UFix64, seller: Address?)`

   Emitted when a seller changes the percentage cut that is taken
   from their sales and sent to a beneficiary.















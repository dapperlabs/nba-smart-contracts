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
 2. Open the `NonFungibleToken.cdc` file from the [flow-nft repo](https://github.com/onflow/flow-nft/blob/master/src/contracts/NonFungibleToken.cdc) and the `TopShot.cdc` file.  Feel free to read as much as you want to
    familiarize yourself with the contract
 3. In `NonFungibleToken.cdc`, click the `deploy contract to account` to deploy it.
 4. Switch to a different account
 5. In `TopShot.cdc`, make sure it imports `NonFungibleToken` from the account you deployed it to
 6. click the `deploy contract to account` button that appears over the 
    `TopShot` contract declaration to deploy it to a new account

This deploys the contract code. It also runs the contracts
`init` function, which initializes the contract storage variables,
stores the `Collection` and `Admin` resources 
in account storage, and creates links to the `Collection`.

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

The directories here are organized into contrats, scripts, and transactions.

Contracts contain the source code for the topshot contracts that are deployed to Flow.

Scripts contain read-only transactions to get information about
the state of someones Collection or about the state of the TopShot contract.

Transactions contain the transactions that various admins and users can use
to performa actions in the smart contract like creating plays and sets,
minting moments, and transfering moments.

 - `contracts/` : Where the TopShot related smart contracts live
 - `scripts/`  : This contains all the read-only Cadence scripts 
 that are used to read information from the smart contract
 or from a resource in account storage
 - `transactions/` : This directory contains all the state-changing transactions
 that are associated with the TopShot smart contracts.


### Marketplace

The `MarketTopShot.cdc` contract allows users to create a marketplace object in their account to sell their moments.

#### Events for Market-related actions

- `pub event MomentListed(id: UInt64, price: UFix64, seller: Address?)`
   
   Emitted when a user lists a moment for sale in their SaleCollection.

- `pub event MomentPriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)`

   Emitted when a user changes the price of their moment.

- `pub event MomentPurchased(id: UInt64, price: UFix64, seller: Address?)`

   Emitted when a user purchases a moment that is for sale.

- `pub event MomentWithdrawn(id: UInt64, owner: Address?)`

   Emitted when a seller withdraws their moment from their SaleCollection

- `pub event CutPercentageChanged(newPercent: UFix64, seller: Address?)`

   Emitted when a seller changes the percentage cut that is taken
   from their sales and sent to a beneficiary.















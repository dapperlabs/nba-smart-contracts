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
    All molds in Top Shot will be stored and modified in the main contract.
    A `Mold` object contains fields for its ID, the number of moments that can
    be minted from each quality, the number that can still be minted
    of each quality, and a field for all the mold's metadata. 
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

### Directory Structure

The directories here are organized into scripts and transactions.

Scripts contain read-only transactions to get information about
the state of someones Collection or about the state of the TopShot contract.

Transactions contain the transactions that various admins and users can use
to performa actions in the smart contract like creating plays and sets,
minting moments, and transfering moments.


### Marketplace

The `topshot_market.cdc` contract allows users to create a marketplace object in their account to sell their moments.

1. Make sure you have followed the steps to get topshot set up.
2. Deploy `fungible-token.cdc` to account 1
3. Deploy `NonFungibleToken.cdc` to account 2 and `TopShot.cdc` to account 3.
4. Deploy `MarketTopShot.cdc` to account 4. Feel free to look at the various 
   fields and functions in the smart contract.

There currently aren't many example transactions for the market but they will be added soon.















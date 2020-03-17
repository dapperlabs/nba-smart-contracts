
pub contract interface NonFungibleToken {

    // The total number of tokens of this type in existance
    pub var totalSupply: UInt64

    pub event ContractInitialized()
    pub event Withdraw(id: UInt64)
    pub event Deposit(id: UInt64)

    pub resource interface INFT {
        // The unique ID that each NFT has
        pub let id: UInt64

        // placeholder for token metadata 
        pub var metadata: {String: String}
    }

    pub resource NFT: INFT {
        pub let id: UInt64

        pub var metadata: {String: String}
    }

    pub resource interface Provider {
        // withdraw removes an NFT from the collection and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NFT {
            post {
                result.id == withdrawID: "The ID of the withdrawn token must be the same as the requested ID"
            }
        }

        pub fun batchWithdraw(ids: [UInt64]): @Collection
    }

    pub resource interface Receiver {

		pub fun deposit(token: @NFT) 

        pub fun batchDeposit(tokens: @Collection)
    }

    pub resource interface Metadata {

		pub fun getIDs(): [UInt64]
	}

    pub resource Collection: Provider, Receiver, Metadata {
        
        pub var ownedNFTs: @{UInt64: NFT}

        // withdraw removes an NFT from the collection and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NFT 

        pub fun batchWithdraw(ids: [UInt64]): @Collection

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @NFT)

        pub fun batchDeposit(tokens: @Collection)

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64]
    }

    pub fun createEmptyCollection(): @Collection {
        post {
            result.getIDs().length == 0: "The created collection must be empty!"
        }
    }
}

// For deploying to vscode extension
pub contract Dummy {}
 
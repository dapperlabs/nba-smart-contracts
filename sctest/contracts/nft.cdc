
access(all) contract interface NonFungibleToken {

    // The total number of tokens of this type in existance
    pub var totalSupply: UInt64

    access(all) event ContractInitialized()
    access(all) event Withdraw(id: UInt64)
    access(all) event Deposit(id: UInt64)

    pub resource interface INFT {
        // The unique ID that each NFT has
        access(all) let id: UInt64

        // placeholder for token metadata 
        access(all) var metadata: {String: String}
    }

    access(all) resource NFT: INFT {
        access(all) let id: UInt64

        access(all) var metadata: {String: String}
    }

    access(all) resource interface Provider {
        // withdraw removes an NFT from the collection and moves it to the caller
        access(all) fun withdraw(withdrawID: UInt64): @NFT {
            post {
                result.id == withdrawID: "The ID of the withdrawn token must be the same as the requested ID"
            }
        }

        access(all) fun batchWithdraw(ids: [UInt64]): @Collection
    }

    access(all) resource interface Receiver {

		access(all) fun deposit(token: @NFT) 

        access(all) fun batchDeposit(tokens: @Collection)
    }

    access(all) resource interface Metadata {

		access(all) fun getIDs(): [UInt64]
	}

    access(all) resource Collection: Provider, Receiver, Metadata {
        
        access(all) var ownedNFTs: @{UInt64: NFT}

        // withdraw removes an NFT from the collection and moves it to the caller
        access(all) fun withdraw(withdrawID: UInt64): @NFT 

        access(all) fun batchWithdraw(ids: [UInt64]): @Collection

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        access(all) fun deposit(token: @NFT)

        access(all) fun batchDeposit(tokens: @Collection)

        // getIDs returns an array of the IDs that are in the collection
        access(all) fun getIDs(): [UInt64]
    }

    access(all) fun createEmptyCollection(): @Collection {
        post {
            result.getIDs().length == 0: "The created collection must be empty!"
        }
    }
}

access(all) contract Tokens: NonFungibleToken {

    access(all) var totalSupply: UInt64

    access(all) event ContractInitialized()
    access(all) event Withdraw(id: UInt64)
    access(all) event Deposit(id: UInt64)

    access(all) resource NFT: NonFungibleToken.INFT {
        access(all) let id: UInt64

        access(all) var metadata: {String: String}

        init(initID: UInt64) {
            self.id = initID
            self.metadata = {}
        }
    }

    access(all) resource Collection: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.Metadata {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `UInt64` ID field
        access(all) var ownedNFTs: @{UInt64: NFT}

        init () {
            self.ownedNFTs <- {}
        }

        // withdraw removes an NFT from the collection and moves it to the caller
        access(all) fun withdraw(withdrawID: UInt64): @NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            emit Withdraw(id: token.id)

            return <-token
        }

        access(all) fun batchWithdraw(ids: [UInt64]): @Collection {
            var i = 0
            var batchCollection: @Collection <- create Collection()

            while i < ids.length {
                batchCollection.deposit(token: <-self.withdraw(withdrawID: ids[i]))

                i = i + 1
            }
            return <-batchCollection
        }

        // deposit takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        access(all) fun deposit(token: @NFT) {
            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id)

            destroy oldToken
        }

        access(all) fun batchDeposit(tokens: @Collection) {
            var i = 0
            let keys = tokens.getIDs()

            while i < keys.length {
                self.deposit(token: <-tokens.withdraw(withdrawID: keys[i]))

                i = i + 1
            }
            destroy tokens
        }

        // idExists checks to see if a NFT with the given ID exists in the collection
        access(all) fun idExists(id: UInt64): Bool {
            return self.ownedNFTs[id] != nil
        }

        // getIDs returns an array of the IDs that are in the collection
        access(all) fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        access(all) fun getMetaData(id: UInt64, field: String): String {
            let token <- self.ownedNFTs.remove(key: id) ?? panic("No NFT!")
            
            let dataOpt = token.metadata[field]

            let oldToken <- self.ownedNFTs[id] <- token
            destroy oldToken

            if let data = dataOpt {
                return data
            } else {
                return "None"
            }
        }

        destroy() {
            destroy self.ownedNFTs
        }
    }

    access(all) fun createNFT(id: UInt64): @NFT {
        return <- create NFT(initID: id)
    }

    access(all) fun createEmptyCollection(): @Collection {
        return <- create Collection()
    }

	access(all) resource NFTFactory {

		// the ID that is used to mint moments
		access(all) var idCount: UInt64

		init() {
			self.idCount = 1
		}

		// mintNFT mints a new NFT with a new ID
		// and deposit it in the recipients colelction using their collection reference
		access(all) fun mintNFT(recipient: &Collection) {

			// create a new NFT
			var newNFT <- create NFT(initID: self.idCount)
			
			// deposit it in the recipient's account using their reference
			recipient.deposit(token: <-newNFT)

			// change the id so that each ID is unique
			self.idCount = self.idCount + UInt64(1)
		}
	}

	init() {
        self.totalSupply = 0
        
		let oldCollection <- self.account.storage[Collection] <- create Collection()
		destroy oldCollection

		self.account.storage[&Collection] = &self.account.storage[Collection] as &Collection
        self.account.published[&NonFungibleToken.Receiver] = &self.account.storage[Collection] as &NonFungibleToken.Receiver

		let oldFactory <- self.account.storage[NFTFactory] <- create NFTFactory()
		destroy oldFactory

		self.account.storage[&NFTFactory] = &self.account.storage[NFTFactory] as &NFTFactory

        emit ContractInitialized()
	}
}

 
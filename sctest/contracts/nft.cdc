
pub contract interface NonFungibleToken {

    // The total number of tokens of this type in existance
    pub var totalSupply: UInt64

    // event ContractInitialized()
    // event Withdraw()
    // event Deposit()

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

		pub fun idExists(id: UInt64): Bool

        pub fun getMetaData(id: UInt64, field: String): String {
            pre {
                field.length != 0: "The requested field is undefined!"
			}
        }
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

        // idExists checks to see if a NFT with the given ID exists in the collection
        pub fun idExists(id: UInt64): Bool 

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64]

        pub fun getMetaData(id: UInt64, field: String): String
    }
}

pub contract CryptoKitties: NonFungibleToken {

    pub var totalSupply: UInt64

    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64

        pub var metadata: {String: String}

        init(initID: UInt64) {
            self.id = initID
            self.metadata = {}
        }
    }

    pub resource Collection: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.Metadata {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `UInt64` ID field
        pub var ownedNFTs: @{UInt64: NFT}

        init () {
            self.ownedNFTs <- {}
        }

        // withdraw removes an NFT from the collection and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            return <-token
        }

        pub fun batchWithdraw(ids: [UInt64]): @Collection {
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
        pub fun deposit(token: @NFT) {
            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            destroy oldToken
        }

        pub fun batchDeposit(tokens: @Collection) {
            var i = 0
            let keys = tokens.getIDs()

            while i < keys.length {
                self.deposit(token: <-tokens.withdraw(withdrawID: keys[i]))

                i = i + 1
            }
            destroy tokens
        }

        // idExists checks to see if a NFT with the given ID exists in the collection
        pub fun idExists(id: UInt64): Bool {
            return self.ownedNFTs[id] != nil
        }

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        pub fun getMetaData(id: UInt64, field: String): String {
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

    pub fun createNFT(id: UInt64): @NFT {
        return <- create NFT(initID: id)
    }

    pub fun createCollection(): @Collection {
        return <- create Collection()
    }

	pub resource NFTFactory {

		// the ID that is used to mint moments
		pub var idCount: UInt64

		init() {
			self.idCount = 1
		}

		// mintNFT mints a new NFT with a new ID
		// and deposit it in the recipients colelction using their collection reference
		pub fun mintNFT(recipient: &Collection) {

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

		self.account.storage[&Collection] = &self.account.storage[Collection] as Collection
        self.account.published[&NonFungibleToken.Receiver] = &self.account.storage[Collection] as NonFungibleToken.Receiver

		let oldFactory <- self.account.storage[NFTFactory] <- create NFTFactory()
		destroy oldFactory

		self.account.storage[&NFTFactory] = &self.account.storage[NFTFactory] as NFTFactory
	}
}


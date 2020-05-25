/*
    Description: Central Collection for a large number of topshot
                 NFTs

    authors: Joshua Hannan joshua.hannan@dapperlabs.com
             Bastian Muller bastian@dapperlabs.com

*/

import NonFungibleToken from 0x02
import TopShot from 0x03

pub contract TopShotShardedCollection {

    // Collection is a resource that every user who owns NFTs 
    // will store in their account to manage their NFTS
    //
    pub resource ShardedCollection: TopShot.MomentCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic { 
        // Dictionary of Moment conforming tokens
        // NFT is a resource type with a UInt64 ID field
        pub var collections: @{UInt64: NonFungibleToken.Collection}

        // the number of buckets to split moments into
        // this makes storage more efficient and performant
        pub let numBuckets: UInt64

        init(numBuckets: UInt64) {
            self.collections <- {}
            self.numBuckets = numBuckets

            // Create a new empty collection for each bucket
            var i: UInt64 = 0
            while i < numBuckets {

                self.collections[i] <-! TopShot.createEmptyCollection()

                i = i + UInt64(1)
            }
        }

        // withdraw removes an Moment from the collection and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            post {
                result.id == withdrawID: "The ID of the withdrawn NFT is incorrect"
            }
            // find the bucket it should be withdrawn from
            let bucket = withdrawID % self.numBuckets

            let token <- self.collections[bucket]?.withdraw(withdrawID: withdrawID)!
            
            return <-token
        }

        // batchWithdraw withdraws multiple tokens and returns them as a Collection
        pub fun batchWithdraw(ids: [UInt64]): @NonFungibleToken.Collection {
            var batchCollection <- TopShot.createEmptyCollection()
            
            // iterate through the ids and withdraw them from the collection
            for id in ids {
                batchCollection.deposit(token: <-self.withdraw(withdrawID: id))
            }
            return <-batchCollection
        }

        // deposit takes a Moment and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @NonFungibleToken.NFT) {

            // find the bucket this corresponds to
            let bucket = token.id % UInt64(self.numBuckets)

            self.collections[bucket]?.deposit(token: <-token)
        }

        // batchDeposit takes a Collection object as an argument
        // and deposits each contained NFT into this collection
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection) {
            let keys = tokens.getIDs()

            // iterate through the keys in the collection and deposit each one
            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }
            destroy tokens
        }

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64] {

            // concatenate IDs in all the collections
            var idArray: [UInt64] = []

            for key in self.collections.keys {
                idArray.concat(self.collections[key]?.getIDs() ?? [])
            }

            return idArray
        }

        // borrowNFT Returns a borrowed reference to a Moment in the collection
        // so that the caller can read data and call methods from it
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            post {
                result.id == id: "The ID of the reference is incorrect"
            }

            let bucket = id % self.numBuckets

            let ref = self.collections[bucket]?.borrowNFT(id: id)!

            // find NFT in the collections and borrow a reference
            return ref
        }

        // If a transaction destroys the Collection object,
        // All the NFTs contained within are also destroyed
        destroy() {
            destroy self.collections
        }
    }

    pub fun createEmptyCollection(numBuckets: UInt64): @ShardedCollection {
        return <-create ShardedCollection(numBuckets: numBuckets)
    }

    init() {

        // Put a new Collection in storage
        self.account.save<@ShardedCollection>(<- create ShardedCollection(numBuckets: 32), to: /storage/ShardedMomentCollection)

        // create a public capability for the collection
        self.account.link<&{TopShot.MomentCollectionPublic}>(/public/ShardedMomentCollection, target: /storage/ShardedMomentCollection)
    }
}
 
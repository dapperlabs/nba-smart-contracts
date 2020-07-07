/*
    Description: Central Collection for a large number of topshot
                 NFTs

    authors: Joshua Hannan joshua.hannan@dapperlabs.com
             Bastian Muller bastian@dapperlabs.com

    This resource object looks and acts exactly like a TopShot MomentCollection
    and (in a sense) shouldn’t have to exist! 
    The problem is that Cadence currently has a limitation where 
    storing more than ~100k objects in a single dictionary or array can fail. 
    Most MomentCollections are likely to be much, much smaller than this, 
    and that limitation will be removed in a future iteration of Cadence, 
    so most people will never need to worry about it.

    However! The main TopShot administration account DOES need to worry about it
    because it frequently needs to mint >10k Moments for sale, 
    and could easily end up needing to hold more than 100k Moments at one time.
    
    Until Cadence gets an update, that leaves in a bit of a pickle!

    This contract bundles together a bunch of MomentCollection objects 
    in a dictionary, and then distributes the individual Moments between them 
    while implementing the same public interface 
    as the default MomentCollection implementation. 
    If we assume that Moment IDs are uniformly distributed, 
    a ShardedCollection with 10 inner Collections should be able 
    to store 10x as many Moments (or ~1M).

    When Cadence is updated to allow larger dictionaries, 
    then this contract can be retired.

*/

import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS

pub contract TopShotShardedCollection {

    // ShardedCollection stores a dictionary of TopShot Collections
    // A moment is stored in the field that corresponds to its id % numBuckets
    pub resource ShardedCollection: TopShot.MomentCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic { 
        
        // Dictionary of topshot collections
        pub var collections: @{UInt64: TopShot.Collection}

        // the number of buckets to split moments into
        // this makes storage more efficient and performant
        pub let numBuckets: UInt64

        init(numBuckets: UInt64) {
            self.collections <- {}
            self.numBuckets = numBuckets

            // Create a new empty collection for each bucket
            var i: UInt64 = 0
            while i < numBuckets {

                self.collections[i] <-! TopShot.createEmptyCollection() as! @TopShot.Collection

                i = i + UInt64(1)
            }
        }

        // withdraw removes a Moment from one of the collections 
        // and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            post {
                result.id == withdrawID: "The ID of the withdrawn NFT is incorrect"
            }
            // find the bucket it should be withdrawn from
            let bucket = withdrawID % self.numBuckets

            // withdraw the moment
            let token <- self.collections[bucket]?.withdraw(withdrawID: withdrawID)!
            
            return <-token
        }

        // batchWithdraw withdraws multiple tokens and returns them as a Collection
        //
        // parameters: ids: an array fo the IDs to be withdrawn from the collection
        //
        // returns: @NonFungibleToken.Collection a collection containing the moments
        //          that were withdrawn
        pub fun batchWithdraw(ids: [UInt64]): @NonFungibleToken.Collection {
            var batchCollection <- TopShot.createEmptyCollection()
            
            // iterate through the ids and withdraw them from the collection
            for id in ids {
                batchCollection.deposit(token: <-self.withdraw(withdrawID: id))
            }
            return <-batchCollection
        }

        // deposit takes a Moment and adds it to the collections dictionary
        pub fun deposit(token: @NonFungibleToken.NFT) {

            // find the bucket this corresponds to
            let bucket = token.id % self.numBuckets

            // remove the collection
            let collection <- self.collections.remove(key: bucket)!

            // deposit the nft into the bucket
            collection.deposit(token: <-token)

            // put the collection back in storage
            self.collections[bucket] <-! collection
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

            var ids: [UInt64] = []
            // concatenate IDs in all the collections
            for key in self.collections.keys {
                for id in self.collections[key]?.getIDs() ?? [] {
                    ids.append(id)
                }
            }
            return ids
        }

        // borrowNFT Returns a borrowed reference to a Moment in the collection
        // so that the caller can read data and call methods from it
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            post {
                result.id == id: "The ID of the reference is incorrect"
            }

            // get the bucket of the nft to be borrowed
            let bucket = id % self.numBuckets

            // find NFT in the collections and borrow a reference
            return self.collections[bucket]?.borrowNFT(id: id)!
        }

        // borrowMoment Returns a borrowed reference to a Moment in the collection
        // so that the caller can read data and call methods from it
        // They can use this to read its setID, playID, serialNumber,
        // or any of the setData or Play Data associated with it by
        // getting the setID or playID and reading those fields from
        // the smart contract
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: A reference to the NFT
        pub fun borrowMoment(id: UInt64): &TopShot.NFT? {

            // get the bucket of the nft to be borrowed
            let bucket = id % self.numBuckets

            return self.collections[bucket]?.borrowMoment(id: id) ?? nil
        }

        // If a transaction destroys the Collection object,
        // All the NFTs contained within are also destroyed
        destroy() {
            destroy self.collections
        }
    }

    // function to create an empty ShardedCollection and return it to the caller
    pub fun createEmptyCollection(numBuckets: UInt64): @ShardedCollection {
        return <-create ShardedCollection(numBuckets: numBuckets)
    }
}
 
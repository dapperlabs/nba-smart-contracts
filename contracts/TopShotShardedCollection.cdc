/*
    Description: Central Collection for a large number of TopShot
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
    
    Until Cadence gets an update, that leaves us in a bit of a pickle!

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
import ViewResolver   from 0xVIEWRESOLVERADDRESS

access(all) contract TopShotShardedCollection {

    // ShardedCollection stores a dictionary of TopShot Collections
    // A Moment is stored in the field that corresponds to its id % numBuckets
    access(all) resource ShardedCollection: TopShot.MomentCollectionPublic, NonFungibleToken.CollectionPublic {
        
        // Dictionary of topshot collections
        access(contract) var collections: @{UInt64: TopShot.Collection}

        // The number of buckets to split Moments into
        // This makes storage more efficient and performant
        access(all) let numBuckets: UInt64

        init(numBuckets: UInt64) {
            self.collections <- {}
            self.numBuckets = numBuckets

            // Create a new empty collection for each bucket
            var i: UInt64 = 0
            while i < numBuckets {

                self.collections[i] <-! TopShot.createEmptyCollection(nftType: Type<@TopShot.NFT>()) as! @TopShot.Collection

                i = i + UInt64(1)
            }
        }

        // withdraw removes a Moment from one of the Collections 
        // and moves it to the caller
        access(NonFungibleToken.Withdraw | NonFungibleToken.Owner) fun withdraw(withdrawID: UInt64): @{NonFungibleToken.NFT} {
            post {
                result.id == withdrawID: "The ID of the withdrawn NFT is incorrect"
            }
            // Find the bucket it should be withdrawn from
            let bucket = withdrawID % self.numBuckets

            // Withdraw the moment
            let token <- self.collections[bucket]?.withdraw(withdrawID: withdrawID)!
            
            return <-token
        }

        // batchWithdraw withdraws multiple tokens and returns them as a Collection
        //
        // Parameters: ids: an array of the IDs to be withdrawn from the Collection
        //
        // Returns: @NonFungibleToken.Collection a Collection containing the moments
        //          that were withdrawn
        access(NonFungibleToken.Withdraw | NonFungibleToken.Owner) fun batchWithdraw(ids: [UInt64]): @{NonFungibleToken.Collection} {
            var batchCollection <- TopShot.createEmptyCollection(nftType: Type<@TopShot.NFT>())
            
            // Iterate through the ids and withdraw them from the Collection
            for id in ids {
                batchCollection.deposit(token: <-self.withdraw(withdrawID: id))
            }
            return <-batchCollection
        }

        // deposit takes a Moment and adds it to the Collections dictionary
        access(all) fun deposit(token: @{NonFungibleToken.NFT}) {

            // Find the bucket this corresponds to
            let bucket = token.id % self.numBuckets

            let collectionRef = (&self.collections[bucket] as &TopShot.Collection?)!

            // Deposit the nft into the bucket
            collectionRef.deposit(token: <-token)
        }

        // batchDeposit takes a Collection object as an argument
        // and deposits each contained NFT into this Collection
        access(all) fun batchDeposit(tokens: @{NonFungibleToken.Collection}) {
            let keys = tokens.getIDs()

            // Iterate through the keys in the Collection and deposit each one
            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }
            destroy tokens
        }

        // getIDs returns an array of the IDs that are in the Collection
        access(all) view fun getIDs(): [UInt64] {
            var ids: [UInt64] = []
            // Concatenate IDs in all the Collections
            for key in self.collections.keys {
                let collectionIDs = self.collections[key]?.getIDs() ?? []
                ids = ids.concat(collectionIDs)
            }
            return ids
        }

        // borrowNFT Returns a borrowed reference to a Moment in the Collection
        // so that the caller can read data and call methods from it
        access(all) view fun borrowNFT(_ id: UInt64): &{NonFungibleToken.NFT}? {

            // Get the bucket of the nft to be borrowed
            let bucket = id % self.numBuckets

            // Find NFT in the collections and borrow a reference
            return self.collections[bucket]?.borrowNFT(id)!
        }

        // borrowMoment Returns a borrowed reference to a Moment in the Collection
        // so that the caller can read data and call methods from it
        // They can use this to read its setID, playID, serialNumber,
        // or any of the setData or Play Data associated with it by
        // getting the setID or playID and reading those fields from
        // the smart contract
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: A reference to the NFT
        access(all) view fun borrowMoment(id: UInt64): &TopShot.NFT? {

            // Get the bucket of the nft to be borrowed
            let bucket = id % self.numBuckets

            return self.collections[bucket]?.borrowMoment(id: id)!
        }

        // Return a list of NFT types that this receiver accepts
        access(all) view fun getSupportedNFTTypes(): {Type: Bool} {
            let supportedTypes: {Type: Bool} = {}
            supportedTypes[Type<@TopShot.NFT>()] = true
            return supportedTypes
        }

        // Return whether or not the given type is accepted by the collection
        // A collection that can accept any type should just return true by default
        access(all) view fun isSupportedNFTType(type: Type): Bool {
            if type == Type<@TopShot.NFT>() {
                return true
            }
            return false
        }

        // Return the amount of NFTs stored in the collection
        access(all) view fun getLength(): Int {
            return self.getIDs().length
        }

    }

    // Creates an empty ShardedCollection and returns it to the caller
    access(all) fun createEmptyCollection(numBuckets: UInt64): @ShardedCollection {
        return <-create ShardedCollection(numBuckets: numBuckets)
    }
}
 
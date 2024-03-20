import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This transaction creates and stores an empty moment collection 
// and creates a public capability for it.
// Moments are split into a number of buckets
// This makes storage more efficient and performant

// Parameters
//
// numBuckets: The number of buckets to split Moments into

transaction(numBuckets: UInt64) {

    prepare(acct: auth(Storage, Capabilities) &Account) {

        if acct.storage.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection) == nil {

            let collection <- TopShotShardedCollection.createEmptyCollection(numBuckets: numBuckets)
            // Put a new Collection in storage
            acct.storage.save(<-collection, to: /storage/ShardedMomentCollection)

            acct.capabilities.unpublish(/public/MomentCollection)
            acct.capabilities.publish(
                acct.capabilities.storage.issue<&TopShotShardedCollection.ShardedCollection>(/storage/ShardedMomentCollection),
                at: /public/MomentCollection
            )
        } else {

            panic("Sharded Collection already exists!")
        }
    }
}
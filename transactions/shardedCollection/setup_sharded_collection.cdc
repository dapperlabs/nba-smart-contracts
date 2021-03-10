import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

// This transaction creates and stores an empty moment collection 
// and creates a public capability for it.
// Moments are split into a number of buckets
// This makes storage more efficient and performant

// Parameters
//
// numBuckets: The number of buckets to split Moments into

transaction(numBuckets: UInt64) {

    prepare(acct: AuthAccount) {

        if acct.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection) == nil {

            let collection <- TopShotShardedCollection.createEmptyCollection(numBuckets: numBuckets)
            // Put a new Collection in storage
            acct.save(<-collection, to: /storage/ShardedMomentCollection)

            // create a public capability for the collection
            if acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/ShardedMomentCollection) == nil {
                acct.unlink(/public/MomentCollection)
            }

            acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/ShardedMomentCollection)
        } else {

            panic("Sharded Collection already exists!")
        }
    }
}
/*

  AdminReceiver.cdc

  This contract defines a function that takes a TopShot admin
  object and stores it in the storage of the contract account
  so it can be used normally

 */

import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

pub contract TopshotAdminReceiver {

    pub fun storeAdmin(newAdmin: @TopShot.Admin) {
        self.account.save(<-newAdmin, to: /storage/TopShotAdmin)
    }
    
    init() {
        if self.account.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection) == nil {
            let collection <- TopShotShardedCollection.createEmptyCollection(numBuckets: 32)
            // Put a new Collection in storage
            self.account.save(<-collection, to: /storage/ShardedMomentCollection)

            self.account.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/ShardedMomentCollection)
        }
    }
}

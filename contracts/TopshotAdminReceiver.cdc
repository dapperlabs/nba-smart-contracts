/*

  AdminReceiver.cdc

  This contract defines a function that takes a TopShot Admin
  object and stores it in the storage of the contract account
  so it can be used.

 */

import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

access(all) contract TopshotAdminReceiver {

    // storeAdmin takes a TopShot Admin resource and 
    // saves it to the account storage of the account
    // where the contract is deployed
    access(all) fun storeAdmin(newAdmin: @TopShot.Admin) {
        self.account.storage.save(<-newAdmin, to: /storage/TopShotAdmin)
    }
    
    init() {
        // Save a copy of the sharded Moment Collection to the account storage
        if self.account.storage.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection) == nil {
            let collection <- TopShotShardedCollection.createEmptyCollection(numBuckets: 32)
            // Put a new Collection in storage
            self.account.storage.save(<-collection, to: /storage/ShardedMomentCollection)
            let cap = self.account.capabilities.storage.issue<&TopShotShardedCollection.ShardedCollection>(/storage/ShardedMomentCollection)
            self.account.capabilities.publish(cap, at: /public/MomentCollection)
        }
    }
}

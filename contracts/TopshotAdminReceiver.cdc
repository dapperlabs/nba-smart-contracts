/*

  AdminReceiver.cdc

  This contract defines a function that takes a TopShot admin
  object and stores it in the storage of the contract account
  so it can be used normally

 */

import TopShot from 0x03

pub contract TopshotAdminReceiver {

    pub fun storeAdmin(newAdmin: @TopShot.Admin) {
        self.account.save(<-newAdmin, to: /storage/TopShotAdmin)
    }
    
    init() {
            let collection <- TopShot.createEmptyCollection() as! @TopShot.Collection
            // Put a new Collection in storage
            self.account.save(<-collection, to: /storage/MomentCollection)

            // create a public capability for the collection
            self.account.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/MomentCollection)
    }
}

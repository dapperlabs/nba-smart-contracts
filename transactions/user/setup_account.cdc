import TopShot from 0xTOPSHOTADDRESS

transaction {

    prepare(acct: AuthAccount) {

        if acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection) == nil {
            let collection <- TopShot.createEmptyCollection() as! @TopShot.Collection
            // Put a new Collection in storage
            acct.save(<-collection, to: /storage/MomentCollection)

            // create a public capability for the collection
            acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/MomentCollection)
        }
    }
}
import TopShotLocking from 0xTOPSHOTADDRESS

transaction(ids: [UInt64]) {

    prepare(acct: auth(Storage) &Account) {
        let adminRef = acct.storage.borrow<&TopShotLocking.Admin>(from: /storage/TopShotLockingAdmin)
    ?? panic("Could not borrow a reference to the Admin resource")
    
        for id in ids {
            adminRef.unlockByID(id: UInt64(id))
        }
    }

}


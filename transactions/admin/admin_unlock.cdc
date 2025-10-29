import TopShotLocking from 0x0b2a3299cc857e29

transaction(ids: [UInt64]) {

    prepare(acct: auth(Storage) &Account) {
        let adminRef = acct.storage.borrow<&TopShotLocking.Admin>(from: /storage/TopShotLockingAdmin)
    ?? panic("Could not borrow a reference to the Admin resource")
        // let adminRef = adminCap.borrow()!
        for id in ids {
            adminRef.unlockByID(id: UInt64(id))
        }
    }

}


import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

transaction(ids: [UInt64], expiryTimestamp: UFix64) {
    let adminRef: &TopShotLocking.Admin

    prepare(acct: auth(BorrowValue) &Account) {
        // Set TopShotLocking admin ref
        self.adminRef = acct.storage.borrow<&TopShotLocking.Admin>(from: /storage/TopShotLockingAdmin)
            ?? panic("Could not find reference to TopShotLocking Admin resource")
    }

    execute {
        for id in ids {
            self.adminRef.setLockExpiryByID(id: id, expiryTimestamp: expiryTimestamp)
        }
    }
}

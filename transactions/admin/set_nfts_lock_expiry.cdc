import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

transaction(id: UInt64, expiryTimestamp: UFix64) {
    let adminRef: &TopShotLocking.Admin

    prepare(acct: AuthAccount) {
        // Set TopShotLocking admin ref
        self.adminRef = acct.borrow<&TopShotLocking.Admin>(from: /storage/TopShotLockingAdmin)
            ?? panic("Could not find reference to TopShotLocking Admin resource")
    }

    execute {
        self.adminRef.setLockExpiryByID(id: id, expiryTimestamp: expiryTimestamp)
    }
}

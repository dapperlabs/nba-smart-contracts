import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

transaction() {
    let adminRef: &TopShotLocking.Admin

    prepare(acct: AuthAccount) {
        // Set TopShotLocking admin ref
        self.adminRef = acct.borrow<&TopShotLocking.Admin>(from: /storage/TopShotLockingAdmin)
            ?? panic("Could not find reference to TopShotLocking Admin resource")
    }

    execute {
        self.adminRef.unlockAll()
    }
}

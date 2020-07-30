import TopShot from 0xTOPSHOTADDRESS

// This transaction is for retiring all plays from a set, which
// makes it so that moments can no longer be minted
// from all the editions with that set

transaction(setID: UInt32) {
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        // borrow a reference to the specified set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // retire all the plays
        setRef.retireAll()
    }
}
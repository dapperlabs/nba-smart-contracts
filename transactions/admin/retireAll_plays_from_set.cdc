import TopShot from 0xTOPSHOTADDRESS

// This transaction is for retiring all plays from a set, which
// makes it so that moments can no longer be minted
// from all the editions with that set

// Parameters:
//
// setID: the ID of the set to be retired entirely

transaction(setID: UInt32) {
    let adminRef: &TopShot.Admin

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.storage.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        // borrow a reference to the specified set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // retire all the plays
        setRef.retireAll()
    }
}
import TopShot from 0xTOPSHOTADDRESS

// This transaction locks a set so that new plays can no longer be added to it

// Parameters:
//
// setID: the ID of the set to be locked

transaction(setID: UInt32) {

    // local variable for the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: auth(BorrowValue) &Account) {
        // borrow a reference to the admin resource
        self.adminRef = acct.storage.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        // borrow a reference to the Set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // lock the set permanently
        setRef.lock()
    }

    post {
        
        TopShot.isSetLocked(setID: setID)!:
            "Set did not lock"
    }
}
import TopShot from 0xTOPSHOTADDRESS

// This transaction locks a set so that new plays can no longer be added to it

transaction(setID: UInt32) {

    // local variable for the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // borrow a reference to the admin resource
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        // borrow a reference to the Set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // lock the set permanently
        setRef.lock()
    }
}
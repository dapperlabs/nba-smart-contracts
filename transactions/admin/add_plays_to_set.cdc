import TopShot from 0xTOPSHOTADDRESS

// This transaction adds multiple plays to a set
		
transaction(setID: UInt32, plays: [UInt32]) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute() {

        // borrow a reference to the set to be added to
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Add the specified play IDs
        setRef.addPlays(playIDs: plays)
    }
}
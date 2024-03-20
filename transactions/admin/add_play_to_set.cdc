import TopShot from 0xTOPSHOTADDRESS

// This transaction is how a Top Shot admin adds a created play to a set

// Parameters:
//
// setID: the ID of the set to which a created play is added
// playID: the ID of the play being added

transaction(setID: UInt32, playID: UInt32) {

    // Local variable for the topshot Admin object
    let adminRef: auth(TopShot.NFTMinter) &TopShot.Admin

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.storage.borrow<auth(TopShot.NFTMinter) &TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
    }

    execute {
        
        // Borrow a reference to the set to be added to
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Add the specified play ID
        setRef.addPlay(playID: playID)
    }

    post {
        TopShot.getPlaysInSet(setID: setID)!.contains(playID): 
            "set does not contain playID"
    }
}
import TopShot from 0xTOPSHOTADDRESS

// This is a transaction an admin would use to retire all the plays in a set
// which makes it so that no more moments can be minted from the retired plays

// Parameters:
//
// setID: the ID of the set to be retired entirely

transaction(setID: UInt32) {

    // local variable for the admin reference
    let adminRef: auth(TopShot.NFTMinter) &TopShot.Admin

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the admin resource
        self.adminRef = acct.storage.borrow<auth(TopShot.NFTMinter) &TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {

        // borrow a reference to the specified set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // retire all the plays permenantely
        setRef.retireAll()
    }
}
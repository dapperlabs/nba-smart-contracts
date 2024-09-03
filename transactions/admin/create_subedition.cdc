import TopShot from 0xTOPSHOTADDRESS

// This transaction creates a new subedition struct
// and stores it in the Top Shot smart contract

// Parameters:
//
// name:  the name of a new Subedition to be created
// metadata: A dictionary of all the play metadata associated

transaction(name:String, metadata:{String:String}) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin
    let currSubeditionID: UInt32

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the admin resource
        self.currSubeditionID = TopShot.getNextSubeditionID();
        self.adminRef = acct.storage.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {

        // Create a subedition with the specified metadata
        self.adminRef.createSubedition(name: name, metadata: metadata)
    }

    post {

        TopShot.getSubeditionByID(subeditionID: self.currSubeditionID) != nil:
            "SubedititonID doesnt exist"
    }
}
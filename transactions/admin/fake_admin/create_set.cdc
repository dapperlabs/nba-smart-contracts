import TopShot from 0xFAKETOPSHOTADDRESS

// This transaction is for the admin to create a new set resource
// and store it in the top shot smart contract

// Parameters:
//
// setName: the name of a new Set to be created

transaction(setName: String) {
    
    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin
    let currSetID: UInt32

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
        self.currSetID = TopShot.nextSetID;
    }

    execute {
        
        // Create a set with the specified name
        self.adminRef.createSet(name: setName)
    }

    post {
        
        TopShot.getSetName(setID: self.currSetID) == setName:
          "Could not find the specified set"
    }
}

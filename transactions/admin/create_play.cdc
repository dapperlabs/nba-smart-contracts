import TopShot from 0xTOPSHOTADDRESS

// This transaction creates a new play struct 
// and stores it in the Top Shot smart contract
// We currently stringify the metadata and insert it into the 
// transaction string, but want to use transaction arguments soon

// Parameters:
//
// metadata: A dictionary of all the play metadata associated

transaction(metadata: {String: String}) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin
    let currPlayID: UInt32

    prepare(acct: AuthAccount) {

        // borrow a reference to the admin resource
        self.currPlayID = TopShot.nextPlayID;
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {

        // Create a play with the specified metadata
        self.adminRef.createPlay(metadata: metadata)
    }

    post {
        
        TopShot.getPlayMetaData(playID: self.currPlayID) != nil:
            "playID doesnt exist"
    }
}
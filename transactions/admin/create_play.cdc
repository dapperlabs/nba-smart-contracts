import TopShot from 0xTOPSHOTADDRESS

// This transaction creates a new play struct 
// and stores it in the Top Shot smart contract
// We currently stringify the metadata and instert it into the 
// transaction string, but want to use transaction arguments soon

transaction(metadata: {String: String}) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the admin resource
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {

        // Create a play with the specified metadata
        self.adminRef.createPlay(metadata: metadata)
    }
}
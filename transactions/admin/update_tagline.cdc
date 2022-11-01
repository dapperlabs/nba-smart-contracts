import TopShot from 0xTOPSHOTADDRESS

// This transaction updates an existing play struct 
// and stores it in the Top Shot smart contract
// We currently stringify the metadata and insert it into the 
// transaction string, but want to use transaction arguments soon

// Parameters:
//
// metadata: A dictionary of all the play metadata associated

transaction(plays: {UInt32: String}) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the admin resource
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        // update a play with the specified metadata
        for key in plays.keys {
            self.adminRef.updatePlayTagline(playID: key, tagline: plays[key] ?? panic("No tagline for play"))
        }
    }

    post {
        for key in plays.keys {
            TopShot.getPlayMetaDataByField(playID: key, field: "tagline") != nil:
                "tagline doesnt exist"
        }
    }
}

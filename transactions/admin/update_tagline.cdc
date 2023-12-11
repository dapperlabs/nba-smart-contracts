import TopShot from 0xTOPSHOTADDRESS

// This transaction updates multiple existing plays' taglines
// and stores them in the Top Shot smart contract

// Parameters:
//
// plays: A dictionary of {playID: tagline} pairs

transaction(plays: {UInt32: String}) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin
    let firstKey: UInt32
    let lastKey: UInt32

    prepare(acct: AuthAccount) {

        // borrow a reference to the admin resource
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
        self.firstKey = plays.keys[0]
        self.lastKey = plays.keys[plays.keys.length - 1]
    }

    execute {
        // update multiple plays with the specified metadata
        for key in plays.keys {
            self.adminRef.updatePlayTagline(playID: key, tagline: plays[key] ?? panic("No tagline for play"))
        }
    }

    post {
        TopShot.getPlayMetaDataByField(playID: self.firstKey, field: "Tagline") != nil:
            "First play's tagline does not exist"
        TopShot.getPlayMetaDataByField(playID: self.lastKey, field: "Tagline") != nil:
            "Last play's tagline does not exist"
    }
}

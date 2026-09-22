import TopShot from 0xTOPSHOTADDRESS

// This transaction sets one or more metadata fields on an existing play.
// Fields not named in the argument are left as they are.

// Parameters:
//
// playID: The ID of the play to update
// metadata: A dictionary of {field: value} pairs to set on the play

transaction(playID: UInt32, metadata: {String: String}) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin
    let firstKey: String

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the admin resource
        self.adminRef = acct.storage.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
        self.firstKey = metadata.keys[0]
    }

    execute {
        for key in metadata.keys {
            self.adminRef.updatePlayMetadata(playID: playID, key: key, value: metadata[key] ?? panic("No value for metadata field"))
        }
    }

    post {
        TopShot.getPlayMetaDataByField(playID: playID, field: self.firstKey) == metadata[self.firstKey]:
            "Play metadata field was not updated"
    }
}

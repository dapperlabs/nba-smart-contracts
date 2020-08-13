import TopShot from 0xTOPSHOTADDRESS

// This transaction creates a new play struct 
// and stores it in the Top Shot smart contract
// We currently stringify the metadata and instert it into the 
// transaction string, but want to use transaction arguments soon

transaction() {
    prepare(acct: AuthAccount) {

        // borrow a reference to the admin resource
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")

        // Create a play with specified metadata
        // The argument is a string template field, so if you are running this manually,
        // you can replace it with a {String: String} mapping
        // Example: {"Name": "TJ Warren", "Position": "Superstar"}
        admin.createPlay(metadata: %s)
    }
}
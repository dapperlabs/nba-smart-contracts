import TopShot from 0xTOPSHOTADDRESS

// This transaction is for the admin to create a new showcase resource
// and store it in the top shot smart contract

transaction() {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
    }

    execute {
        self.adminRef.createSubEditionResource()
    }
}
import TopShot from 0xTOPSHOTADDRESS

// This transaction is for the admin to create a new set resource
// and store it in the top shot smart contract

transaction(setName: String) {
    prepare(acct: AuthAccount) {
        // borrow a reference to the Admin resource in storage
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")

        // Create a set with the specified name
        admin.createSet(name: setName)
    }
}
import TopShot from 0x03

// This transaction allows an admin to create a play struct
// with metadata

transaction {

    // temporary reference to the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // create a reference to the Admin resource
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        
        // Create two new plays
        let id1 = self.adminRef.createPlay(metadata: {"Name": "Lebron"})
        let id2 = self.adminRef.createPlay(metadata: {"Name": "Oladipo"})

        log("PlayData 1 and 2 Succcesfully created!")
    }
}
 
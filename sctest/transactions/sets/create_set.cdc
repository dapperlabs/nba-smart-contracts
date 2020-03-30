import TopShot from 0x03

// Transaction to create a new Set and add it to the contract

transaction {

    // temporary reference to the Admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // create a reference to Admin
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        // create a new set with a name
        self.adminRef.createSet(name: "Genesis")
    }
}
 
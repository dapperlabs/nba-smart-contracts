import TopShot from 0x03

// This transaction locks a set so that no more plays can be added to it

transaction {

    // temporary reference to the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // create a reference to the private admin resource object
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        // create a reference to a specific set
        let setRef = self.adminRef.borrowSet(setID: 0)

        // lock the set
        setRef.lock()
    }
}
 
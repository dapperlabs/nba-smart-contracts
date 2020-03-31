import TopShot from 0x03

// This transaction locks a set so that no more plays can be added to it

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // create a temporary reference to the private admin resource
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        // create a reference to a specific set
        let setRef = self.adminRef.borrowSet(setID: 1)

        // lock the set
        setRef.lock()
    }
}
 
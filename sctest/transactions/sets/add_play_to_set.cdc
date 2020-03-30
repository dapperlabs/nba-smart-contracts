import TopShot from 0x03

// transaction for an admin to add a new play to a set
// to create an edition

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        // get a reference to the private set resource
        let setRef = self.adminRef.borrowSet(setID: 0)

        // add a play to the set using its reference
        setRef.addPlay(playID: 1)
    }
}
 
import TopShot from 0x03

// This transaction retires all the plays from a single set
// so that no more moments can be minted for those plays anymore

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // create a temporary admin reference
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        // create a reference to a specific set
        let setRef = self.adminRef.borrowSet(setID: 1)

        // retire all the plays
        setRef.retireAll()
    }
}
 
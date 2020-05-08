import TopShot from 0x03

// This transaction retires a single play from a set
// so that moments cannot be minted for it anymore

transaction {

    // temporary reference to the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // create a reference to the private admin resource object
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {
        // get a reference to a specific set
        let setRef = self.adminRef.borrowSet(setID: 0)

        // retire a play from the set
        setRef.retirePlay(playID: 0)
    }
}
 
import TopShot from 0x03

// This transaction allows an admin to increment the series number
// of the TopShot smart contract

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // get a reference to the stored admin resource
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {
        
        // Call the startNewSeries function with the Admin reference
        let newSeriesNumber = self.adminRef.startNewSeries()

    }
}
 
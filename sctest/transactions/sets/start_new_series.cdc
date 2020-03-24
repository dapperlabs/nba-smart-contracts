import TopShot from 0x03

transaction {

    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {
        
        let newSeriesNumber = self.adminRef.startNewSeries()

    }
}
 
import TopShot from 0xTOPSHOTADDRESS

// This transaction is for an Admin to start a new Top Shot series

transaction {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin
    let currentSeries: UInt32

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")

        self.currentSeries = TopShot.currentSeries
    }

    execute {
        
        // Increment the series number
        self.adminRef.startNewSeries()
    }

    post {
    
        TopShot.currentSeries == self.currentSeries + 1 as UInt32:
            "new series not started"
    }
}
 
import TopShot from 0xTOPSHOTADDRESS

// This transaction is for an Admin to start a new Top Shot series

transaction {
    prepare(acct: AuthAccount) {
        // borrow a reference to the Admin resource in storage
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")

        // Increment the series number
        admin.startNewSeries()
    }
}
 
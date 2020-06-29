import TopShot from 0xTOPSHOTADDRESS

transaction {
    prepare(acct: AuthAccount) {
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
        admin.startNewSeries()
    }
}
 
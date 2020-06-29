import TopShot from 0xTOPSHOTADDRESS

transaction(setName: String) {
    prepare(acct: AuthAccount) {
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
        admin.createSet(name: setName)
    }
}
 
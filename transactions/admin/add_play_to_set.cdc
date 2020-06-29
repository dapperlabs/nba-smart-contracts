import TopShot from 0xTOPSHOTADDRESS

transaction(setID: UInt32, playID: UInt32) {

    prepare(acct: AuthAccount) {
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
        let setRef = admin.borrowSet(setID: setID)
        setRef.addPlay(playID: playID)
    }
}
 
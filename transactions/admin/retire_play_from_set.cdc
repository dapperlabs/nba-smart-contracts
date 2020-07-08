import TopShot from 0xTOPSHOTADDRESS

transaction(setID: UInt32, playID: UInt32) {
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        let setRef = self.adminRef.borrowSet(setID: setID)

        setRef.retirePlay(playID: playID)
    }
}
 
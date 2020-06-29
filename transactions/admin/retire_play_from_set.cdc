import TopShot from 0x%s

transaction {
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        let setRef = self.adminRef.borrowSet(setID: %d)

        setRef.retirePlay(playID: UInt32(%d))
    }
}
 
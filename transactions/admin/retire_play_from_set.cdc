import TopShot from 0xTOPSHOTADDRESS

// This transaction is for retiring a play from a set, which
// makes it so that moments can no longer be minted from that edition

transaction(setID: UInt32, playID: UInt32) {
    
    // local variable for storing the reference to the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {

        // borrow a reference to the specified set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // retire the play
        setRef.retirePlay(playID: playID)
    }
}
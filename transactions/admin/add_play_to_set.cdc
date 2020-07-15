import TopShot from 0xTOPSHOTADDRESS

// This transaction is how a Top Shot admin adds a created play to a set

transaction(setID: UInt32, playID: UInt32) {

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
        
        // Borrow a reference to the set to be added to
        let setRef = admin.borrowSet(setID: setID)

        // Add the specified play ID
        setRef.addPlay(playID: playID)
    }
}
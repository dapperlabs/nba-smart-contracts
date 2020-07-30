import TopShot from 0xTOPSHOTADDRESS

// This transaction is what an admin would use to mint a single new moment
// and deposit it in a user's collection

transaction(setID: UInt32, playID: UInt32, recipientAddr: Address) {
    // local variable for the admin reference
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {
        // Borrow a reference to the specified set
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Mint a new NFT
        let moment1 <- setRef.mintMoment(playID: playID)

        // get the public account object for the recipient
        let recipient = getAccount(recipientAddr)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Cannot borrow a reference to the recipient's moment collection")

        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-moment1)
    }
}
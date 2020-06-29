import TopShot from 0xTOPSHOTADDRESS

transaction(setID: UInt32, playID: UInt32, recipientAddr: Address) {
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Mint a new NFT
        let moment1 <- setRef.mintMoment(playID: playID)
        let recipient = getAccount(recipientAddr)
        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!
        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-moment1)
    }
}
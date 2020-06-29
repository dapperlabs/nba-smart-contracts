import TopShot from 0xTOPSHOTADDRESS

transaction(setID: UInt32, playID: UInt32, quantity: UInt64, recipientAddr: Address) {
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Mint a new NFT
        let collection <- setRef.batchMintMoment(playID: playID, quantity: quantity)
        let recipient = getAccount(recipientAddr)
        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!
        // deposit the NFT in the receivers collection
        receiverRef.batchDeposit(tokens: <-collection)
    }
}
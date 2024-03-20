import TopShot from 0xTOPSHOTADDRESS

// This transaction mints multiple moments
// from a single set/play/subedition combination

// Parameters:
//
// setID: the ID of the set to be minted from
// playID: the ID of the Play from which the Moments are minted
// subeditionID: the ID of play's subedition
// quantity: the quantity of Moments to be minted
// recipientAddr: the Flow address of the account receiving the collection of minted moments

transaction(setID: UInt32, playID: UInt32, quantity: UInt64, subeditionID: UInt32, recipientAddr: Address) {

    // Local variable for the topshot Admin object
    let adminRef: auth(TopShot.NFTMinter) &TopShot.Admin

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.storage.borrow<auth(TopShot.NFTMinter) &TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {

        // borrow a reference to the set to be minted from
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Mint all the new NFTs with Subeditions
        let collection <- setRef.batchMintMomentWithSubedition(playID: playID, quantity: quantity, subeditionID: subeditionID)

        // Get the account object for the recipient of the minted tokens
        let recipient = getAccount(recipientAddr)

        // get the Collection reference for the receiver
        let receiverRef = recipient.capabilities.borrow<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection)
            ?? panic("Cannot borrow a reference to the recipient's collection")

        // deposit the NFT in the receivers collection
        receiverRef.batchDeposit(tokens: <-collection)
    }
}
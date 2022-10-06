import TopShot from 0xTOPSHOTADDRESS

// This transaction mints multiple moments
// from a single set/play/subEdition combination

// Parameters:
//
// setID: the ID of the set to be minted from
// playID: the ID of the Play from which the Moments are minted
// subEditionID: the ID of play's subEdition
// quantity: the quantity of Moments to be minted
// recipientAddr: the Flow address of the account receiving the collection of minted moments

transaction(setID: UInt32, playID: UInt32, quantity: UInt64, subEditionID: UInt32, recipientAddr: Address) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
    }

    execute {

        // borrow a reference to the set to be minted from
        let setRef = self.adminRef.borrowSet(setID: setID)

        // Mint all the new NFTs with SubEditions
        let collection <- setRef.batchMintMomentWithSubEdition(playID: playID, quantity: quantity, subEditionID: subEditionID)

        // Get the account object for the recipient of the minted tokens
        let recipient = getAccount(recipientAddr)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Cannot borrow a reference to the recipient's collection")

        // deposit the NFT in the receivers collection
        receiverRef.batchDeposit(tokens: <-collection)
    }
}
import TopShot from 0xTOPSHOTADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This transaction unlocks a list of TopShot NFTs

// Parameters
//
// ids: array of TopShot moment Flow IDs

transaction(ids: [UInt64]) {
    prepare(acct: auth(BorrowValue) &Account) {
        let collectionRef = acct.storage.borrow<auth(NonFungibleToken.Update) &TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.batchUnlock(ids: ids)
    }
}

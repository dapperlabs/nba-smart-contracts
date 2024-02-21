import TopShot from 0xTOPSHOTADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This transaction unlocks a TopShot NFT removing it from the locked dictionary
// and re-enabling the ability to withdraw, sell, and transfer the moment

// Parameters
//
// id: the Flow ID of the TopShot moment
transaction(id: UInt64) {
    prepare(acct: auth(BorrowValue) &Account) {
        let collectionRef = acct.storage.borrow<auth(NonFungibleToken.Update) &TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.unlock(id: id)
    }
}

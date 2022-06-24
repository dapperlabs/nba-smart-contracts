import TopShot from 0xTOPSHOTADDRESS

// This transaction unlocks a TopShot NFT removing it from the locked dictionary
// and re-enabling the ability to withdraw, sell, and transfer the moment

// Parameters
//
// id: the Flow ID of the TopShot moment
transaction(id: UInt64) {
    prepare(acct: AuthAccount) {
        let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.unlock(id: id)
    }
}

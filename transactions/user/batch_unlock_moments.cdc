import TopShot from 0xTOPSHOTADDRESS

// This transaction unlocks a list of TopShot NFTs

// Parameters
//
// ids: array of TopShot moment Flow IDs

transaction(ids: [UInt64]) {
    prepare(acct: AuthAccount) {
        let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.batchUnlock(ids: ids)
    }
}

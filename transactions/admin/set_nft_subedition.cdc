import TopShot from 0xTOPSHOTADDRESS

// This transaction links nft to subedititon

// Parameters:
//
// nftID:  the unique ID of nft
// subeditionID: the unique ID of subedition

transaction(nftID: UInt64, subeditionID: UInt32) {

    // Local variable for the topshot Admin object
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // borrow a reference to the admin resource
        self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
            ?? panic("No admin resource in storage")
    }

    execute {
        // Create a subedition with the specified metadata
        self.adminRef.setMomentsSubedition(nftID: nftID, subeditionID: subeditionID)
    }
}
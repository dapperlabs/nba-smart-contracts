import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

// This script determines if a moment is locked

// Parameters:
//
// account: The Flow Address of the account who owns the moment
// id: The unique ID for the moment

// Returns: Bool
// Whether the moment is locked

pub fun main(account: Address, id: UInt64): Bool {

    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    let nftRef = collectionRef.borrowNFT(id: id)

    return TopShotLocking.isLocked(nftRef: nftRef)
}

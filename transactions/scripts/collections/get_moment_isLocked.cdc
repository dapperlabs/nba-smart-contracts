import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

// This script determines if a moment is locked

// Parameters:
//
// account: The Flow Address of the account who owns the moment
// id: The unique ID for the moment

// Returns: Bool
// Whether the moment is locked

access(all) fun main(account: Address, id: UInt64): Bool {

    let collectionRef = getAccount(account).capabilities.borrow<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection)
        ?? panic("Could not get public moment collection reference")

    let nftRef = collectionRef.borrowNFT(id)!

    return TopShotLocking.isLocked(nftRef: nftRef)
}

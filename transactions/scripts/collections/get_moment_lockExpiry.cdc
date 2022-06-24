import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

// This script gets the time at which a moment will be eligible for unlocking

// Parameters:
//
// account: The Flow Address of the account who owns the moment
// id: The unique ID for the moment

// Returns: UFix64
// The unix timestamp when the moment is unlockable

pub fun main(account: Address, id: UInt64): UFix64 {

    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    let nftRef = collectionRef.borrowNFT(id: id)

    return TopShotLocking.getLockExpiry(nftRef: nftRef)
}

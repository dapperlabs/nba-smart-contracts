import TopShot from 0xTOPSHOTADDRESS

// This script returns true if a moment with the specified ID
// exists in a user's collection

// Parameters:
//
// account: The Flow Address of the account whose moment data needs to be read
// id: The unique ID for the moment whose data needs to be read

// Returns: Bool
// Whether a moment with specified ID exists in user's collection

pub fun main(account: Address, id: UInt64): Bool {

    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    return collectionRef.borrowNFT(id: id) != nil
}
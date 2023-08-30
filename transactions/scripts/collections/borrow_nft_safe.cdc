import TopShot from 0xTOPSHOTADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This is a script to get a boolean value safely to see if a moment exists in a collection
// We expect this will not panic if the NFT is not in the collection
// Change the `account` to whatever account you want
// and as long as they have a published Collection receiver, you can 
// get reference to the NFTs they own.

// Parameters:
//
// account: The Flow Address of the account whose moment data needs to be read
// nftID: The ID of the NFT to return

// Returns: Boolean value indicating if the NFT is in the collection

pub fun main(account: Address, nftID: UInt64 ): Bool {

    let acct = getAccount(account)

    let collectionRef = acct.getCapability(/public/MomentCollection)
                            .borrow<&{NonFungibleToken.CollectionPublic}>()!

    let optionalNFT = collectionRef.borrowNFTSafe(id: nftID)

    // optional binding
    if let nft = optionalNFT {
        return true
    } else {
        return false
    }
}
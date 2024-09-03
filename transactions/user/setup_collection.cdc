import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

// This transaction sets up an account to use Top Shot
// by storing an empty moment collection and creating
// a public capability for it

transaction {

    prepare(acct: auth(Storage, Capabilities) &Account) {

        // First, check to see if a moment collection already exists
        if acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection) == nil {
            // create a new TopShot Collection
            let collection <- TopShot.createEmptyCollection(nftType: Type<@TopShot.NFT>()) as! @TopShot.Collection
            // Put the new Collection in storage
            acct.storage.save(<-collection, to: /storage/MomentCollection)
        }

        acct.capabilities.unpublish(/public/MomentCollection)
        acct.capabilities.publish(
            acct.capabilities.storage.issue<&TopShot.Collection>(/storage/MomentCollection),
            at: /public/MomentCollection
        )
    }
}

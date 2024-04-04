import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS
import PackNFT from 0xPACKNFTADDRESS

// This transaction sets up an account to use Top Shot
// by storing an empty moment collection and creating
// a public capability for it

transaction {

    prepare(acct: auth(Storage, Capabilities) &Account) {

        // First, check to see if a moment collection already exists
        if acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection) == nil {
            // create a new TopShot Collection
            let collection <- TopShot.createEmptyCollection() as! @TopShot.Collection
            // Put the new Collection in storage
            acct.storage.save(<-collection, to: /storage/MomentCollection)
        }

        acct.capabilities.unpublish(/public/MomentCollection)
        acct.capabilities.publish(
            acct.capabilities.storage.issue<&TopShot.Collection>(/storage/MomentCollection),
            at: /public/MomentCollection
        )

        // Create a PackNFT collection in the signer account if it doesn't already have one
        if acct.storage.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath) == nil {
            acct.storage.save(<- PackNFT.createEmptyCollection(), to: PackNFT.CollectionStoragePath);
        }

        // Create collection public capability if it doesn't already exist
        acct.capabilities.unpublish(PackNFT.CollectionPublicPath)
        acct.capabilities.publish(
            acct.capabilities.storage.issue<&PackNFT.Collection>(PackNFT.CollectionStoragePath),
            at: PackNFT.CollectionPublicPath
        )
    }
}

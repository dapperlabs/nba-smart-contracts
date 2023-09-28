import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS
import PackNFT from 0xPACKNFTADDRESS

// This transaction sets up an account to use Top Shot
// by storing an empty moment collection and creating
// a public capability for it

transaction {

    prepare(acct: AuthAccount) {

        // First, check to see if a moment collection already exists
        if acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection) == nil {
            // create a new TopShot Collection
            let collection <- TopShot.createEmptyCollection() as! @TopShot.Collection
            // Put the new Collection in storage
            acct.save(<-collection, to: /storage/MomentCollection)
        }

        if !acct.getCapability<&{
            NonFungibleToken.Receiver,
            NonFungibleToken.CollectionPublic,
            TopShot.MomentCollectionPublic,
            MetadataViews.ResolverCollection
        }>(/public/MomentCollection).check() {
            acct.unlink(/public/MomentCollection)
            // create a public capability for the collection
            acct.link<&TopShot.Collection{NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic, TopShot.MomentCollectionPublic, MetadataViews.ResolverCollection}>(/public/MomentCollection, target: /storage/MomentCollection) ??  panic("Could not link Topshot Collection Public Path");
        }

        // Create a PackNFT collection in the signer account if it doesn't already have one
        if acct.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath) == nil {
            acct.save(<- PackNFT.createEmptyCollection(), to: PackNFT.CollectionStoragePath);
            acct.link<&{NonFungibleToken.CollectionPublic}>(PackNFT.CollectionPublicPath, target: PackNFT.CollectionStoragePath)
        }

        // Create collection public capability if it doesn't already exist
        if !acct.getCapability<&{
            NonFungibleToken.CollectionPublic
        }>(PackNFT.CollectionPublicPath).check() {
            acct.unlink(PackNFT.CollectionPublicPath)
            acct.link<&{NonFungibleToken.CollectionPublic}>(PackNFT.CollectionPublicPath, target: PackNFT.CollectionStoragePath)
            ??  panic("Could not link Topshot PackNFT Collection Public Path");
        }
    }
}

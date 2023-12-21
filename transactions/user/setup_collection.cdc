import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

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
    }
}

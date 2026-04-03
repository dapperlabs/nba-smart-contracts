import TopShot from "TopShot"
import MetadataViews from "MetadataViews"

// Returns the MetadataViews.Medias for a given TopShot moment NFT,
// which includes IPFS media entries when CIDs have been set in TopShotIPFSResolver.

access(all) fun main(address: Address, momentID: UInt64): MetadataViews.Medias {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection)
        ?? panic("Could not borrow MomentCollectionPublic")

    let nft = collectionRef.borrowMoment(id: momentID)
        ?? panic("Could not borrow moment")

    return nft.resolveView(Type<MetadataViews.Medias>())! as! MetadataViews.Medias
}

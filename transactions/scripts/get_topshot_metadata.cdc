import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS


access(all) fun main(address: Address, id: UInt64): TopShot.TopShotMomentMetadataView {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)!

    let nft = collectionRef.borrowMoment(id: id)!
    
    // Get the Top Shot specific metadata for this NFT
    let view = nft.resolveView(Type<TopShot.TopShotMomentMetadataView>())!

    let metadata = view as! TopShot.TopShotMomentMetadataView
    
    return metadata
}
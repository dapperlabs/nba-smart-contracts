import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

access(all) struct NFT {
    access(all) let name: String
    access(all) let description: String
    access(all) let thumbnail: String
    access(all) let owner: Address
    access(all) let type: String
    access(all) let externalURL: String
    access(all) let storagePath: String
    access(all) let publicPath: String
    access(all) let collectionName: String
    access(all) let collectionDescription: String
    access(all) let collectionSquareImage: String
    access(all) let collectionBannerImage: String
    access(all) let royaltyReceiversCount: UInt32
    access(all) let traitsCount: UInt32
    access(all) let videoURL: String

    init(
            name: String,
            description: String,
            thumbnail: String,
            owner: Address,
            type: String,
            externalURL: String,
            storagePath: String,
            publicPath: String,
            collectionName: String,
            collectionDescription: String,
            collectionSquareImage: String,
            collectionBannerImage: String,
            royaltyReceiversCount: UInt32,
            traitsCount: UInt32,
            videoURL: String
    ) {
        self.name = name
        self.description = description
        self.thumbnail = thumbnail
        self.owner = owner
        self.type = type
        self.externalURL = externalURL
        self.storagePath = storagePath
        self.publicPath = publicPath
        self.collectionName = collectionName
        self.collectionDescription = collectionDescription
        self.collectionSquareImage = collectionSquareImage
        self.collectionBannerImage = collectionBannerImage
        self.royaltyReceiversCount = royaltyReceiversCount
        self.traitsCount = traitsCount
        self.videoURL = videoURL
    }
}

access(all) fun main(address: Address, id: UInt64): NFT {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection)!

    let nft = collectionRef.borrowMoment(id: id)!
    
    // Get all core views for this TopShot NFT
    let displayView = nft.resolveView(Type<MetadataViews.Display>())! as! MetadataViews.Display
    let collectionDisplayView = nft.resolveView(Type<MetadataViews.NFTCollectionDisplay>())! as! MetadataViews.NFTCollectionDisplay
    let collectionDataView = nft.resolveView(Type<MetadataViews.NFTCollectionData>())! as! MetadataViews.NFTCollectionData
    let royaltiesView = nft.resolveView(Type<MetadataViews.Royalties>())! as! MetadataViews.Royalties
    let externalURLView = nft.resolveView(Type<MetadataViews.ExternalURL>())! as! MetadataViews.ExternalURL
    let traitsView = nft.resolveView(Type<MetadataViews.Traits>())! as! MetadataViews.Traits
    let mediasView = nft.resolveView(Type<MetadataViews.Medias>())! as! MetadataViews.Medias

    let owner: Address = nft.owner!.address!
    let nftType = nft.getType()

    return NFT(
        name: displayView.name,
        description: displayView.description,
        thumbnail: displayView.thumbnail.uri(),
        owner: owner,
        type: nftType.identifier,
        externalURL: externalURLView.url,
        storagePath: collectionDataView.storagePath.toString(),
        publicPath: collectionDataView.publicPath.toString(),
        collectionName: collectionDisplayView.name,
        collectionDescription: collectionDisplayView.description,
        collectionSquareImage: collectionDisplayView.squareImage.file.uri(),
        collectionBannerImage: collectionDisplayView.bannerImage.file.uri(),
        royaltyReceiversCount: UInt32(royaltiesView.getRoyalties().length),
        traitsCount: UInt32(traitsView.traits.length),
        videoURL: mediasView.items[1].file.uri()
    )
}

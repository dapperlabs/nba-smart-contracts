import TopShot from 0xTOPSHOTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

pub struct NFT {
    pub let name: String
    pub let description: String
    pub let thumbnail: String
    pub let owner: Address
    pub let type: String
    pub let externalURL: String
    pub let storagePath: String
    pub let publicPath: String
    pub let privatePath: String
    pub let collectionName: String
    pub let collectionDescription: String
    pub let collectionSquareImage: String
    pub let collectionBannerImage: String
    pub let royaltyReceiversCount: UInt32
    pub let traitsCount: UInt32

    init(
            name: String,
            description: String,
            thumbnail: String,
            owner: Address,
            type: String,
            externalURL: String,
            storagePath: String,
            publicPath: String,
            privatePath: String,
            collectionName: String,
            collectionDescription: String,
            collectionSquareImage: String,
            collectionBannerImage: String,
            royaltyReceiversCount: UInt32,
            traitsCount: UInt32
    ) {
        self.name = name
        self.description = description
        self.thumbnail = thumbnail
        self.owner = owner
        self.type = type
        self.externalURL = externalURL
        self.storagePath = storagePath
        self.publicPath = publicPath
        self.privatePath = privatePath
        self.collectionName = collectionName
        self.collectionDescription = collectionDescription
        self.collectionSquareImage = collectionSquareImage
        self.collectionBannerImage = collectionBannerImage
        self.royaltyReceiversCount = royaltyReceiversCount
        self.traitsCount = traitsCount
    }
}

pub fun main(address: Address, id: UInt64): NFT {
    let account = getAccount(address)

    let collectionRef = account.getCapability(/public/MomentCollection)
                            .borrow<&{TopShot.MomentCollectionPublic}>()!

    let nft = collectionRef.borrowMoment(id: id)!
    
    // Get all core views for this TopShot NFT
    let displayView = nft.resolveView(Type<MetadataViews.Display>())! as! MetadataViews.Display
    let collectionDisplayView = nft.resolveView(Type<MetadataViews.NFTCollectionDisplay>())! as! MetadataViews.NFTCollectionDisplay
    let collectionDataView = nft.resolveView(Type<MetadataViews.NFTCollectionData>())! as! MetadataViews.NFTCollectionData
    let royaltiesView = nft.resolveView(Type<MetadataViews.Royalties>())! as! MetadataViews.Royalties
    let externalURLView = nft.resolveView(Type<MetadataViews.ExternalURL>())! as! MetadataViews.ExternalURL
    let traitsView = nft.resolveView(Type<MetadataViews.Traits>())! as! MetadataViews.Traits
    
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
        privatePath: collectionDataView.providerPath.toString(),
        collectionName: collectionDisplayView.name,
        collectionDescription: collectionDisplayView.description,
        collectionSquareImage: collectionDisplayView.squareImage.file.uri(),
        collectionBannerImage: collectionDisplayView.bannerImage.file.uri(),
        royaltyReceiversCount: UInt32(royaltiesView.getRoyalties().length),
        traitsCount: UInt32(traitsView.traits.length)
    )
}

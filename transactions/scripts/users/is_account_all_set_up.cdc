import TopShot from 0xTOPSHOTADDRESS
import NonFungibleToken from 0xNFTADDRESS
import PackNFT from 0xPACKNFTADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

// Check to see if an account looks like it has been set up to hold Pinnacle NFTs and PackNFTs.
pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.Receiver,
            NonFungibleToken.CollectionPublic,
            TopShot.MomentCollectionPublic,
            MetadataViews.ResolverCollection
        }>(/public/MomentCollection).check() &&
        account.getCapability<&{
            NonFungibleToken.CollectionPublic
        }>(PackNFT.CollectionPublicPath).check()
}
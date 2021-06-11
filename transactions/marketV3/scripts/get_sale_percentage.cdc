import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

pub fun main(sellerAddress: Address): UFix64 {
    let acct = getAccount(sellerAddress)
    let collectionRef = acct.getCapability(TopShotMarketV3.marketPublicPath).borrow<&{Market.SalePublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.cutPercentage
}
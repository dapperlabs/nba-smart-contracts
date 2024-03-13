import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

access(all) fun main(sellerAddress: Address, momentID: UInt64): UFix64 {

    let acct = getAccount(sellerAddress)
    let collectionRef = acct.capabilities.borrow<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketPublicPath)
        ?? panic("Could not borrow capability from public collection")
    
    let price = collectionRef.getPrice(tokenID: UInt64(momentID))
        ?? panic("Could not find price")

    return price
    
}
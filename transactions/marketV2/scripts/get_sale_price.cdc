import TopShotMarketV2 from 0xMARKETV2ADDRESS

pub fun main(sellerAddress: Address, momentID: UInt64): UFix64 {

    let acct = getAccount(sellerAddress)
    let collectionRef = acct.getCapability(/public/topshotSaleCollection).borrow<&{TopShotMarketV2.SalePublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getPrice(tokenID: UInt64(momentID))!
    
}
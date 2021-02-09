import TopShotMarketV2 from 0xMARKETV2ADDRESS

pub fun main(sellerAddress: Address, momentID: UInt64): UInt32 {
    let saleRef = getAccount(sellerAddress).getCapability(/public/topshotSaleCollection)
        .borrow<&{TopShotMarketV2.SalePublic}>()
        ?? panic("Could not get public sale reference")

    let token = saleRef.borrowMoment(id: momentID)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return data.setID
}
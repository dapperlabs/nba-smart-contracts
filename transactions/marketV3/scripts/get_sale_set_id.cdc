import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

access(all) fun main(sellerAddress: Address, momentID: UInt64): UInt32 {
    let saleRef = getAccount(sellerAddress).capabilities.borrow<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketPublicPath)
        ?? panic("Could not get public sale reference")

    let token = saleRef.borrowMoment(id: momentID)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return data.setID
}
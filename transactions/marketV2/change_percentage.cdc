import TopShotMarketV2 from 0xMARKETV2ADDRESS

transaction(newPercentage: UFix64) {
    prepare(acct: AuthAccount) {

        let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changePercentage(newPercentage)
    }
}
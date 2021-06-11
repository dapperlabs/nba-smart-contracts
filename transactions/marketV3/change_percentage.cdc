import TopShotMarketV3 from 0xMARKETV3ADDRESS

transaction(newPercentage: UFix64) {
    prepare(acct: AuthAccount) {

        let topshotSaleCollection = acct.borrow<&TopShotMarketV3.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changePercentage(newPercentage)
    }
}
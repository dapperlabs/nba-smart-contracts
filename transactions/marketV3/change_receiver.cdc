import TopShotMarketV3 from 0xMARKETV3ADDRESS

transaction(receiverPath: PublicPath) {
    prepare(acct: AuthAccount) {

        let topshotSaleCollection = acct.borrow<&TopShotMarketV3.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changeOwnerReceiver(acct.getCapability(receiverPath))
    }
}
import TopShotMarketV3 from 0xMARKETV3ADDRESS

transaction(receiverPath: PublicPath) {
    prepare(acct: auth(BorrowValue) &Account) {

        let topshotSaleCollection = acct.storage.borrow<auth(TopShotMarketV3.Update) &TopShotMarketV3.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changeOwnerReceiver(acct.capabilities.get<&{FungibleToken.Receiver}>(receiverPath)!)
    }
}
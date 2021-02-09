import Market from 0xMARKETADDRESS

transaction(receiverPath: PublicPath) {
    prepare(acct: AuthAccount) {

        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changeOwnerReceiver(acct.getCapability(receiverPath))
    }
}
import Market from 0xMARKETADDRESS

transaction(newPercentage: UFix64) {
    prepare(acct: AuthAccount) {

        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changePercentage(newPercentage)
    }
}
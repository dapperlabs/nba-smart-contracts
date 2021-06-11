import TopShot from 0xTOPSHOTADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction changes the price of a moment that a user has for sale

// Parameters:
//
// tokenID: the ID of the moment whose price is being changed
// newPrice: the new price of the moment

transaction(tokenID: UInt64, newPrice: UFix64) {
    prepare(acct: AuthAccount) {

        // borrow a reference to the owner's sale collection
        let topshotSaleCollection = acct.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath)
            ?? panic("Could not borrow from sale in storage")

        // Change the price of the moment
        topshotSaleCollection.listForSale(tokenID: tokenID, price: newPrice)
    }
}
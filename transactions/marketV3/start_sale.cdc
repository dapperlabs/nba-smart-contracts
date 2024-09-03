import TopShot from 0xTOPSHOTADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction is for a user to put a new moment up for sale
// They must have TopShot Collection and a TopShotMarketV2 Sale Collection already
// stored in their account

// Parameters
//
// momentId: the ID of the moment to be listed for sale
// price: the sell price of the moment

transaction(momentID: UInt64, price: UFix64) {
    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the topshot Sale Collection
        let topshotSaleCollection = acct.storage.borrow<auth(TopShotMarketV3.Create) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath)
            ?? panic("Could not borrow from sale in storage")

        // List the specified moment for sale
        topshotSaleCollection.listForSale(tokenID: momentID, price: price)
    }
}
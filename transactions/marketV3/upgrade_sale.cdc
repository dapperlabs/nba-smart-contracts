import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction is for a user to change a moment sale from
// the first version of the market contract to the third version

// Parameters
//
// tokenID: the ID of the moment whose sale is to be upgraded

transaction(tokenID: UInt64, price: UFix64) {

    prepare(acct: auth(BorrowValue) &Account) {

        // Borrow a reference to the NFT collection in the signers account	
        let nftCollection = acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")	

        // borrow a reference to the owner's sale collection
        let topshotSaleCollection = acct.storage.borrow<auth(Market.Withdraw) &Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        let topshotSaleV3Collection = acct.storage.borrow<auth(TopShotMarketV3.Create) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath)
            ?? panic("Could not borrow reference to sale V2 in storage")

        // withdraw the moment from the sale, thereby de-listing it
        let token <- topshotSaleCollection.withdraw(tokenID: tokenID)

        // deposit the moment into the owner's collection	
        nftCollection.deposit(token: <-token)

        // List the specified moment for sale
        topshotSaleV3Collection.listForSale(tokenID: tokenID, price: price)

    }
}
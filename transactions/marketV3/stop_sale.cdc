import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction is for a user to stop a moment sale in their account

// Parameters
//
// tokenID: the ID of the moment whose sale is to be delisted

transaction(tokenID: UInt64) {

    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the owner's sale collection
        if let topshotSaleV3Collection = acct.storage.borrow<auth(TopShotMarketV3.Cancel) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) {

            // cancel the moment from the sale, thereby de-listing it
            topshotSaleV3Collection.cancelSale(tokenID: tokenID)
            
        } else if let topshotSaleCollection = acct.storage.borrow<auth(Market.Withdraw) &Market.SaleCollection>(from: /storage/topshotSaleCollection) {
            // Borrow a reference to the NFT collection in the signers account
            let collectionRef = acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
                ?? panic("Could not borrow from MomentCollection in storage")
        
            // withdraw the moment from the sale, thereby de-listing it
            let token <- topshotSaleCollection.withdraw(tokenID: tokenID)

            // deposit the moment into the owner's collection
            collectionRef.deposit(token: <-token)
        }
    }
}
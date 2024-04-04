import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

// This transaction changes the price of a moment that a user has for sale

// Parameters:
//
// tokenID: the ID of the moment whose price is being changed
// newPrice: the new price of the moment

transaction(tokenID: UInt64, newPrice: UFix64) {

    // Local variable for the account's topshot sale collection
    let topshotSaleCollectionRef: auth(Market.Update) &Market.SaleCollection

    prepare(acct: auth(Storage, Capabilities) &Account) {

        // borrow a reference to the owner's sale collection
        self.topshotSaleCollectionRef = acct.storage.borrow<auth(Market.Update) &Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
    }

    execute {

        // Change the price of the moment
        self.topshotSaleCollectionRef.changePrice(tokenID: tokenID, newPrice: newPrice)
    }

    
}
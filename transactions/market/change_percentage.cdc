import Market from 0xMARKETADDRESS

// This transaction changes the percentage cut of a moment's sale given to beneficiary

// Parameters:
//
// newPercentage: new percentage of tokens the beneficiary will receive from the sale

transaction(newPercentage: UFix64) {

    // Local variable for the account's topshot sale collection
    let topshotSaleCollectionRef: auth(Market.Update) &Market.SaleCollection

    prepare(acct: auth(Storage, Capabilities) &Account) {

        // borrow a reference to the owner's sale collection
        self.topshotSaleCollectionRef = acct.storage.borrow<auth(Market.Update) &Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
    }

    execute {

        // Change the percentage of the moment
        self.topshotSaleCollectionRef.changePercentage(newPercentage)
    }

    post {

        self.topshotSaleCollectionRef.cutPercentage! == newPercentage: 
            "cutPercentage not changed"
    }
    
}
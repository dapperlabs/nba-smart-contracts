import Market from 0xMARKETADDRESS
import TopShot from 0xTOPSHOTADDRESS

// This transaction puts a moment owned by the user up for sale

// Parameters:
//
// tokenReceiverPath: token capability for the account who will receive tokens for purchase
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale
// momentID: ID of moment to be put on sale
// price: price of moment

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64, momentID: UInt64, price: UFix64) {

    // Local variables for the topshot collection and market sale collection objects
    let collectionRef: &TopShot.Collection
    let marketSaleCollectionRef: &Market.SaleCollection
    
    prepare(acct: AuthAccount) {

        // check to see if a sale collection already exists
        if acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection) == nil {

            // get the fungible token capabilities for the owner and beneficiary

            let ownerCapability = acct.getCapability(tokenReceiverPath)

            let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability(tokenReceiverPath)

            // create a new sale collection
            let topshotSaleCollection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
            
            // save it to storage
            acct.save(<-topshotSaleCollection, to: /storage/topshotSaleCollection)
        
            // create a public link to the sale collection
            acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
        }
        
        // borrow a reference to the seller's moment collection
        self.collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // borrow a reference to the sale
        self.marketSaleCollectionRef = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
    }

    execute {

        // withdraw the moment to put up for sale
        let token <- self.collectionRef.withdraw(withdrawID: momentID) as! @TopShot.NFT
        
        // the the moment for sale
        self.marketSaleCollectionRef.listForSale(token: <-token, price: UFix64(price))
    }
}
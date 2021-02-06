import Market from 0xMARKETADDRESS
import TopShot from 0xTOPSHOTADDRESS

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64, momentID: UInt64, price: UFix64) {
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
        let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // withdraw the moment to put up for sale
        let token <- nftCollection.withdraw(withdrawID: momentID) as! @TopShot.NFT

        // borrow a reference to the sale
        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
        
        // the the moment for sale
        topshotSaleCollection.listForSale(token: <-token, price: UFix64(price))
        
    }
}
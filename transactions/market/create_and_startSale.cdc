
import FungibleToken from 0xee82856bf20e2aa6
import TopShot from 0x179b6b1cb6755e31
import Market from 0xf3fcd2c1a78f5eee

transaction {
    prepare(acct: AuthAccount) {
        // check to see if a sale collection already exists
        if acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection) == nil {
            // get the fungible token capabilities for the owner and beneficiary
            let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(/public/sReceiver)!
            let beneficiaryCapability = getAccount(0x01).getCapability<&{FungibleToken.Receiver}>(/public/sReceiver)!

            let ownerCollection = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

            // create a new sale collection
            let topshotSaleCollection <- Market.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: 0.25)
            
            // save it to storage
            acct.save(<-topshotSaleCollection, to: /storage/topshotSaleCollection)
        
            // create a public link to the sale collection
            acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
        }

        // borrow a reference to the sale
        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        // set the new cut percentage
        topshotSaleCollection.changePercentage(0.30)
        
        // put the moment up for sale
        topshotSaleCollection.listForSale(tokenID: 2, price: 19.0)
        
    }
}
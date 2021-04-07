import FungibleToken from 0xFUNGIBLETOKENADDRESS
import TopShotMarketV2 from 0xMARKETV2ADDRESS
import TopShot from 0xTOPSHOTADDRESS

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64, momentID: UInt64, price: UFix64) {
    prepare(acct: AuthAccount) {
        // check to see if a sale collection already exists
        if acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath) == nil {
            // get the fungible token capabilities for the owner and beneficiary
            let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(tokenReceiverPath)
            let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability<&{FungibleToken.Receiver}>(tokenReceiverPath)

            let ownerCollection = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

            // create a new sale collection
            let topshotSaleCollection <- TopShotMarketV2.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
            
            // save it to storage
            acct.save(<-topshotSaleCollection, to: TopShotMarketV2.marketStoragePath)
        
            // create a public link to the sale collection
            acct.link<&TopShotMarketV2.SaleCollection{TopShotMarketV2.SalePublic}>(TopShotMarketV2.marketPublicPath, target: TopShotMarketV2.marketStoragePath)
        }

        // borrow a reference to the sale
        let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
            ?? panic("Could not borrow from sale in storage")
        
        // put the moment up for sale
        topshotSaleCollection.listForSale(tokenID: momentID, price: price)
        
    }
}
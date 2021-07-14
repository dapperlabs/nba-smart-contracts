import FungibleToken from 0xFUNGIBLETOKENADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64, momentID: UInt64, price: UFix64) {

    prepare(acct: AuthAccount) {
        // check to see if a v3 sale collection already exists
        if acct.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) == nil {
            // get the fungible token capabilities for the owner and beneficiary
            let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(tokenReceiverPath)
            let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability<&{FungibleToken.Receiver}>(tokenReceiverPath)

            let ownerCollection = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

            // get a capability for the v1 collection
            var v1SaleCollection: Capability<&Market.SaleCollection>? = nil
            if acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection) != nil {
                v1SaleCollection = acct.link<&Market.SaleCollection>(/private/topshotSaleCollection, target: /storage/topshotSaleCollection)!
            }

            // create a new sale collection
            let topshotSaleCollection <- TopShotMarketV3.createSaleCollection(ownerCollection: ownerCollection,
                                                                             ownerCapability: ownerCapability,
                                                                             beneficiaryCapability: beneficiaryCapability,
                                                                             cutPercentage: cutPercentage,
                                                                             marketV1Capability: v1SaleCollection)
            
            // save it to storage
            acct.save(<-topshotSaleCollection, to: TopShotMarketV3.marketStoragePath)
        
            // create a public link to the sale collection
            acct.link<&TopShotMarketV3.SaleCollection{Market.SalePublic}>(TopShotMarketV3.marketPublicPath, target: TopShotMarketV3.marketStoragePath)
        }

        // borrow a reference to the sale
        let topshotSaleCollection = acct.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath)
            ?? panic("Could not borrow from sale in storage")
        
        // put the moment up for sale
        topshotSaleCollection.listForSale(tokenID: momentID, price: price)
        
    }
}
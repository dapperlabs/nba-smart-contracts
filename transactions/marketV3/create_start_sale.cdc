import FungibleToken from 0xFUNGIBLETOKENADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import NonFungibleToken from 0xNFTADDRESS

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64, momentID: UInt64, price: UFix64) {

    prepare(acct: auth(Storage, Capabilities) &Account) {
        // check to see if a v3 sale collection already exists
        if acct.storage.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) == nil {
            // get the fungible token capabilities for the owner and beneficiary
            let ownerCapability = acct.capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)!
            let beneficiaryCapability = getAccount(beneficiaryAccount).capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)!

            var ownerCollection = acct.storage.copy<Capability<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>>(from: /storage/MomentCollectionCap)
            if ownerCollection == nil {
                ownerCollection = acct.capabilities.storage.issue<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>(/storage/MomentCollection)
                acct.storage.save(ownerCollection, to: /storage/MomentCollectionCap)
            }

            // get a capability for the v1 collection
            var v1SaleCollection = acct.storage.copy<Capability<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>>(from: /storage/topshotSaleCollectionCap)
            if v1SaleCollection == nil {
                v1SaleCollection = acct.capabilities.storage.issue<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>(/storage/topshotSaleCollection)
                acct.storage.save(v1SaleCollection, to: /storage/topshotSaleCollectionCap)
            }

            // create a new sale collection
            let topshotSaleCollection <- TopShotMarketV3.createSaleCollection(ownerCollection: ownerCollection,
                                                                             ownerCapability: ownerCapability,
                                                                             beneficiaryCapability: beneficiaryCapability,
                                                                             cutPercentage: cutPercentage,
                                                                             marketV1Capability: v1SaleCollection)
            
            // save it to storage
            acct.storage.save(<-topshotSaleCollection, to: TopShotMarketV3.marketStoragePath)
        
            // create a public link to the sale collection
           acct.capabilities.publish(
                acct.capabilities.storage.issue<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketStoragePath),
                at: TopShotMarketV3.marketPublicPath
            )
        }

        // borrow a reference to the sale
        let topshotSaleCollection = acct.storage.borrow<auth(TopShotMarketV3.Create) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath)
            ?? panic("Could not borrow from sale in storage")
        
        // put the moment up for sale
        topshotSaleCollection.listForSale(tokenID: momentID, price: price)
        
    }
}
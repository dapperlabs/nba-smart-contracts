import FungibleToken from 0xFUNGIBLETOKENADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import NonFungibleToken from 0xNFTADDRESS

/// This transaction creates a V3 Sale Collection
/// in a user's account and lists a Moment for Sale in that collection
/// If a user already has a V3 Sale Collection
/// the transaction only lists the moment for sale
///
/// When creating a V3 Sale Collection, if the user already has a V1 Sale Collection,
/// the transaction will create and store a provider capability for that V1 Sale Collection
/// to be used with the V3 Sale Collection

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64, momentID: UInt64, price: UFix64) {

    prepare(acct: auth(Storage, Capabilities) &Account) {
        // check to see if a v3 sale collection already exists
        if acct.storage.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) == nil {
            // If the V3 Sale Collection does not exist, set up a new one

            // get the fungible token capabilities for the owner and beneficiary
            let ownerCapability = acct.capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)
            if !ownerCapability.check() {
                panic("Could not get the owner's FungibleToken.Receiver capability from ".concat(tokenReceiverPath.toString()))
            }
            let beneficiaryCapability = getAccount(beneficiaryAccount).capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)
            if !beneficiaryCapability.check() {
                panic("Could not get the beneficiary's FungibleToken.Receiver capability from ".concat(tokenReceiverPath.toString()))
            }

            // Get the owner's TopShot Collection Provider Capability that
            // allows the V3 sale collection to withdraw when sales are made
            var ownerCollection = acct.storage.copy<Capability<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>>(from: /storage/MomentCollectionCap)
            if ownerCollection == nil {
                // If the moment collection capabilitity does not already exist,
                // Issue a new one and store it in the standard private moment collection capability path
                ownerCollection = acct.capabilities.storage.issue<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>(/storage/MomentCollection)
                acct.storage.save(ownerCollection, to: /storage/MomentCollectionCap)
            }

            // get a capability for the v1 collection
            // Only accounts that existed before the V3 Sale Collection contract was deployed
            // will have this
            var v1SaleCollection = acct.storage.copy<Capability<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>>(from: /storage/topshotSaleCollectionCap)
            if v1SaleCollection == nil {
                // If the account doesn't have a V1 Sale Collection capability already,
                // first check if they even have a V1 Sale Collection at all
                if acct.storage.borrow<auth(Market.Create) &Market.SaleCollection>(from: /storage/topshotSaleCollection) != nil {
                    // If they have a V1 Sale Collection, issue a capability for it
                    // and store it in storage
                    v1SaleCollection = acct.capabilities.storage.issue<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>(/storage/topshotSaleCollection)
                    acct.storage.save(v1SaleCollection, to: /storage/topshotSaleCollectionCap)
                }
            }

            // create a new sale collection
            // V1SaleCollection will still be `nil` here if a V1 Sale Collection
            // did not exist in the authorizer's account
            // We can force-unwrap `ownerCollection` because it was already guaranteed to be non-`nil` above
            let topshotV3SaleCollection <- TopShotMarketV3.createSaleCollection(ownerCollection: ownerCollection!,
                                                                             ownerCapability: ownerCapability,
                                                                             beneficiaryCapability: beneficiaryCapability,
                                                                             cutPercentage: cutPercentage,
                                                                             marketV1Capability: v1SaleCollection)
            
            // save it to storage
            acct.storage.save(<-topshotV3SaleCollection, to: TopShotMarketV3.marketStoragePath)
        
            // create a public link to the sale collection
           acct.capabilities.publish(
                acct.capabilities.storage.issue<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketStoragePath),
                at: TopShotMarketV3.marketPublicPath
            )
        }

        // borrow a reference to the sale
        let topshotSaleCollection = acct.storage.borrow<auth(TopShotMarketV3.Create) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath)
            ?? panic("Could not borrow the owner's Top Shot V3 Sale Collection in storage from ".concat(TopShotMarketV3.marketStoragePath.toString()))
        
        // put the moment up for sale
        topshotSaleCollection.listForSale(tokenID: momentID, price: price)
    }
}
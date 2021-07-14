import FungibleToken from 0xFUNGIBLETOKENADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction creates a sale collection and stores it in the signer's account
// It does not put an NFT up for sale

// Parameters
// 
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64) {
    prepare(acct: AuthAccount) {
        let ownerCapability = acct.getCapability<&AnyResource{FungibleToken.Receiver}>(tokenReceiverPath)

        let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability<&AnyResource{FungibleToken.Receiver}>(tokenReceiverPath)

        let ownerCollection = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

        let collection <- TopShotMarketV3.createSaleCollection(ownerCollection: ownerCollection,
                                                               ownerCapability: ownerCapability,
                                                               beneficiaryCapability: beneficiaryCapability,
                                                               cutPercentage: cutPercentage,
                                                               marketV1Capability: nil)
        
        acct.save(<-collection, to: TopShotMarketV3.marketStoragePath)
        
        acct.link<&TopShotMarketV3.SaleCollection{Market.SalePublic}>(TopShotMarketV3.marketPublicPath, target: TopShotMarketV3.marketStoragePath)
    }
}

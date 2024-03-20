import FungibleToken from 0xFUNGIBLETOKENADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This transaction creates a sale collection and stores it in the signer's account
// It does not put an NFT up for sale

// Parameters
// 
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64) {
    prepare(acct: auth(Storage, Capabilities) &Account) {
        let ownerCapability = acct.capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)!

        let beneficiaryCapability = getAccount(beneficiaryAccount).capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)!

        let ownerCollection = acct.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &TopShot.Collection>(/storage/MomentCollection)

        let collection <- TopShotMarketV3.createSaleCollection(ownerCollection: ownerCollection,
                                                               ownerCapability: ownerCapability,
                                                               beneficiaryCapability: beneficiaryCapability,
                                                               cutPercentage: cutPercentage,
                                                               marketV1Capability: nil)
        
        acct.storage.save(<-collection, to: TopShotMarketV3.marketStoragePath)

        acct.capabilities.publish(
            acct.capabilities.storage.issue<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketStoragePath),
            at: TopShotMarketV3.marketPublicPath
        )
    }
}


import FungibleToken from 0xee82856bf20e2aa6
import TopShot from 0x179b6b1cb6755e31
import Market from 0xf3fcd2c1a78f5eee

// This transaction creates a sale collection and stores it in the signer's account
// It does not put an NFT up for sale

// Parameters
// 
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale

transaction(beneficiaryAccount: Address, cutPercentage: UFix64) {
    prepare(acct: AuthAccount) {
        let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(/public/%sReceiver)
        let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability<&{FungibleToken.Receiver}>(/public/%sReceiver)

        let ownerCollection: Capability<&TopShot.Collection> = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

        let collection <- Market.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
        
        acct.save(<-collection, to: Market.marketStoragePath)
        
        acct.link<&Market.SaleCollection{Market.SalePublic}>(Market.marketPublicPath, target: Market.marketStoragePath)
    }
}

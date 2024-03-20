import Market from 0xMARKETADDRESS
import FungibleToken from 0xFUNGIBLETOKENADDRESS

// This transaction creates a public sale collection capability that any user can interact with

// Parameters:
//
// tokenReceiverPath: token capability for the account who will receive tokens for purchase
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64) {

    prepare(acct: auth(Storage, Capabilities) &Account) {
        
        let ownerCapability = acct.capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)!

        let beneficiaryCapability = getAccount(beneficiaryAccount).capabilities.get<&{FungibleToken.Receiver}>(tokenReceiverPath)!

        let collection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
        
        acct.storage.save(<-collection, to: /storage/topshotSaleCollection)
        acct.capabilities.publish(
            acct.capabilities.storage.issue<&Market.SaleCollection>(/storage/topshotSaleCollection),
            at: /public/topshotSaleCollection
        )
    }
}
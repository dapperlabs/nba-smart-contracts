import Market from 0xMARKETADDRESS

// This transaction creates a public sale collection capability that any user can interact with

// Parameters:
//
// tokenReceiverPath: token capability for the account who will receive tokens for purchase
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale

transaction(tokenReceiverPath: PublicPath, beneficiaryAccount: Address, cutPercentage: UFix64) {

    prepare(acct: AuthAccount) {
        
        let ownerCapability = acct.getCapability(tokenReceiverPath)

        let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability(tokenReceiverPath)

        let collection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
        
        acct.save(<-collection, to: /storage/topshotSaleCollection)
        
        acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
    }
}
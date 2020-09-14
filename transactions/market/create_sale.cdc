import Market from 0xMARKETADDRESS

// Parameters
// 
// beneficiaryAccount: the Flow address of the account where a cut of the purchase will be sent
// cutPercentage: how much in percentage the beneficiary will receive from the sale

transaction(beneficiaryAccount: Address, cutPercentage: UFix64) {
    prepare(acct: AuthAccount) {
        let ownerCapability = acct.getCapability(/public/%sReceiver)!
        let beneficiaryCapability = getAccount(beneficiaryAccount).getCapability(/public/%sReceiver)!

        let collection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
        
        acct.save(<-collection, to: /storage/topshotSaleCollection)
        
        acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
    }
}
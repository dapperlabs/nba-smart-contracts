import Market from 0xMARKETADDRESS

// This script gets the percentage cut that beneficiary will take
// of moments in an account's sale collection

// Parameters:
//
// sellerAddress: The Flow Address of the account whose sale collection needs to be read

// Returns: UFix64
// The percentage cut of an account's sale collection

pub fun main(sellerAddress: Address): UFix64 {

    let acct = getAccount(sellerAddress)

    let collectionRef = acct.getCapability(/public/topshotSaleCollection).borrow<&{Market.SalePublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.cutPercentage
}
import Market from 0xMARKETADDRESS

// This script gets the number of moments an account has for sale

// Parameters:
//
// sellerAddress: The Flow Address of the account whose sale collection needs to be read

// Returns: Int
// Number of moments up for sale in an account

access(all) fun main(sellerAddress: Address): Int {

    let acct = getAccount(sellerAddress)

    let collectionRef = acct.capabilities.borrow<&Market.SaleCollection>(/public/topshotSaleCollection)
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}
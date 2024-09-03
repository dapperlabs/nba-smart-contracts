import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction is for a user to purchase a moment that another user
// has for sale in their sale collection

// Parameters
//
// sellerAddress: the Flow address of the account issuing the sale of a moment
// tokenID: the ID of the moment being purchased
// purchaseAmount: the amount for which the user is paying for the moment; must not be less than the moment's price

transaction(sellerAddress: Address, tokenID: UInt64, purchaseAmount: UFix64) {
    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the signer's collection
        let collection = acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow reference to the Moment Collection")

        // borrow a reference to the signer's fungible token Vault
        let provider = acct.storage.borrow<auth(FungibleToken.Withdraw) &DapperUtilityCoin.Vault>(from: /storage/dapperUtilityCoinVault)!
        
        // withdraw tokens from the signer's vault
        let tokens <- provider.withdraw(amount: purchaseAmount) as! @DapperUtilityCoin.Vault

        // get the seller's public account object
        let seller = getAccount(sellerAddress)

        // borrow a public reference to the seller's sale collection
        let topshotSaleCollection = seller.capabilities.borrow<&TopShotMarketV3.SaleCollection>(/public/topshotSalev3Collection)
            ?? panic("Could not borrow public sale reference")
    
        // purchase the moment
        let purchasedToken <- topshotSaleCollection.purchase(tokenID: tokenID, buyTokens: <-tokens)

        // deposit the purchased moment into the signer's collection
        collection.deposit(token: <-purchasedToken)
    }
}
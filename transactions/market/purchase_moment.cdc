import FungibleToken from 0xFUNGIBLETOKENADDRESS
import %[1]s from 0x%[2]s
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

transaction(sellerAddress: Address, tokenID: UInt64, purchaseAmount: UFix64) {
    prepare(acct: AuthAccount) {
        let seller = getAccount(sellerAddress)

        let collection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow reference to the Moment Collection")

        let provider = acct.borrow<&%[1]s.Vault{FungibleToken.Provider}>(from: /storage/%[5]sVault)!
        
        let tokens <- provider.withdraw(amount: purchaseAmount) as! @%[1]s.Vault

        let topshotSaleCollection = seller.getCapability(/public/topshotSaleCollection)!
            .borrow<&{Market.SalePublic}>()
            ?? panic("Could not borrow public sale reference")
    
        let purchasedToken <- topshotSaleCollection.purchase(tokenID: tokenID, buyTokens: <-tokens)

        collection.deposit(token: <-purchasedToken)
    }
}
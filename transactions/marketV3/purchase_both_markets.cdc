import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction purchases a moment from the v3 sale collection
// The v3 sale collection will also check the v1 collection for for sale moments as part of the purchase
// If there is no v3 sale collection, the transaction will just purchase it from v1 anyway

transaction(seller: Address, recipient: Address, momentID: UInt64, purchaseAmount: UFix64) {

    let purchaseTokens: @DapperUtilityCoin.Vault

    prepare(acct: AuthAccount) {

        // Borrow a provider reference to the buyers vault
        let provider = acct.borrow<&DapperUtilityCoin.Vault{FungibleToken.Provider}>(from: /storage/dapperUtilityCoinVault)
            ?? panic("Could not borrow a reference to the buyers FlowToken Vault")
        
        // withdraw the purchase tokens from the vault
        self.purchaseTokens <- provider.withdraw(amount: purchaseAmount) as! @DapperUtilityCoin.Vault
        
    }

    execute {

        // get the accounts for the seller and recipient
        let seller = getAccount(seller)
        let recipient = getAccount(recipient)

        // Get the reference for the recipient's nft receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Could not borrow a reference to the recipients moment collection")

        if let marketV3CollectionRef = seller.getCapability(/public/topshotSalev3Collection)
                .borrow<&{Market.SalePublic}>() {

            let purchasedToken <- marketV3CollectionRef.purchase(tokenID: momentID, buyTokens: <-self.purchaseTokens)
            receiverRef.deposit(token: <-purchasedToken)

        } else if let marketV1CollectionRef = seller.getCapability(/public/topshotSaleCollection)
            .borrow<&{Market.SalePublic}>() {
            // purchase the moment
            let purchasedToken <- marketV1CollectionRef.purchase(tokenID: momentID, buyTokens: <-self.purchaseTokens)

            // deposit the purchased moment into the signer's collection
            receiverRef.deposit(token: <-purchasedToken)

        } else {
            destroy self.purchaseTokens // REMOVE THIS
            panic("Could not borrow reference to either Sale collection")
        }
    }
}
 
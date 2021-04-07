import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV2 from 0xMARKETV2ADDRESS

// This transaction purchases a moment by first checking if it is in the first version of the market collecion
// If it isn't in the first version, it checks if it is in the second and purchases it there

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

        // Check if the V1 market collection exists
        if let marketCollection = seller.getCapability(/public/topshotSaleCollection)!
            .borrow<&{Market.SalePublic}>() {
                
            // Check if the V1 market has the moment for sale
            if marketCollection.borrowMoment(id: momentID) != nil {

                // purchase from the V1 market
                let purchasedToken <- marketCollection.purchase(tokenID: momentID, buyTokens: <-tokens)
                receiverRef.deposit(token: <-purchasedToken)
            }
            
        // Check if the V2 market collection exists
        } else if let marketV2Collection = seller.getCapability(/public/topshotSalev2Collection)!
                .borrow<&{MarketV2.SalePublic}>() {
        
            // Check if the V2 market has the moment for sale
            if marketV2Collection.borrowMoment(id: momentID) != nil {

                // Purchase from the V2 market
                let purchasedToken <- marketV2Collection.purchase(tokenID: momentID, buyTokens: <-tokens)
                receiverRef.deposit(token: <-purchasedToken)
            } else {
                panic("Could not find the moment sale in either collection")
            }
        } else {
            panic("Could not find either sale collection")
        }
    }
}
 
import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction is for a user to purchase a group of moments from
// one more or more sellers

// Parameters
//
// momentsBySeller: An object consisting of a key of the sellers address,
//  and an array of the moments being purchased from this seller
//
// purchaseAmount: the amount the user is paying for all moments within
//  the group

transaction(momentsBySeller: {Address: [UInt64]}, purchaseAmount: UFix64) {

    // Local variables for the topshot collection object and token provider
    let collectionRef: &TopShot.Collection
    let providerRef: auth(FungibleToken.Withdraw) &DapperUtilityCoin.Vault
    
    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the signer's collection
        self.collectionRef = acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow reference to the Moment Collection")

        // borrow a reference to the signer's fungible token Vault
        self.providerRef = acct.storage.borrow<auth(FungibleToken.Withdraw) &DapperUtilityCoin.Vault>(from: /storage/dapperUtilityCoinVault)!
    }

    execute {
        // Obtain a list of seller addresses
        var sellerAddresses = momentsBySeller.keys

        // Initialize the sum price of all moments
        var sumMomentPrices: UFix64 = 0.00

        for sellerAddress in sellerAddresses {
            // Get all moments we are purchasing from this seller
            var sellerMoments = momentsBySeller[sellerAddress]!
            
            for sellerMoment in sellerMoments {
                // Get the seller account
                let seller = getAccount(sellerAddress)
                // Check if we can obtain a reference to the sellers marketV3 collection
                if let marketV3CollectionRef = seller.capabilities.borrow<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketPublicPath) {

                    // Check the moments sale price
                    var momentPrice = marketV3CollectionRef.getPrice(tokenID: sellerMoment) ?? panic("Moment not for sale")
                    // Add the sale price to the sum of all moment prices
                    sumMomentPrices = sumMomentPrices + momentPrice
                    // Withdraw fungible tokens for payment
                    let tokens <- self.providerRef.withdraw(amount: momentPrice) as! @DapperUtilityCoin.Vault
                    // Purchase non-fungible token with payment via fungible tokens
                    let purchasedToken <- marketV3CollectionRef.purchase(tokenID: sellerMoment, buyTokens: <-tokens)
                    // Deposit purchased non-fungible token to purchasers collection
                    self.collectionRef.deposit(token: <-purchasedToken)


                // If we could not obtain reference to sellers marketV3 collection, try V1
                } else if let topshotSaleCollection = seller.capabilities.borrow<&Market.SaleCollection>(/public/topshotSaleCollection) {

                // Check the moments sale price
                var momentPrice = topshotSaleCollection.getPrice(tokenID: sellerMoment) ?? panic("Moment not for sale")
                // Add the sale price to the sum of all moment prices
                sumMomentPrices = sumMomentPrices + momentPrice
                // Withdraw fungible tokens for payment
                let tokens <- self.providerRef.withdraw(amount: momentPrice) as! @DapperUtilityCoin.Vault
                // Purchase non-fungible token with payment via fungible tokens
                let purchasedToken <- topshotSaleCollection.purchase(tokenID: sellerMoment, buyTokens: <-tokens)
                // Deposit purchased non-fungible token to purchasers collection
                self.collectionRef.deposit(token: <-purchasedToken)

                } else {
                    // Could not borrow a reference to sellers marketV1 or V3 sale collection
                    panic("Could not borrow reference to either Sale collection")
                }

            }
        }
        if sumMomentPrices > purchaseAmount {
            // Revert the transaction if the amount of fungible tokens required 
            // are larger than the users purchaseAmount
            panic("Sum of all moment prices is greater than purchaseAmount!")
        }

    }
}



import TopShot from 0x03
import Market from 0x04
import FungibleToken, FlowToken from 0x01

// this is the transacion a user would run if they want
// to put one of their moments up for sale

// values that would be configurable are
// withdrawID of the token for sale 
// price of the token
transaction {

    let collectionRef: &TopShot.Collection

    let saleRef: &Market.SaleCollection

    prepare(acct: AuthAccount) {

        // remove the sale collection from storage
        if acct.borrow<&Market.SaleCollection>(from: /storage/MomentSale) == nil {
            // this branch executes if there isn't a sale already
            // it creates a new sale collection and puts it in storage

            // get the FlowToken receiver reference to the Vault
            let receiverRef = acct.getCapability(/public/flowTokenReceiver)!
                                .borrow<&{FungibleToken.Receiver}>()!

            // create a new empty sale collection
            let sale <- Market.createSaleCollection(ownerVault: receiverRef, cutPercentage: 0.10)

            // put the sale back into storage
            acct.save(<-sale, to: /storage/MomentSale)

            // publish a reference to the sale
            acct.link<&{Market.SalePublic}>(/public/MomentSale, target: /storage/MomentSale)
        }

        // create a temporary reference to the stored collection
        self.collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!

        // create a temporary reference to the sale
        self.saleRef = acct.borrow<&Market.SaleCollection>(from: /storage/MomentSale)!
    }

    execute {
        // withdraw the token from the collection
        let token <- self.collectionRef.withdraw(withdrawID: 0)

        // put the token up for sale
        self.saleRef.listForSale(token: <-token, price: 30.00)

        log("Token put up for sale")
        log(1)
        log("Price:")
        log(30)
    }

    post {
        self.saleRef.getPrice(tokenID: 0) != nil:
            "Token should have been put up for sale and should not be nil"
    }
}
 
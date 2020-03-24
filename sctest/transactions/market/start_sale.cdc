import TopShot from 0x02
import Market from 0x03
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
        if acct.storage[Market.SaleCollection] == nil {
            // this branch executes if there isn't a sale already
            // it creates a new sale collection and puts it in storage

            // get the FlowToken receiver reference to the Vault
            let receiverRef = acct.published[&FungibleToken.Receiver] ?? panic("No receiver ref!")

            // create a new empty sale collection
            let sale <- Market.createSaleCollection(ownerVault: receiverRef, cutPercentage: 10)

            // put the sale back into storage
            let oldSale <- acct.storage[Market.SaleCollection] <- sale
            destroy oldSale

            // publish a reference to the sale
            acct.published[&Market.SalePublic] = &acct.storage[Market.SaleCollection] as &Market.SalePublic
        }

        // create a temporary reference to the stored collection
        self.collectionRef = &acct.storage[TopShot.Collection] as &TopShot.Collection

        // create a temporary reference to the sale
        self.saleRef = &acct.storage[Market.SaleCollection] as &Market.SaleCollection
    }

    execute {
        // withdraw the token from the collection
        let token <- self.collectionRef.withdraw(withdrawID: 1)

        // put the token up for sale
        self.saleRef.listForSale(token: <-token, price: 30)

        log("Token put up for sale")
        log(1)
        log("Price:")
        log(30)
    }

    post {
        self.saleRef.idPrice(tokenID: 1) != nil:
            "Token should have been put up for sale and should not be nil"
    }
}
 
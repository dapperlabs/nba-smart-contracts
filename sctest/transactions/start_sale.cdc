import TopShot from 0x02
import Market from 0x03
import FungibleToken, FlowToken from 0x01

// this is the transacion a user would run if they want
// to put one of their moments up for sale

// values that would be configurable are
// withdrawID of the token for sale 
// price of the token
transaction {

    prepare(acct: Account) {

        // create a reference to the stored collection
        let collectionRef = &acct.storage[TopShot.Collection] as TopShot.Collection

        // remove the sale collection from storage
        if acct.storage[Market.SaleCollection] != nil {

            let saleRef = &acct.storage[Market.SaleCollection] as Market.SaleCollection

            // withdraw the token from the collection
            let token <- collectionRef.withdraw(withdrawID: 1)

            // put the token up for sale
            saleRef.listForSale(token: <-token, price: 30)

        } else {
            // this branch executes if there isn't a sale already
            // it creates a new sale collection and puts the same moment up for sale

            // get the FlowToken receiver reference to the Vault
            let receiverRef = acct.published[&FungibleToken.Receiver] ?? panic("No receiver ref!")

            // create a new empty sale collection
            let sale <- Market.createSaleCollection(ownerVault: receiverRef)

            // withdraw the token from the collection
            let token <- collectionRef.withdraw(withdrawID: 1)

            // put the token up for sale
            sale.listForSale(token: <-token, price: 10)

            // put the sale back into storage
            let oldSale <- acct.storage[Market.SaleCollection] <- sale
            destroy oldSale

            // publish a reference to the sale
            acct.published[&Market.SalePublic] = &acct.storage[Market.SaleCollection] as Market.SalePublic
        }

        log("Token put up for sale")
        log(1)
    }
}
 
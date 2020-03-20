import TopShot from 0x02
import Market from 0x03
import FungibleToken, FlowToken from 0x01

// this is the transacion a user would run if they want
// to remove an NFT from their Sale and cancel the sale

// values that would be configurable are
// tokenID that is being removed from the sale
transaction {

    let collectionRef: &TopShot.Collection

    let saleRef: &Market.SaleCollection

    prepare(acct: Account) {

        // create a temporary reference to the stored collection
        self.collectionRef = &acct.storage[TopShot.Collection] as &TopShot.Collection

        // create a temporary reference to the sale
        self.saleRef = &acct.storage[Market.SaleCollection] as &Market.SaleCollection
    }

    execute {
        // withdraw the token from the sale
        let token <- self.saleRef.withdraw(tokenID: 1)

        // put the token up for sale
        self.collectionRef.deposit(token: <-token)

        log("Token withdrawn from the sale")
        log(1)
    }

    post {
        self.saleRef.idPrice(tokenID: 1) == nil:
            "Moment should have been removed from the sale!"
    }
}
 
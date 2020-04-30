import TopShot from 0x03
import Market from 0x04
import FungibleToken, FlowToken from 0x01

// this is the transacion a user would run if they want
// to remove an NFT from their Sale and cancel the sale

// values that would be configurable are
// tokenID that is being removed from the sale
transaction {

    let collectionRef: &TopShot.Collection

    let saleRef: &Market.SaleCollection

    prepare(acct: AuthAccount) {

        // create a temporary reference to the stored collection
        self.collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!

        // create a temporary reference to the sale
        self.saleRef = acct.borrow<&Market.SaleCollection>(from: /storage/MomentSale)!
    }

    execute {
        // withdraw the token from the sale
        let token <- self.saleRef.withdraw(tokenID: 0)

        // put the token back in the normal collection
        self.collectionRef.deposit(token: <-token)

        log("Token withdrawn from the sale")
        log(1)
    }

    post {
        self.saleRef.getPrice(tokenID: 0) == nil:
            "Moment should have been removed from the sale!"
    }
}
 
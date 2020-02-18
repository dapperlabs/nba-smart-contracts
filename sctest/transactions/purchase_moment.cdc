import TopShot from 0x02
import Market from 0x03
import FungibleToken, FlowToken from 0x01

// this is the transacion a user would run if they want
// to purchase a moment from the marketplace

// values that would be configurable are
// withdrawID of the token being bought
// price of the token
transaction {

    prepare(acct: Account) {

        // create a reference to the stored collection
        let collectionRef = &acct.storage[TopShot.Collection] as &TopShot.Collection

        let seller = getAccount(0x02)

        // remove the sale collection from storage
        if let saleRef = seller.published[&Market.SalePublic] {

            let vaultRef = &acct.storage[FlowToken.Vault] as &FlowToken.Vault
            let buyTokens <- vaultRef.withdraw(amount: 30)

            saleRef.purchase(tokenID: 1, recipient: collectionRef, buyTokens: <-buyTokens)

            log("token bought!")
        } else {
            // this branch executes if there isn't a sale
            panic("No sale!")
        }
    }
}
 
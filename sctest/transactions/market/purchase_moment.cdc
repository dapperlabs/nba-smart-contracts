import TopShot from 0x03
import Market from 0x04
import FungibleToken, FlowToken from 0x01

// this is the transacion a user would run if they want
// to purchase a moment from the marketplace

// values that would be configurable are
// withdrawID of the token being bought
// price of the token
transaction {

    // temporary reference for the signers moment collection
    let collectionRef: &TopShot.Collection

    // temp reference for the signer's Vault
    let vaultRef: &FlowToken.Vault

    prepare(acct: AuthAccount) {

        // create a temporary reference to the stored collection
        self.collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!

        // create reference to Vault
        self.vaultRef = acct.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)!
    }

    execute {
        // get the sellers public account object
        let seller = getAccount(0x03)

        // get the public capability and reference to the sellers Sale
        if let saleRef = seller.getCapability(/public/MomentSale)!.borrow<&{Market.SalePublic}>() {

            let buyTokens <- self.vaultRef.withdraw(amount: 40.00)

            saleRef.purchase(tokenID: 0, recipient: self.collectionRef, buyTokens: <-buyTokens)

            log("token bought!")
        } else {
            // this branch executes if there isn't a sale
            panic("No sale!")
        }

    }
}
 
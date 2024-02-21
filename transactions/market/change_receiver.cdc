import Market from 0xMARKETADDRESS
import FungibleToken from 0xFUNGIBLETOKENADDRESS

// This transaction changes the path which receives tokens for purchases of an account

// Parameters:
//
// receiverPath: The new fungible token capability for the account who receives tokens for purchases

transaction(receiverPath: PublicPath) {

    // Local variables for the sale collection object and receiver
    let saleCollectionRef: auth(Market.Update) &Market.SaleCollection
    let receiverPathRef: Capability<&{FungibleToken.Receiver}>

    prepare(acct: auth(BorrowValue) &Account) {

        self.saleCollectionRef = acct.storage.borrow<auth(Market.Update) &Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
        self.receiverPathRef = acct.capabilities.get<&{FungibleToken.Receiver}>(receiverPath)!
    }

    execute {

        self.saleCollectionRef.changeOwnerReceiver(self.receiverPathRef)

    }
}
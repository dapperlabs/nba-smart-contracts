import Market from 0xMARKETADDRESS

transaction(receiverPath: PublicPath) {

    // Local variables for the sale collection object and receiver
    let saleCollectionRef: &Market.SaleCollection
    let receiverPathRef: Capability

    prepare(acct: AuthAccount) {

        self.saleCollectionRef = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        self.receiverPathRef = acct.getCapability(receiverPath)
    }

    execute {

        self.saleCollectionRef.changeOwnerReceiver(self.receiverPathRef)

    }
}
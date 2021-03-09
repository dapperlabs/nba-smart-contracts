import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

// This transaction is for a user to stop a moment sale in their account
// by withdrawing that moment from their sale collection and depositing
// it into their normal moment collection

// Parameters
//
// tokenID: the ID of the moment whose sale is to be delisted

transaction(tokenID: UInt64) {

    let collectionRef: &TopShot.Collection
    let saleCollectionRef: &Market.SaleCollection

    prepare(acct: AuthAccount) {

        // Borrow a reference to the NFT collection in the signers account
        self.collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // borrow a reference to the owner's sale collection
        self.saleCollectionRef = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
    }

    execute{
    
        // withdraw the moment from the sale, thereby de-listing it
        let token <- self.saleCollectionRef.withdraw(tokenID: tokenID)

        // deposit the moment into the owner's collection
        self.collectionRef.deposit(token: <-token)
    }
}   
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import NonFungibleToken from 0xNFTADDRESS


// This transaction is for a user to stop a moment sale in their account
// by withdrawing that moment from their sale collection and depositing
// it into their normal moment collection

// Parameters
//
// tokenID: the ID of the moment whose sale is to be delisted

transaction(tokenID: UInt64) {

    let collectionRef: &TopShot.Collection
    let saleCollectionRef: auth(NonFungibleToken.Withdraw) &Market.SaleCollection

    prepare(acct: auth(Storage, Capabilities) &Account) {

        // Borrow a reference to the NFT collection in the signers account
        self.collectionRef = acct.storage.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // borrow a reference to the owner's sale collection
        self.saleCollectionRef = acct.storage.borrow<auth(NonFungibleToken.Withdraw) &Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
    }

    execute {
    
        // withdraw the moment from the sale, thereby de-listing it
        let token <- self.saleCollectionRef.withdraw(tokenID: tokenID)

        // deposit the moment into the owner's collection
        self.collectionRef.deposit(token: <-token)
    }
}   
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

// This transaction is for a user to stop a moment sale in their account
// by withdrawing that moment from their sale collection and depositing
// it into their normal moment collection

transaction(tokenID: UInt64) {

    prepare(acct: AuthAccount) {

        // Borrow a reference to the NFT collection in the signers account
        let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // borrow a reference to the owner's sale collection
        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        // withdraw the moment from the sale, thereby de-listing it
        let token <- topshotSaleCollection.withdraw(tokenID: tokenID)

        // deposit the moment into the owner's collection
        nftCollection.deposit(token: <-token)
    }
}
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

// This transaction is for a user to put a new moment up for sale
// They must have TopShot Collection and a Market Sale Collection
// stored in their account

transaction(momentID: UInt64, price: UFix64) {
    prepare(acct: AuthAccount) {

        // borrow a reference to the Top Shot Collection
        let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // withdraw the specified token from the collection
        let token <- nftCollection.withdraw(withdrawID: momentID) as! @TopShot.NFT

        // borrow a reference to the topshot Sale Collection
        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        // List the specified moment for sale
        topshotSaleCollection.listForSale(token: <-token, price: price)
    }
}
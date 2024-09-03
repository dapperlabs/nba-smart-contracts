import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This transaction is for a user to put a new moment up for sale
// They must have TopShot Collection and a Market Sale Collection
// stored in their account

// Parameters
//
// momentId: the ID of the moment to be listed for sale
// price: the sell price of the moment

transaction(momentID: UInt64, price: UFix64) {

    let collectionRef: auth(NonFungibleToken.Withdraw) &TopShot.Collection
    let saleCollectionRef: auth(Market.Create) &Market.SaleCollection

    prepare(acct: auth(Storage, Capabilities) &Account) {

        // borrow a reference to the Top Shot Collection
        self.collectionRef = acct.storage.borrow<auth(NonFungibleToken.Withdraw) &TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        // borrow a reference to the topshot Sale Collection
        self.saleCollectionRef = acct.storage.borrow<auth(Market.Create) &Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")
    }

    execute {

        // withdraw the specified token from the collection
        let token <- self.collectionRef.withdraw(withdrawID: momentID) as! @TopShot.NFT

        // List the specified moment for sale
        self.saleCollectionRef.listForSale(token: <-token, price: price)
    }
}
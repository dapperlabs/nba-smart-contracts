import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

transaction(momentID: UInt64, price: UFix64) {
    prepare(acct: AuthAccount) {
        let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        let token <- nftCollection.withdraw(withdrawID: momentID) as! @TopShot.NFT

        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.listForSale(token: <-token, price: price)
    }
}
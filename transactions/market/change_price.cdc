import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

transaction(tokenID: UInt64, newPrice: UFix64) {
    prepare(acct: AuthAccount) {

        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        topshotSaleCollection.changePrice(tokenID: tokenID, newPrice: newPrice)
    }
}
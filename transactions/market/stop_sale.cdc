import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

transaction(tokenID: UInt64) {
    prepare(acct: AuthAccount) {
        let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
            ?? panic("Could not borrow from sale in storage")

        let token <- topshotSaleCollection.withdraw(tokenID: tokenID)

        nftCollection.deposit(token: <-token)
    }
}
import TopShot from 0xTOPSHOTADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS
import NonFungibleToken from 0xNFTADDRESS

// This transaction locks a TopShot NFT rendering it unable to be withdrawn, sold, or transferred

// Parameters
//
// id: the Flow ID of the TopShot moment
// duration: number of seconds that the moment will be locked for

transaction(id: UInt64, duration: UFix64) {
    prepare(acct: auth(BorrowValue) &Account) {
        if let saleRef = acct.storage.borrow<auth(TopShotMarketV3.Cancel) &TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) {
            saleRef.cancelSale(tokenID: id)
        }

        let collectionRef = acct.storage.borrow<auth(NonFungibleToken.Update) &TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.lock(id: id, duration: duration)
    }
}

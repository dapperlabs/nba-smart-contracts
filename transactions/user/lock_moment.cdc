import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

transaction(id: UInt64, duration: UFix64) {
    prepare(acct: AuthAccount) {
        if let saleRef = acct.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) {
            saleRef.cancelSale(tokenID: id)
        }

        let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.lock(id: id, duration: duration)
    }
}

import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

// This transaction destroys a number of moments owned by a user

// Parameters
//
// momentIDs: an array of moment IDs of NFTs to be destroyed

transaction(momentIDs: [UInt64]) {

    let collectionRef: &TopShot.Collection
    
    prepare(acct: AuthAccount) {
        // delist any of the moments that are listed (this delists for both MarketV1 and Marketv3)
        if let topshotSaleV3Collection = acct.borrow<&TopShotMarketV3.SaleCollection>(from: TopShotMarketV3.marketStoragePath) {
            for id in momentIDs {
                if topshotSaleV3Collection.borrowMoment(id: id) != nil{
                    // cancel the moment from the sale, thereby de-listing it
                    topshotSaleV3Collection.cancelSale(tokenID: id)
                }
            }
        }

        self.collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")
    }

    execute {
        let tokens <- self.collectionRef.batchWithdraw(ids: momentIDs)
        // destroy the NFTs
        destroy tokens
    }
}
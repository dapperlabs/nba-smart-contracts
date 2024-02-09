

import FastBreakV1 from 0xFASTBREAKADDRESS

pub fun main(addr: Address): [UInt64] {

    let recipientAccount = getAccount(addr)
    let collectionRef = recipientAccount.getCapability(FastBreakV1.CollectionPublicPath).borrow<&{FastBreakV1.FastBreakNFTCollectionPublic}>()
        ?? panic("Could Not borrow Reference")

    var arrNFTs = collectionRef.getIDs()

    var scores = []

    for nftID in arrNFTs { 
        let nft = collectionRef.borrowFastBreakNFT(id: nftID) ??  
            panic("Couldn't borrow FastBreakNFT")

        scores.append(nft.points())
    }

    return scores
}


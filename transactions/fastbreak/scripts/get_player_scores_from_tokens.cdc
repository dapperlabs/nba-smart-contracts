

import FastBreakV1 from 0xFASTBREAKADDRESS

pub fun main(addr: Address): [{String: UInt64}] {

    let recipientAccount = getAccount(addr)
    let collectionRef = recipientAccount.getCapability(FastBreakV1.CollectionPublicPath).borrow<&{FastBreakV1.FastBreakNFTCollectionPublic}>()
        ?? panic("Could Not borrow Reference")

    var arrNFTs = collectionRef.getIDs()

    var scores: [{String: UInt64}] = []

    for nftID in arrNFTs { 
        let nft = collectionRef.borrowFastBreakNFT(id: nftID) ??  
            panic("Couldn't borrow FastBreakNFT")

        scores.append({nft.fastBreakGameID: nft.points()})
    }

    return scores
}


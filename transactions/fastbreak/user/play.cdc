import NonFungibleToken from 0xNFTADDRESS
import FastBreak from 0xFASTBREAKADDRESS

transaction(
    fastBreakGameID: String,
    topShots: [UInt64]
) {

    let gameRef: &FastBreak.Collection
    let recipient: &{FastBreak.FastBreakNFTCollectionPublic}

    prepare(acct: AuthAccount) {
        self.gameRef = acct
            .borrow<&FastBreak.Collection>(from: FastBreak.CollectionStoragePath)
            ?? panic("could not borrow a reference to the owner's collection")

        self.recipient = acct.getCapability(FastBreak.CollectionPublicPath)
            .borrow<&{FastBreak.FastBreakNFTCollectionPublic}>()
            ?? panic("could not borrow a reference to the collection receiver")
    }

    execute {

        let nft <- self.gameRef.play(
            fastBreakGameID: fastBreakGameID,
            topShots: topShots
        )
        self.recipient.deposit(token: <- (nft as @NonFungibleToken.NFT))
    }
}
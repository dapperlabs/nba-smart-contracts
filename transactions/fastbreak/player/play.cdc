import NonFungibleToken from 0xNFTADDRESS
import FastBreak from 0xFASTBREAKADDRESS

transaction(
    fastBreakGameID: String,
    topShots: [UInt64]
) {

    let gameRef: &FastBreak.Player
    let recipient: &{FastBreak.FastBreakNFTCollectionPublic}

    prepare(acct: AuthAccount) {
        self.gameRef = acct
            .borrow<&FastBreak.Player>(from: FastBreak.PlayerStoragePath)
            ?? panic("could not borrow a reference to the accounts player")

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
import NonFungibleToken from 0xNFTADDRESS
import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(
    fastBreakGameID: String,
    topShots: [UInt64]
) {

    let gameRef: &FastBreakV1.Player
    let recipient: &{FastBreakV1.FastBreakNFTCollectionPublic}

    prepare(acct: AuthAccount) {
        self.gameRef = acct
            .borrow<&FastBreakV1.Player>(from: FastBreakV1.PlayerStoragePath)
            ?? panic("could not borrow a reference to the accounts player")

        self.recipient = acct.getCapability(FastBreakV1.CollectionPublicPath)
            .borrow<&{FastBreakV1.FastBreakNFTCollectionPublic}>()
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
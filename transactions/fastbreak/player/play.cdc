import NonFungibleToken from 0xNFTADDRESS
import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(
    fastBreakGameID: String,
    topShots: [UInt64]
) {

    let gameRef: auth(FastBreakV1.Play) &FastBreakV1.Player
    let recipient: &{FastBreakV1.FastBreakNFTCollectionPublic}

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.gameRef = acct.storage
            .borrow<auth(FastBreakV1.Play) &FastBreakV1.Player>(from: FastBreakV1.PlayerStoragePath)
            ?? panic("could not borrow a reference to the accounts player")

        self.recipient = acct.capabilities.borrow<&FastBreakV1.Collection>(FastBreakV1.CollectionPublicPath)
            ?? panic("could not borrow a reference to the collection receiver")

    }

    execute {

        let nft <- self.gameRef.play(
            fastBreakGameID: fastBreakGameID,
            topShots: topShots
        )
        self.recipient.deposit(token: <- (nft as @{NonFungibleToken.NFT}))
    }
}
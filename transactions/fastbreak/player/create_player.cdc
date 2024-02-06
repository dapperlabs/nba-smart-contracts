import NonFungibleToken from 0xNFTADDRESS
import FastBreak from 0xFASTBREAKADDRESS

transaction(playerName: String) {

    prepare(signer: AuthAccount) {
        if signer.borrow<&FastBreak.Collection>(from: FastBreak.CollectionStoragePath) == nil {

            let collection <- FastBreak.createEmptyCollection()
            signer.save(<-collection, to: FastBreak.CollectionStoragePath)
            signer.link<&FastBreak.Collection{NonFungibleToken.CollectionPublic, FastBreak.FastBreakNFTCollectionPublic}>(
                FastBreak.CollectionPublicPath,
                target: FastBreak.CollectionStoragePath
            )

        }

        if signer.borrow<&FastBreak.Player>(from: FastBreak.PlayerStoragePath) == nil {

            let player <- FastBreak.createPlayer(playerName: playerName)
            signer.save(<-player, to: FastBreak.PlayerStoragePath)
            signer.link<&FastBreak.Player{FastBreak.FastBreakPlayer}>(
                FastBreak.PlayerPrivatePath,
                target: FastBreak.PlayerStoragePath
            )
        }
    }
}
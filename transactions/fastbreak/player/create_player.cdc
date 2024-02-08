import NonFungibleToken from 0xNFTADDRESS
import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(playerName: String) {

    prepare(signer: AuthAccount) {
        if signer.borrow<&FastBreakV1.Collection>(from: FastBreakV1.CollectionStoragePath) == nil {

            let collection <- FastBreakV1.createEmptyCollection()
            signer.save(<-collection, to: FastBreakV1.CollectionStoragePath)
            signer.link<&FastBreakV1.Collection{NonFungibleToken.CollectionPublic, FastBreakV1.FastBreakNFTCollectionPublic}>(
                FastBreakV1.CollectionPublicPath,
                target: FastBreakV1.CollectionStoragePath
            )

        }

        if signer.borrow<&FastBreakV1.Player>(from: FastBreakV1.PlayerStoragePath) == nil {

            let player <- FastBreakV1.createPlayer(playerName: playerName)
            signer.save(<-player, to: FastBreakV1.PlayerStoragePath)
            signer.link<&FastBreakV1.Player{FastBreakV1.FastBreakPlayer}>(
                FastBreakV1.PlayerPrivatePath,
                target: FastBreakV1.PlayerStoragePath
            )
        }
    }
}
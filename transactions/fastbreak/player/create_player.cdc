import NonFungibleToken from 0xNFTADDRESS
import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(playerName: String) {

    prepare(signer: auth(Storage, Capabilities) &Account) {
        if signer.storage.borrow<&FastBreakV1.Collection>(from: FastBreakV1.CollectionStoragePath) == nil {

            let collection <- FastBreakV1.createEmptyCollection(nftType: Type<@FastBreakV1.NFT>())
            signer.storage.save(<-collection, to: FastBreakV1.CollectionStoragePath)
            signer.capabilities.unpublish(FastBreakV1.CollectionPublicPath)
            signer.capabilities.publish(
                signer.capabilities.storage.issue<&FastBreakV1.Collection>(FastBreakV1.CollectionStoragePath),
                at: FastBreakV1.CollectionPublicPath
            )

        }

        if signer.storage.borrow<&FastBreakV1.Player>(from: FastBreakV1.PlayerStoragePath) == nil {

            let player <- FastBreakV1.createPlayer(playerName: playerName)
            signer.storage.save(<-player, to: FastBreakV1.PlayerStoragePath)
        }
    }
}
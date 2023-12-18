import NonFungibleToken from 0xNFTADDRESS
import FastBreak from 0xFASTBREAKADDRESS

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&FastBreak.Collection>(from: FastBreak.CollectionStoragePath) == nil {

            let collection <- FastBreak.createEmptyCollection()
            signer.save(<-collection, to: FastBreak.CollectionStoragePath)
            signer.link<&FastBreak.Collection{NonFungibleToken.CollectionPublic, FastBreak.FastBreakNFTCollectionPublic}>(
                FastBreak.CollectionPublicPath,
                target: FastBreak.CollectionStoragePath
            )

        }
    }
}
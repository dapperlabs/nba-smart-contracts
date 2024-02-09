import FastBreakV1 from 0xFASTBREAKADDRESS
import NonFungibleToken from 0xNFTADDRESS

transaction(fastBreakGameID: String, topShotMomentIds: [UInt64]) {

    let gameRef: &FastBreakV1.Player

    prepare(acct: AuthAccount) {

        self.gameRef = acct
            .borrow<&FastBreakV1.Player>(from: FastBreakV1.PlayerStoragePath)
            ?? panic("could not borrow a reference to the accounts player")

    }

    execute {
        self.gameRef.updateSubmission(fastBreakGameID: fastBreakGameID, topShots: topShotMomentIds)
    }
}
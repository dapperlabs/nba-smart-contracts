import FastBreakV1 from 0xFASTBREAKADDRESS
import NonFungibleToken from 0xNFTADDRESS

transaction(fastBreakGameID: String, topShotMomentIds: [UInt64]) {

    let gameRef: auth(FastBreakV1.Update) &FastBreakV1.Player

    prepare(acct: auth(BorrowValue) &Account) {

        self.gameRef = acct.storage
            .borrow<auth(FastBreakV1.Update) &FastBreakV1.Player>(from: FastBreakV1.PlayerStoragePath)
            ?? panic("could not borrow a reference to the accounts player")

    }

    execute {
        self.gameRef.updateSubmission(fastBreakGameID: fastBreakGameID, topShots: topShotMomentIds)
    }
}
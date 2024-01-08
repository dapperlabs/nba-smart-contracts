import FastBreak from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, playerId: UInt64, points: UInt64, win: Bool) {

    let oracleRef: &FastBreak.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.updateFastBreakScore(
            fastBreakGameID: fastBreakGameID,
            playerId: playerId,
            points: points,
            win: win
        )
    }

}
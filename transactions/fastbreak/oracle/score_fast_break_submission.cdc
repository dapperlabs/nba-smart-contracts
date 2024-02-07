import FastBreak from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, playerAddress: Address, points: UInt64, win: Bool) {

    let oracleRef: &FastBreak.FastBreakDaemon
    let playerId: UInt64

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("could not borrow a reference to the oracle resource")

        self.playerId = FastBreak.getPlayerIdByAccount(accountAddress: playerAddress)
    }

    execute {

        self.oracleRef.updateFastBreakScore(
            fastBreakGameID: fastBreakGameID,
            playerId: self.playerId,
            points: points,
            win: win
        )
    }

}
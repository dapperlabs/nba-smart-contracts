import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, playerAddress: Address, points: UInt64, win: Bool) {

    let oracleRef: &FastBreakV1.FastBreakDaemon
    let playerId: UInt64

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
            ?? panic("could not borrow a reference to the oracle resource")

        self.playerId = FastBreakV1.getPlayerIdByAccount(accountAddress: playerAddress)
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
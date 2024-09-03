import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, playerAddress: Address, points: UInt64, win: Bool) {

    let oracleRef: auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon
    let playerId: UInt64

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.oracleRef = acct.storage.borrow<auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
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
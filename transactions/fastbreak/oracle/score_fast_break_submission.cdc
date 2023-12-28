import FastBreak from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, accountAddress: Address, points: UInt64, win: Bool) {

    let oracleRef: &FastBreak.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.updateFastBreakScore(
            fastBreakGameID: fastBreakGameID,
            accountAddress: accountAddress,
            points: points,
            win: win
        )
    }

}
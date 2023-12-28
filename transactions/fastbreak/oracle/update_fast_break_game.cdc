import FastBreak from 0xFASTBREAKADDRESS

transaction(id: String, status: UInt8, winner: Address?) {

    let oracleRef: &FastBreak.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.updateFastBreakGame(
            id: id,
            status: status,
            winner: winner
        )
    }
}
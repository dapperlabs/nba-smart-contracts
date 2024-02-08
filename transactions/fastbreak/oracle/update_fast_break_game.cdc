import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(id: String, status: UInt8, winner: UInt64) {

    let oracleRef: &FastBreakV1.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
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
import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(
    id: String,
    name: String,
    fastBreakRunID: String,
    submissionDeadline: UInt64,
    numPlayers: UInt64
) {

    let oracleRef: &FastBreakV1.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.createFastBreakGame(
            id: id,
            name: name,
            fastBreakRunID: fastBreakRunID,
            submissionDeadline: submissionDeadline,
            numPlayers: numPlayers
        )
    }
}
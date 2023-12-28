import FastBreak from 0xFASTBREAKADDRESS

transaction(
    id: String,
    name: String,
    fastBreakRunID: String,
    isPublic: Bool,
    submissionDeadline: UInt64,
    numPlayers: UInt64
) {

    let oracleRef: &FastBreak.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.createFastBreakGame(
            id: id,
            name: name,
            fastBreakRunID: fastBreakRunID,
            isPublic: isPublic,
            submissionDeadline: submissionDeadline,
            numPlayers: numPlayers
        )
    }
}
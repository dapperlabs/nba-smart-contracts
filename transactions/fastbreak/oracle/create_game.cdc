import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(
    id: String,
    name: String,
    fastBreakRunID: String,
    submissionDeadline: UInt64,
    numPlayers: UInt64
) {

    let oracleRef: auth(FastBreakV1.Create) &FastBreakV1.FastBreakDaemon

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.oracleRef = acct.storage.borrow<auth(FastBreakV1.Create) &FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
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
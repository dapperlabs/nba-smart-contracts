import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, submissionDeadline: UInt64) {

    let oracleRef: auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.oracleRef = acct.storage.borrow<auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.updateSubmissionDeadline(
            fastBreakGameID: fastBreakGameID,
            submissionDeadline: submissionDeadline
        )
    }
}
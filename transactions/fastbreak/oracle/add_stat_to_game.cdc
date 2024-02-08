import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64) {

    let oracleRef: &FastBreakV1.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.addStatToFastBreakGame(
            fastBreakGameID: fastBreakGameID,
            name: name,
            rawType: rawType,
            valueNeeded: valueNeeded
        )
    }
}
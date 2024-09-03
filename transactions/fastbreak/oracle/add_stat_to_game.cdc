import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, name: String, rawType: UInt8, valueNeeded: UInt64) {

    let oracleRef: auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.oracleRef = acct.storage.borrow<auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
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
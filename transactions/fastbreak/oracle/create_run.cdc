import FastBreakV1 from 0xFASTBREAKADDRESS


transaction(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {

    let oracleRef: auth(FastBreakV1.Create) &FastBreakV1.FastBreakDaemon

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.oracleRef = acct.storage.borrow<auth(FastBreakV1.Create) &FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {
        self.oracleRef.createFastBreakRun(id: id, name: name, runStart: runStart, runEnd: runEnd, fatigueModeOn: fatigueModeOn)
    }

    post {
        FastBreakV1.getFastBreakRun(id: id)?.name! == name: "could not find fast break run"
    }
}
import FastBreak from 0xFASTBREAKADDRESS


transaction(id: String, name: String, runStart: UInt64, runEnd: UInt64, fatigueModeOn: Bool) {

    let oracleRef: &FastBreak.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {
        self.oracleRef.createFastBreakRun(id: id, name: name, runStart: runStart, runEnd: runEnd, fatigueModeOn: fatigueModeOn)
    }

    post {
        FastBreak.getFastBreakRun(id: id)?.name! == name: "could not find fast break run"
    }
}
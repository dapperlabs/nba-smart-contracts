import FastBreak from 0xFASTBREAKADDRESS

transaction(fastBreakGameID: String, name: String, type: String, valueNeeded: UInt64) {

    let oracleRef: &FastBreak.FastBreakDaemon

    prepare(acct: AuthAccount) {
        self.oracleRef = acct.borrow<&FastBreak.FastBreakDaemon>(from: FastBreak.OracleStoragePath)
            ?? panic("Could not borrow a reference to the oracle resource")
    }

    execute {

        self.oracleRef.addStatToFastBreakGame(
            fastBreakGameID: fastBreakGameID,
            name: name,
            type: type,
            valueNeeded: valueNeeded
        )
    }

    post {
        FastBreak.getFastBreakGameStats(id: fastBreakGameID).length != 0: "could not find fast break stat"
    }
}
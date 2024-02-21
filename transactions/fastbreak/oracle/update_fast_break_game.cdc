import FastBreakV1 from 0xFASTBREAKADDRESS

transaction(id: String, status: UInt8, winner: UInt64) {

    let oracleRef: auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon

    prepare(acct: auth(Storage, Capabilities) &Account) {
        self.oracleRef = acct.storage.borrow<auth(FastBreakV1.Update) &FastBreakV1.FastBreakDaemon>(from: FastBreakV1.OracleStoragePath)
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
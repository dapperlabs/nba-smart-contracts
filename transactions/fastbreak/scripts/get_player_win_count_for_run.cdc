import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(runId: String, playerAddress: Address): UInt64 {
    let playerId = FastBreakV1.getPlayerIdByAccount(accountAddress: playerAddress)
    if let run = FastBreakV1.getFastBreakRun(id: runId) {
        return run.runWinCount[playerId] ?? 0
    }
    return 0
}
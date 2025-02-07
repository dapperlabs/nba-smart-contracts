import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(runId: String, playerAddress: Address): UInt64 {
    let playerId = FastBreakV1.getPlayerIdByAccount(accountAddress: playerAddress)
    let fastBreakRun = FastBreakV1.getFastBreakRun(id: runId)!

    return fastBreakRun.runWinCount[playerId] ?? 0
}
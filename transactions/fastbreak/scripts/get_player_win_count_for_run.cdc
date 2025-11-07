import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(runId: String, playerAddress: Address): UInt64 {
    let playerId = FastBreakV1.getPlayerIdByAccount(accountAddress: playerAddress)
    // New runs are stored in year-based storage
    // Try current year first, then previous year, then legacy
    let currentTimestamp = UInt64(getCurrentBlock().timestamp)
    let currentYearNum = FastBreakV1.getYearFromTimestamp(timestamp: currentTimestamp)
    let currentYearString = currentYearNum.toString()
    
    var fastBreakRun: &FastBreakV1.FastBreakRun? = nil
    
    // Check current year
    if let run = FastBreakV1.getFastBreakRunByYear(id: runId, year: currentYearString) {
        fastBreakRun = run
    } else if currentYearNum > 1970 {
        // Check previous year
        let previousYearString = (currentYearNum - 1).toString()
        fastBreakRun = FastBreakV1.getFastBreakRunByYear(id: runId, year: previousYearString)
    }
    
    // Fallback to legacy (deprecated)
    if fastBreakRun == nil {
        fastBreakRun = FastBreakV1.getFastBreakRun(id: runId)
    }
    
    if let run = fastBreakRun {
        return run.runWinCount[playerId] ?? 0
    }
    return 0
}
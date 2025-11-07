import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(id: String, playerAddress: Address): UInt64 {
    let playerId = FastBreakV1.getPlayerIdByAccount(accountAddress: playerAddress)
    // New games are stored in year-based storage
    // Try current year first, then previous year, then legacy
    let currentTimestamp = UInt64(getCurrentBlock().timestamp)
    let currentYearNum = FastBreakV1.getYearFromTimestamp(timestamp: currentTimestamp)
    let currentYearString = currentYearNum.toString()
    
    var fastBreak: &FastBreakV1.FastBreakGame? = nil
    
    // Check current year
    if let game = FastBreakV1.getFastBreakGameByYear(id: id, year: currentYearString) {
        fastBreak = game
    } else if currentYearNum > 1970 {
        // Check previous year
        let previousYearString = (currentYearNum - 1).toString()
        fastBreak = FastBreakV1.getFastBreakGameByYear(id: id, year: previousYearString)
    }
    
    // Fallback to legacy (deprecated)
    if fastBreak == nil {
        fastBreak = FastBreakV1.getFastBreakGame(id: id)
    }
    
    let submission = fastBreak?.getFastBreakSubmissionByPlayerId(playerId: playerId)!

    return submission?.points ?? 0
}
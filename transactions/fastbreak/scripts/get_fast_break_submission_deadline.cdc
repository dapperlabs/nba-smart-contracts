import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(id: String): UInt64? {
    // New games are stored in year-based storage
    // Try current year first, then previous year, then legacy
    let currentTimestamp = UInt64(getCurrentBlock().timestamp)
    let currentYearNum = FastBreakV1.getYearFromTimestamp(timestamp: currentTimestamp)
    let currentYearString = currentYearNum.toString()
    
    // Check current year
    if let game = FastBreakV1.getFastBreakGameByYear(id: id, year: currentYearString) {
        return game.submissionDeadline
    }
    
    // Check previous year
    if currentYearNum > 1970 {
        let previousYearString = (currentYearNum - 1).toString()
        if let game = FastBreakV1.getFastBreakGameByYear(id: id, year: previousYearString) {
            return game.submissionDeadline
        }
    }
    
    // Fallback to legacy (deprecated)
    return FastBreakV1.getFastBreakGame(id: id)?.submissionDeadline
}
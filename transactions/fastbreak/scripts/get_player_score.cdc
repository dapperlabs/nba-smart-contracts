import FastBreak from 0xFASTBREAKADDRESS

pub fun main(id: String, playerAddress: Address): UInt64 {
    let playerId = FastBreak.getPlayerIdByAccount(accountAddress: playerAddress)
    let fastBreak = FastBreak.getFastBreakGame(id: id)!
    let submission : FastBreak.FastBreakSubmission = fastBreak.getFastBreakSubmissionByPlayerId(playerId: playerId)!

    return submission.points
}
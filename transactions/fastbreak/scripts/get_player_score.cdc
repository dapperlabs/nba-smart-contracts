import FastBreak from 0xFASTBREAKADDRESS

pub fun main(id: String, playerId: UInt64): UInt64 {
    let fastBreak = FastBreak.getFastBreakGame(id: id)!
    let submission : FastBreak.FastBreakSubmission = fastBreak.getFastBreakSubmissionByPlayerId(playerId: playerId)!

    return submission.points
}
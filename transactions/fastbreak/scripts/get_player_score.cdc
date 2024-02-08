import FastBreakV1 from 0xFASTBREAKADDRESS

pub fun main(id: String, playerAddress: Address): UInt64 {
    let playerId = FastBreakV1.getPlayerIdByAccount(accountAddress: playerAddress)
    let fastBreak = FastBreakV1.getFastBreakGame(id: id)!
    let submission : FastBreakV1.FastBreakSubmission = fastBreak.getFastBreakSubmissionByPlayerId(playerId: playerId)!

    return submission.points
}
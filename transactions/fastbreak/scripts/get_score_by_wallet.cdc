import FastBreak from 0xFASTBREAKADDRESS

pub fun main(id: String, accountAddress: Address): UInt64 {
    let fastBreak = FastBreak.getFastBreakGame(id: id)!
    let submission : FastBreak.FastBreakSubmission = fastBreak.getFastBreakSubmissionByAccount(accountAddress: accountAddress)!

    return submission.points
}
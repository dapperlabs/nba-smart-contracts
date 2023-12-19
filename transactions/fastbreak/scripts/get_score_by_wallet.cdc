import FastBreak from 0xFASTBREAKADDRESS

pub fun main(id: String, wallet: Address): UInt64 {
    let fastBreak = FastBreak.getFastBreakGame(id: id)!
    let submission : FastBreak.FastBreakSubmission = fastBreak.getFastBreakSubmissionByWallet(wallet: wallet)!

    return submission.points
}
import FastBreak from 0xFASTBREAKADDRESS

pub fun main(fastBreakGameID: String): [FastBreak.FastBreakStat] {

    return FastBreak.getFastBreakGameStats(id: fastBreakGameID)
}
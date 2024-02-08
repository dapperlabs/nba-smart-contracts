import FastBreakV1 from 0xFASTBREAKADDRESS

pub fun main(fastBreakGameID: String): [FastBreakV1.FastBreakStat] {

    return FastBreakV1.getFastBreakGameStats(id: fastBreakGameID)
}
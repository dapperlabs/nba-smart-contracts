import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(fastBreakGameID: String): [FastBreakV1.FastBreakStat] {

    return FastBreakV1.getFastBreakGameStats(id: fastBreakGameID)
}
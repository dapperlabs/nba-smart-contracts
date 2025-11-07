import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(id: String, year: String): &FastBreakV1.FastBreakGame? {
    return FastBreakV1.getFastBreakGameByYear(id: id, year: year)
}


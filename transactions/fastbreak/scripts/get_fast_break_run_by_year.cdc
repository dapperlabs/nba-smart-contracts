import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(id: String, year: String): &FastBreakV1.FastBreakRun? {
    return FastBreakV1.getFastBreakRunByYear(id: id, year: year)
}


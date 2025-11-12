import FastBreakV1 from 0xFASTBREAKADDRESS

access(all) fun main(id: String): UInt64? {
    return FastBreakV1.getFastBreakGame(id: id)?.submissionDeadline
}
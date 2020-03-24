import TopShot from 0x03

// This script reads the next Set ID from the TopShot contract and 
// returns that number to the caller
pub fun main(setID: UInt32): UInt32 {
    log(TopShot.setDatas[setID]!.series)
    return TopShot.setDatas[setID]!.series
}






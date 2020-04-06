import TopShot from 0x03

// This script reads the next Set ID from the TopShot contract and 
// returns that number to the caller
pub fun main(): UInt32 {
    log(TopShot.setDatas[UInt32(0)]!.series)
    return TopShot.setDatas[UInt32(0)]!.series
}






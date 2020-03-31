import TopShot from 0x03

// This script reads the next Set ID from the TopShot contract and 
// returns that number to the caller
pub fun main(): String {
    log(TopShot.setDatas[UInt32(0)]!.name)
    return TopShot.setDatas[UInt32(0)]!.name
}






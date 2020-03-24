import TopShot from 0x03

// This script reads the next Set ID from the TopShot contract and 
// returns that number to the caller
pub fun main(setID: UInt32): String {
    log(TopShot.setDatas[setID]!.name)
    return TopShot.setDatas[setID]!.name
}






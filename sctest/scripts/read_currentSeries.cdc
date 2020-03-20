import TopShot from 0x03

// This script reads the current series from the TopShot contract and 
// returns that number to the caller
pub fun main(): UInt32 {
    log(TopShot.currentSeries)
    return TopShot.currentSeries
}






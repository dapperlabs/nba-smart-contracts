import TopShot from 0xTOPSHOTADDRESS

// This script reads the current series from the TopShot contract and 
// returns that number to the caller

// Returns: UInt32
// currentSeries field in TopShot contract

access(all) fun main(): UInt32 {

    return TopShot.currentSeries
}
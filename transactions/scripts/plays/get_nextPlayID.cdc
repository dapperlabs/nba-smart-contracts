import TopShot from 0xTOPSHOTADDRESS

// This script reads the public nextPlayID from the TopShot contract and 
// returns that number to the caller

// Returns: UInt32
// the nextPlayID field in TopShot contract

access(all) fun main(): UInt32 {

    log(TopShot.nextPlayID)

    return TopShot.nextPlayID
}
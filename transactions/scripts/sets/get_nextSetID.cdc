import TopShot from 0xTOPSHOTADDRESS

// This script reads the next Set ID from the TopShot contract and 
// returns that number to the caller

// Returns: UInt32
// Value of nextSetID field in TopShot contract

access(all) fun main(): UInt32 {

    log(TopShot.nextSetID)

    return TopShot.nextSetID
}
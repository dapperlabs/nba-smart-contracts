import TopShot from 0xTOPSHOTADDRESS

// This script reads the current number of moments that have been minted
// from the TopShot contract and returns that number to the caller

pub fun main(): UInt64 {
    return TopShot.totalSupply
}
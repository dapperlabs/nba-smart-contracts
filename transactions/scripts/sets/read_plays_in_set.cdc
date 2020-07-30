import TopShot from 0xTOPSHOTADDRESS

// This script returns an array of the play IDs that are
// in the specified set

pub fun main(setID: UInt32): [UInt32] {
    let plays = TopShot.getPlaysInSet(setID: setID)?
        ?? panic("Could not find the specified set")

    return plays
}
import TopShot from 0xTOPSHOTADDRESS

// This script returns an array of the setIDs
// that have the specified name

// Parameters:
//
// setName: The name of the set whose data needs to be read

// Returns: [UInt32]
// Array of setIDs that have specified set name

pub fun main(setName: String): [UInt32] {

    let ids = TopShot.getSetIDsByName(setName: setName)
        ?? panic("Could not find the specified set name")

    return ids
}
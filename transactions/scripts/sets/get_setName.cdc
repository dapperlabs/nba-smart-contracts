import TopShot from 0xTOPSHOTADDRESS

// This script gets the setName of a set with specified setID

// Parameters:
//
// setID: The unique ID for the set whose data needs to be read

// Returns: String
// Name of set with specified setID

pub fun main(setID: UInt32): String {

    let name = TopShot.getSetName(setID: setID)
        ?? panic("Could not find the specified set")
        
    return name
}
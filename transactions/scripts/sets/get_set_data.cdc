import TopShot from 0xTOPSHOTADDRESS

// This script returns all the metadata about the specified set

// Parameters:
//
// setID: The unique ID for the set whose data needs to be read

// Returns: TopShot.QuerySetData

pub fun main(setID: UInt32): TopShot.QuerySetData {

    let data = TopShot.getSetData(setID: setID)
        ?? panic("Could not get data for the specified set ID")

    return data
}
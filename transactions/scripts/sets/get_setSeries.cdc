import TopShot from 0xTOPSHOTADDRESS

// This script reads the series of the specified set and returns it

// Parameters:
//
// setID: The unique ID for the set whose data needs to be read

// Returns: UInt32
// unique ID of series

access(all) fun main(setID: UInt32): UInt32 {

    let series = TopShot.getSetSeries(setID: setID)
        ?? panic("Could not find the specified set")

    return series
}
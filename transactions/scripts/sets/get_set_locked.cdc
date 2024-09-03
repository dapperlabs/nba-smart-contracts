import TopShot from 0xTOPSHOTADDRESS

// This script returns a boolean indicating if the specified set is locked
// meaning new plays cannot be added to it

// Parameters:
//
// setID: The unique ID for the set whose data needs to be read

// Returns: Bool
// Whether specified set is locked

access(all) fun main(setID: UInt32): Bool {

    let isLocked = TopShot.isSetLocked(setID: setID)
        ?? panic("Could not find the specified set")

    return isLocked
}
import TopShot from 0xTOPSHOTADDRESS

// This script returns a boolean indicating if the specified set is locked
// meaning new plays cannot be added to it

pub fun main(setID: UInt32): Bool {
    let isLocked = TopShot.isSetLocked(setID: setID)
        ?? panic("Could not find the specified set")

    return isLocked
}
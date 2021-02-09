import TopShot from 0xTOPSHOTADDRESS

// this script returns an array of the setIDs
// that have the specified name

pub fun main(setName: String): [UInt32] {
    let ids = TopShot.getSetIDsByName(setName: setName)
        ?? panic("Could not find the specified set name")

    return ids
}
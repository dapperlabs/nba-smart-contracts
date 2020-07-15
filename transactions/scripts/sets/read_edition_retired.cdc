import TopShot from 0xTOPSHOTADDRESS

// This transaction reads if a specified edition is retired

pub fun main(setID: UInt32, playID: UInt32): Bool {
    let isRetired = TopShot.isEditionRetired(setID: setID, playID: playID)
        ?? panic("Could not find the specified edition")
    
    return isRetired
}
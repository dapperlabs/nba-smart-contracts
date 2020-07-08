import TopShot from 0xTOPSHOTADDRESS

pub fun main(setID: UInt32, playID: UInt32): Bool {
    let isRetired = TopShot.isEditionRetired(setID: setID, playID: playID)!
    
    return isRetired
}
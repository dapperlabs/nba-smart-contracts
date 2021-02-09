import TopShot from 0xTOPSHOTADDRESS

// This script returns the number of moments that have been
// minted for the specified edition

pub fun main(setID: UInt32, playID: UInt32): UInt32 {
    let numMoments = TopShot.getNumMomentsInEdition(setID: setID, playID: playID)
        ?? panic("Could not find the specified edition")

    return numMoments
}
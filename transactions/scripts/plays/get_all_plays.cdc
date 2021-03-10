import TopShot from 0xTOPSHOTADDRESS

// This script returns an array of all the plays 
// that have ever been created for Top Shot

// Returns: [TopShot.Play]
// array of all plays created for Topshot

pub fun main(): [TopShot.Play] {

    return TopShot.getAllPlays()
}
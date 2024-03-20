import TopShot from 0xTOPSHOTADDRESS

// This script returns the full metadata associated with a play
// in the TopShot smart contract

// Parameters:
//
// playID: The unique ID for the play whose data needs to be read

// Returns: {String:String}
// A dictionary of all the play metadata associated
// with the specified playID

access(all) fun main(playID: UInt32): {String:String} {

    let metadata = TopShot.getPlayMetaData(playID: playID) ?? panic("Play doesn't exist")

    log(metadata)

    return metadata
}
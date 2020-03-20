import TopShot from 0x03

// This script returns the full metadata associated with a play
// in the TopShot smart contract
//
pub fun main(playID: UInt32): {String:String} {
    let metadata = TopShot.getPlayMetaData(playID: playID) ?? panic("Play doesn't exist")
    return metadata
}






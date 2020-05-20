import TopShot from 0x03

// This script returns the full metadata associated with a play
// in the TopShot smart contract
//
pub fun main(): {String:String} {
    let metadata = TopShot.getPlayMetaData(playID: 0) ?? panic("Play doesn't exist")
    log(metadata)
    return metadata
}






import TopShot from 0x03

// This script returns the full metadata associated with a play
// in the TopShot smart contract
//
pub fun main(playID: UInt32, field: String): String {
    let field = TopShot.getPlayMetaDataByField(playID: playID, field: field) ?? panic("Play doesn't exist")

    return field
}






import TopShot from 0x03

// This script returns the full metadata associated with a play
// in the TopShot smart contract
//
pub fun main(): String {
    let field = TopShot.getPlayMetaDataByField(playID: 0, field: "Name") ?? panic("Play doesn't exist")

    log(field)

    return field
}






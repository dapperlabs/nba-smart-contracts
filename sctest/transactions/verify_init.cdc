import TopShot from 0x02

// this script checks to see that the IDs are a certain number
// feel free to change them to execute
pub fun main(): UInt32 {
    log("Mold ID")
    log(TopShot.moldID)
    log("Moment ID")
    log(TopShot.totalSupply)
    if TopShot.totalSupply != UInt64(0) && TopShot.moldID != UInt32(0) {
        panic("Wrong initialization!")
    }
    return TopShot.moldID
}
import TopShot from 0x02

// this script checks to see that the IDs are a certain number
// feel free to change them to execute
pub fun main() {
    log("Mold ID")
    log(TopShot.moldID)
    log("Moment ID")
    log(TopShot.totalSupply)
    // if TopShot.totalSupply != UInt64(1) && TopShot.moldID != UInt64(1) {
    //     panic("Wrong initialization!")
    // } 
}
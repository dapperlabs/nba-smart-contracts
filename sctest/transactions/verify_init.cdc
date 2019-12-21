import TopShot from 0x01

// this script checks to see that the IDs are a certain number
// feel free to change them to execute
pub fun main() {
    log("Mold ID")
    log(TopShot.moldID)
    log("Moment ID")
    log(TopShot.momentID)
    if TopShot.momentID != 1 && TopShot.moldID != 1 {
        panic("Wrong initialization!")
    } 
}
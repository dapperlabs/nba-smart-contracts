import TopShot from 0x02

// this script checks to see that the IDs are a certain number
// feel free to change them to execute
pub fun main() {
    let acct = getAccount(0x01)
    let receiver = acct.published[&TopShot.MomentCollectionPublic] ?? panic("missing ref")

    log(receiver.getIDs())

    // if let field = receiver.getMomentMetadataField(id: 1, field: "Name") {
    //     log(field)
    // }
}
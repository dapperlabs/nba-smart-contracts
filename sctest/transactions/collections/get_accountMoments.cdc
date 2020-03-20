import TopShot from 0x02

pub fun main() {

    let acct = getAccount(0x02)

    if let collectionRef = acct.published[&TopShot.MomentCollectionPublic] {

        let collectionIDs = collectionRef.getIDs()

        log(collectionRef.getIDs)

    } else {
        log("No collection reference!")
    }
}
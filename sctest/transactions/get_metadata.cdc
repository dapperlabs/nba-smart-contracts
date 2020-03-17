import TopShot from 0x02

pub fun main() {
    logMoldMetaData(moldID: 0)

    logMomentMetaData(address: 0x02, momentID: 1)
}

// This is how you would log the metadata of a mold
// if you already knew the id of the mold you were looking for
pub fun logMoldMetaData(moldID: UInt32) {
    if let mold = TopShot.plays[UInt32(1)] {
        log(mold.metadata)
    }
}

// This is how you would log the metadata of a moment
// if you don't already know what the mold ID is
pub fun logMomentMetaData(address: Address, momentID: UInt64) {
    let acct = getAccount(address)
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getMetaData(id: momentID))
}
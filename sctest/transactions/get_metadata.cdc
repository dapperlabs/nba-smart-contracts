import TopShot from 0x02

access(all) fun main() {
    logMoldMetaData(moldID: 0)

    logMomentMetaData(address: 0x02, momentID: 1)

    logMomentQuality(address: 0x02, momentID: 1)

    logMomentPlace(address:0x02, momentID: 1)
}

// This is how you would log the metadata of a mold
// if you already knew the id of the mold you were looking for
access(all) fun logMoldMetaData(moldID: UInt32) {
    if let mold = TopShot.molds[UInt32(1)] {
        log(mold.metadata)
    }
}

// This is how you would log the metadata of a moment
// if you don't already know what the mold ID is
access(all) fun logMomentMetaData(address: Address, momentID: UInt64) {
    let acct = getAccount(address)
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getMetaData(id: momentID))
}

access(all) fun logMomentQuality(address: Address, momentID: UInt64) {
    let acct = getAccount(address)
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getQuality(id: momentID))
}

access(all) fun logMomentPlace(address: Address, momentID: UInt64) {
    let acct = getAccount(address)
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getPlaceInQuality(id: momentID))
}
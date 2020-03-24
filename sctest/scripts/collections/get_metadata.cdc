import TopShot from 0x03

pub fun main(address: Address, momentID: UInt64): {String:String} {

    let acct = getAccount(address)

    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getMetaData(id: momentID))

    return collectionRef.getMetaData(id: momentID)!
}
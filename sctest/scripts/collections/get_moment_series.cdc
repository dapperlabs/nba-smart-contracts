import TopShot from 0x03

pub fun main(address: Address, momentID: UInt64): UInt32 {

    let acct = getAccount(address)

    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getSetID(id: momentID))

    return collectionRef.getSetID(id: momentID)!
}
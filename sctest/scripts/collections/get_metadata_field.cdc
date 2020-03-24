import TopShot from 0x03

pub fun main(address: Address, momentID: UInt64, field: String): String {

    let acct = getAccount(address)

    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    let metadata = collectionRef.getMetaData(id: momentID)!

    log(metadata[field])

    return metadata[field]!
}
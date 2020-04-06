import TopShot from 0x03

// This transaction gets the metadata associated with a moment
// in a collection by looking up its playID and then searching
// for that play's metadata in the TopShot contract

pub fun main(address: Address, momentID: UInt64, field: String): String {

    // get the Address of the account with the Moment
    let acct = getAccount(address)

    // Get that account's published collectionRef
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    // Get a reference to a specific NFT in the collection
    let ref = collectionRef.borrowNFT(id: momentID)

    // Get the metadata from the play
    let metadata = TopShot.playDatas[ref.data.playID]!.metadata

    log(metadata[field])

    return metadata[field]!
}
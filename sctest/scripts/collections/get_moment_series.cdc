import TopShot from 0x03

// This transaction gets the series associated with a moment
// in a collection by geting a reference to the moment
// and then looking up its series

pub fun main(): UInt32 {

    // get the Address of the account with the Moment
    let acct = getAccount(0x01)

    // Get that account's published collectionRef
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    // Get a reference to a specific NFT in the collection
    let ref = collectionRef.borrowNFT(id: 1)

    log(TopShot.setDatas[ref.data.setID]!.series)

    return TopShot.setDatas[ref.data.setID]!.series
}
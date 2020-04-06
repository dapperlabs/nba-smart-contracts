import TopShot from 0x03

// This transaction gets the playID associated with a moment
// in a collection by geting a reference to the moment
// and then looking up its playID 

pub fun main(address: Address, momentID: UInt64): UInt32 {

    // get the Address of the account with the Moment
    let acct = getAccount(address)

    // Get that account's published collectionRef
    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    // Get a reference to a specific NFT in the collection
    let ref = collectionRef.borrowNFT(id: momentID)

    log(ref.data.playID)

    return ref.data.playID
}
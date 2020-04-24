import TopShot from 0x03

// This transaction gets the playID associated with a moment
// in a collection by geting a reference to the moment
// and then looking up its playID 

pub fun main(): UInt32 {

    // get the Address of the account with the Moment
    let acct = getAccount(0x01)

    // Get that account's published collectionRef
    let collectionRef = acct.getCapability(/public/MomentCollection)!
                            .borrow<&{TopShot.MomentCollectionPublic}>()!

    // Get a reference to a specific NFT in the collection
    let ref = collectionRef.borrowNFT(id: 1)

    log(ref.data.playID)

    return ref.data.playID
}
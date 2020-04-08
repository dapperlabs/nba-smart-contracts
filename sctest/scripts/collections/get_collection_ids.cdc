import TopShot from 0x03

// This is the script to get a list of all the moments an account owns
// Just change the argument to `getAccount` to whatever account you want
// and as long as they have a published Collection receiver, you can see
// the moments they own.

pub fun main(): [UInt64] {

    let acct = getAccount(0x01)

    let collectionRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no reference")

    log(collectionRef.getIDs())

    return collectionRef.getIDs()
}
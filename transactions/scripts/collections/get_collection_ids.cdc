import TopShot from 0xTOPSHOTADDRESS

// This is the script to get a list of all the moments an account owns
// Just change the argument to `getAccount` to whatever account you want
// and as long as they have a published Collection receiver, you can see
// the moments they own.

pub fun main(account: Address): [UInt64] {

    let acct = getAccount(account)

    let collectionRef = acct.getCapability(/public/MomentCollection)!
                            .borrow<&{TopShot.MomentCollectionPublic}>()!

    log(collectionRef.getIDs())

    return collectionRef.getIDs()
}
import TopShot from 0xTOPSHOTADDRESS

// This transaction gets the playID associated with a moment
// in a collection by geting a reference to the moment
// and then looking up its playID 


pub fun main(account: Address, id: UInt64): UInt32 {
    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)!
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return data.playID
}
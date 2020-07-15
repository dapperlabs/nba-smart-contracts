import TopShot from 0xTOPSHOTADDRESS

// This transaction gets the setID associated with a moment
// in a collection by geting a reference to the moment
// and then looking up its setID 

pub fun main(account: Address, id: UInt64): UInt32 {

    // borrow a public reference to the owner's moment collection 
    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)!
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    // borrow a reference to the specified moment in the collection
    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return data.setID
}
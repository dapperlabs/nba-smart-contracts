import TopShot from 0xTOPSHOTADDRESS

// This transaction gets the set name associated with a moment
// in a collection by geting a reference to the moment
// and then looking up its name

pub fun main(account: Address, id: UInt64): String {
    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)!
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return TopShot.getSetName(setID: data.setID)!
}
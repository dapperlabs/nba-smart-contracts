import TopShot from 0xTOPSHOTADDRESS

// This script gets the series associated with a moment
// in a collection by getting a reference to the moment
// and then looking up its series

// Parameters:
//
// account: The Flow Address of the account whose moment data needs to be read
// id: The unique ID for the moment whose data needs to be read

// Returns: UInt32
// The series associated with a moment with a specified ID

access(all) fun main(account: Address, id: UInt64): UInt32 {

    let collectionRef = getAccount(account).capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)
        ?? panic("Could not get public moment collection reference")

    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    return TopShot.getSetSeries(setID: data.setID)!
}
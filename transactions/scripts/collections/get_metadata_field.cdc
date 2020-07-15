import TopShot from 0xTOPSHOTADDRESS

// This script gets the metadata associated with a moment
// in a collection by looking up its playID and then searching
// for that play's metadata in the TopShot contract. It returns
// the value for the specified metadata field

pub fun main(account: Address, fieldToSearch: String): String {

    // borrow a public reference to the owner's moment collection 
    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)!
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    // borrow a reference to the specified moment in the collection
    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    // Get the tokens data
    let data = token.data

    // Get the metadata field associated with the specific play
    let field = TopShot.getPlayMetaDataByField(playID: data.playID, field: fieldToSearch) ?? panic("Play doesn't exist")

    log(field)

    return field
}
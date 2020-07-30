import TopShot from 0xTOPSHOTADDRESS

// This script gets the metadata associated with a moment
// in a collection by looking up its playID and then searching
// for that play's metadata in the TopShot contract

// Paramters
// account: The Flow Address of the account whose moment data needs to be read
// id: The unique ID for the moment whose data needs to be read
//
// Returns: {String: String} A dictionary of all the play metadata associated
// with the specified moment

pub fun main(account: Address, id: UInt64): {String: String} {

    // get the public capability for the owner's moment collection
    // and borrow a reference to it
    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)!
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    // Borrow a reference to the specified moment
    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    // Get the moment's metadata to access its play and Set IDs
    let data = token.data

    // Use the moment's play ID 
    // to get all the metadata associated with that play
    let metadata = TopShot.getPlayMetaData(playID: data.playID) ?? panic("Play doesn't exist")

    log(metadata)

    return metadata
}
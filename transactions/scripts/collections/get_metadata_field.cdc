import TopShot from 0xTOPSHOTADDRESS

// This script gets the metadata associated with a moment
// in a collection by looking up its playID and then searching
// for that play's metadata in the TopShot contract

pub fun main(account: Address, fieldToSearch: String): String {

    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)!
        .borrow<&{TopShot.MomentCollectionPublic}>()
        ?? panic("Could not get public moment collection reference")

    let token = collectionRef.borrowMoment(id: id)
        ?? panic("Could not borrow a reference to the specified moment")

    let data = token.data

    let field = TopShot.getPlayMetaDataByField(playID: data.playID, field: fieldToSearch) ?? panic("Play doesn't exist")

    log(field)

    return field
}
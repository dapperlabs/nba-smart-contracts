import TopShot from 0xTOPSHOTADDRESS

pub fun main(account: Address, setIDs: [UInt32], playIDs: [UInt32]): Bool {

    assert(
        setIDs.length == playIDs.length,
        message: "set and play ID arrays have mismatched lengths"
    )

    let collectionRef = getAccount(account).getCapability(/public/MomentCollection)
                .borrow<&{TopShot.MomentCollectionPublic}>()
                ?? panic("Could not get public moment collection reference")

    let momentIDs = collectionRef.getIDs()

    // For each SetID/PlayID combo, loop over each moment in the account
    // to see if they own a moment matching that SetPlay.
    var i = 0
    while i < setIDs.length {
        var hasMatchingMoment = false
        for momentID in momentIDs {
            let token = collectionRef.borrowMoment(id: momentID)
                ?? panic("Could not borrow a reference to the specified moment")

            let momentData = token.data
            if momentData.setID == setIDs[i] && momentData.playID == playIDs[i] {
                hasMatchingMoment = true
                break
            }
        }
        if !hasMatchingMoment {
            return false
        }
        i = i + 1
    }
    
    return true
}
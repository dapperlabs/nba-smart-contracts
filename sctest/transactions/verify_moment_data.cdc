import TopShot from 0x02

pub fun main() {
    if verifyIDs(supply: 0, moldID: 2) { log("PASS") 
    } else { log("FAIL") }

    if verifyCollection(account: 0x02, ids: [UInt64(0), UInt64(1)]) { log("PASS") 
    } else { log("FAIL") }
}

pub fun verifyIDs(supply: UInt64, moldID: UInt32): Bool  {
    log("verifyIDs")

    if TopShot.totalSupply != supply && TopShot.moldID != moldID {
        log("Wrong IDs")
        log("Mold ID")
        log(TopShot.moldID)
        log("Moment ID")
        log(TopShot.totalSupply)
        return false
    } else {
        return true
    }
}

pub fun verifyCollection(account: Address, ids: [UInt64]): Bool {
    log("verifyCollection")

    let acct = getAccount(account)

    if let collectionRef = acct.published[&TopShot.MomentCollectionPublic] {
        var i = 0
        let collectionIDs = collectionRef.getIDs()

        while i < ids.length {
            if ids[i] != collectionIDs[i] {
                log("ID does not exist in the collection!")
                log(ids[i])
                return false
            }

            i = i + 1
        }

        return true
    } else {
        log("No collection reference!")
        return false
    }
}
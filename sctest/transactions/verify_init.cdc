import TopShot from 0x02

// This script is meant to be run after initialization of the TopShot
// contract.  it verifies that everything was initialized correctly.
pub fun main() {
    if verifyIDs(supply: 0, moldID: 0) { log("PASS") 
    } else { log("FAIL") }

    if verifyMoldLen(0) { log("PASS") 
    } else { log("FAIL") }

    if numMomentsLeft(id: 0, quality: 1, expected: 0) { log("PASS") 
    } else { log("FAIL") }

    if numMinted(id: 0, quality: 1, expected: 0) { log("PASS") 
    } else { log("FAIL") }

    if verifyCollectionLength(account: 0x02, 0) { log("PASS") 
    } else { log("FAIL") }

    if verifyCollectionIDs(account: 0x02, ids: []) { log("PASS") 
    } else { log("FAIL") }

    if verifyCreateCollection() { log("PASS") 
    } else { log("FAIL") }

    if verifyAdminNonExistence(account: 0x02) { log("PASS") 
    } else { log("FAIL") }

    if verifyMintingAllowed(id: 0, quality: 1, expected: false) { log("PASS")
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

pub fun verifyMoldLen(_ length: Int): Bool  {
    log("verifyMoldLen")

    if TopShot.molds.length != length {
        log("Incorrect nuber of molds!")
        return false
    } else {
        return true
    }
}

pub fun numMomentsLeft(id: UInt32, quality: Int, expected: UInt32): Bool  {
    log("numMomentsLeft")

    let num = TopShot.getNumMomentsLeftInQuality(id: id, quality: quality)
    if num != expected {
        log("Incorrect number of moments left in specified quality")
        log(num)
        return false
    } else {
        return true
    }
}

pub fun numMinted(id: UInt32, quality: Int, expected: UInt32): Bool {
    log("numMinted")

    let num = TopShot.getNumMintedInQuality(id: id, quality: quality)

    if num != expected {
        log("Incorrect number of moments minted in specified quality")
        log(num)
        return false
    } else {
        return true
    }
}

pub fun verifyCollectionLength(account: Address, _ length: Int): Bool  {
    log("verifyCollectionLength")

    let acct = getAccount(account)

    if let collectionRef = acct.published[&TopShot.MomentCollectionPublic] {
        let collectionIDs = collectionRef.getIDs()
        if  collectionIDs.length != length {
            log("Collection length does not match expected length!")
            log(collectionIDs.length)
            return false
        } else {
            return true
        }
    } else {
        log("No collection!")
        return false
    }
}

pub fun verifyCollectionIDs(account: Address, ids: [UInt64]): Bool {
    log("verifyCollectionIDs")

    let acct = getAccount(account)

    if let collectionRef = acct.published[&TopShot.MomentCollectionPublic] {
        let collectionIDs = collectionRef.getIDs()

        var i = 0

        while i < ids.length {
            if collectionRef.ownedNFTs[ids[i]] == nil {
                log(ids[i])
                return false
            }
            i = i + 1
        }
    } else {
        log("FAIL: No collection!")
        return false
    }
    return true
}

pub fun verifyCreateCollection(): Bool {
    log("verifyCreateCollection")

    let collection <- TopShot.createEmptyCollection()

    if collection.getIDs().length != 0 {
        log("Created collection length should be zero!")
        destroy collection
        return false
    }

    destroy collection
    return true
}

pub fun verifyAdminNonExistence(account: Address): Bool {
    log("verifyAdminNonExistence")

    let acct = getAccount(account)

    if let adminRef = acct.published[&TopShot.Admin] {
        log("Admin should not exist in published!")
        return false
    }
    return true
}

pub fun verifyMintingAllowed(id: UInt32, quality: Int, expected: Bool): Bool {
    log("verifyMintingAllowed")

    if (TopShot.mintingAllowed(id: id, quality: quality) != expected) {
        log("MintingAllowed is incorrect for this ID and quality!")
        return false
    }
    return true
}




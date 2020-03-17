import TopShot from 0x02

transaction {

    // Reference for the collection who will own the minted NFT
    let receiverRef: &TopShot.MomentCollectionPublic

    prepare(acct: Account) {
        // Get the two references from storage
        self.receiverRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no ref!")
    }

    execute {

        if verifyIDs(playID: 0, setID: 0, supply: 0) { log("PASS") 
        } else { log("FAIL") }

        if verifyPlaysLen(0) { log("PASS") 
        } else { log("FAIL") }

        if verifyCollectionLength(account: 0x02, 0) { log("PASS") 
        } else { log("FAIL") }

        if verifyCollectionIDs(account: 0x02, ids: []) { log("PASS") 
        } else { log("FAIL") }

        if verifyCreateCollection() { log("PASS") 
        } else { log("FAIL") }

        if verifySetNonExistence(account: 0x02) { log("PASS") 
        } else { log("FAIL") }
        
        let id1 = TopShot.createPlayData(metadata: {"Name": "Lebron"})

        let id2 = TopShot.createPlayData(metadata: {"Name": "Oladipo"})

        // log("Plays 1 and 2 Succcesfully created!")

        // if verifyIDs(playID: 2, setID: 0, supply: 0) { log("PASS") 
        // } else { log("FAIL") }

        // if verifyPlaysLen(2) { log("PASS") 
        // } else { log("FAIL") }

        // if verifyPlayMetaData(id: 0, key: "Name", value: "Lebron") { log("PASS") 
        // } else { log("FAIL") }

        // Mint two new NFTs from different play IDs
        // let moment1 <- self.adminRef.mintMoment(playID: 0, quality: 1)
        // let moment2 <- self.adminRef.mintMoment(playID: 1, quality: 2)

        // // deposit them into the owner's account
        // self.receiverRef.deposit(token: <-moment1)
        // self.receiverRef.deposit(token: <-moment2)

        // log("Minted Moments successfully!")
        // log("You own these moments!")
        // log(self.receiverRef.getIDs())

        // if verifyIDs(playID: 2, setID: 0, supply: 2) { log("PASS") 
        // } else { log("FAIL") }

        // if numMomentsLeft(id: 0, quality: 1, expected: 2999999999) { log("PASS") 
        // } else { log("FAIL") }

        // if numMinted(id: 0, quality: 1, expected: 1) { log("PASS") 
        // } else { log("FAIL") }

        // if verifyCollection(account: 0x02, ids: [UInt64(0), UInt64(1)]) { log("PASS") 
        // } else { log("FAIL") }
    }
}








// Initialization tests
//
//

pub fun verifyIDs(playID: UInt32, setID: UInt32, supply: UInt64): Bool  {
    log("verifyIDs")

    if TopShot.totalSupply != supply && TopShot.playID != playID && TopShot.setID != setID {
        log("Wrong IDs")
        log("Play ID")
        log(TopShot.playID)
        log("Set ID")
        log(TopShot.setID)
        log("Moment ID")
        log(TopShot.totalSupply)
        return false
    } else {
        return true
    }
}

pub fun verifyPlaysLen(_ length: Int): Bool  {
    log("verifyPlayLen")

    if TopShot.plays.length != length {
        log("Incorrect nuber of plays!")
        return false
    } else {
        return true
    }
}

// pub fun numMinted(id: UInt32, quality: Int, expected: UInt32): Bool {
//     log("numMinted")

//     let num = TopShot.getNumMintedInQuality(id: id, quality: quality)

//     if num != expected {
//         log("Incorrect number of moments minted in specified quality")
//         log(num)
//         return false
//     } else {
//         return true
//     }
// }

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
            if collectionIDs[ids[i]] == nil {
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

pub fun verifySetNonExistence(account: Address): Bool {
    log("verifySetNonExistence")

    let acct = getAccount(account)

    if let adminRef = acct.published[&TopShot.Set] {
        log("Set should not exist in published!")
        return false
    }
    return true
}

// Play Creating Tests
//
//

pub fun verifyPlayIDs(ids: [UInt32]): Bool  {
    log("verifyPlayIDs")

    let i = 0

    while i < ids.length {
        if TopShot.plays[ids[i]] == nil {
            log("play Id doesn't exist")
            log(TopShot.plays[ids[i]])
            return false
        }
    }

    return true
}

pub fun verifyPlayMetaData(id: UInt32, key: String, value: String): Bool {
    log("verifyPlayMetaData")

    if let play = TopShot.plays[id] {
        if play.metadata[key] != value || play.id != id {
            log("Metadata is not what was expected!")
            log(play.metadata[key])
            return false
        } else {
            return true
        }
    } else {
        log("Incorrect play ID")
        return false
    }
}

// Moment Minting Tests
//
//

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
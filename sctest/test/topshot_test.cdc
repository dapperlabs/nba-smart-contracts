import TopShot from 0x02

transaction {

    let adminRef: &TopShot.Admin

    // Reference for the collection who will own the minted NFT
    let receiverRef: &TopShot.MomentCollectionPublic

    prepare(acct: Account) {
        // Get the two references from storage
        self.receiverRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no ref!")
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin
    }

    execute {

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
        
        let id1 = self.adminRef.castMold(metadata: {"Name": "Lebron"}, 
                                         qualityCounts: [UInt32(3000000000), UInt32(1000000000), UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0), UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0), UInt32(100), UInt32(0), 
                                                         UInt32(0), UInt32(10), UInt32(0), UInt32(3)])

        let id2 = self.adminRef.castMold(metadata: {"Name": "Oladipo"}, 
                                         qualityCounts: [UInt32(3000000000), UInt32(1000000000), UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0), UInt32(0), UInt32(0), 
                                                         UInt32(0), UInt32(0), UInt32(100), UInt32(0), 
                                                         UInt32(0), UInt32(10), UInt32(0), UInt32(3)])

        log("Molds 1 and 2 Succcesfully cast!")

        if verifyIDs(supply: 0, moldID: 2) { log("PASS") 
        } else { log("FAIL") }

        if verifyMoldLen(2) { log("PASS") 
        } else { log("FAIL") }

        if numMomentsLeft(id: 0, quality: 1, expected: 3000000000) { log("PASS") 
        } else { log("FAIL") }

        if numMinted(id: 0, quality: 1, expected: 0) { log("PASS") 
        } else { log("FAIL") }

        if verifyMoldMetaData(id: 0, key: "Name", value: "Lebron") { log("PASS") 
        } else { log("FAIL") }

        if verifyMoldQualityCounts(id: 0, counts: [UInt32(3000000000), UInt32(1000000000), UInt32(0), UInt32(0), 
                                                UInt32(0), UInt32(0), UInt32(0), UInt32(0), 
                                                UInt32(0), UInt32(0), UInt32(100), UInt32(0), 
                                                UInt32(0), UInt32(10), UInt32(0), UInt32(3)]) { log("PASS")
        } else { log("FAIL") }

        // Mint two new NFTs from different mold IDs
        let moment1 <- self.adminRef.mintMoment(moldID: 0, quality: 1)
        let moment2 <- self.adminRef.mintMoment(moldID: 1, quality: 2)

        // deposit them into the owner's account
        self.receiverRef.deposit(token: <-moment1)
        self.receiverRef.deposit(token: <-moment2)

        log("Minted Moments successfully!")
        log("You own these moments!")
        log(self.receiverRef.getIDs())

        if verifyIDs(supply: 0, moldID: 2) { log("PASS") 
        } else { log("FAIL") }

        if numMomentsLeft(id: 0, quality: 1, expected: 2999999999) { log("PASS") 
        } else { log("FAIL") }

        if numMinted(id: 0, quality: 1, expected: 1) { log("PASS") 
        } else { log("FAIL") }

        if verifyCollection(account: 0x02, ids: [UInt64(0), UInt64(1)]) { log("PASS") 
        } else { log("FAIL") }
    }
}








// Initialization tests
//
//

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


// Mold Casting Tests
//
//

pub fun verifyMoldIDs(ids: [UInt32]): Bool  {
    log("verifyMoldIDs")

    let i = 0

    while i < ids.length {
        if TopShot.molds[ids[i]] == nil {
            log("mold Id doesn't exist")
            log(TopShot.molds[ids[i]])
            return false
        }
    }

    return true
}

pub fun verifyMoldMetaData(id: UInt32, key: String, value: String): Bool {
    log("verifyMoldMetaData")

    if let mold = TopShot.molds[id] {
        if mold.metadata[key] != value || mold.id != id {
            log("Metadata is not what was expected!")
            log(mold.metadata[key])
            return false
        } else {
            return true
        }
    } else {
        log("Incorrect mold ID")
        return false
    }
}

pub fun verifyMoldQualityCounts(id: UInt32, counts: [UInt32]): Bool {
    log("verifyMoldQualityCounts")

    if let mold = TopShot.molds[id] {
        var i = 0

        while i < counts.length {
            if mold.qualityCounts[i+1] != counts[i] {
                log("Quality count is not what was expected!")
                log("Quality:")
                log(i + 1)
                log("mold Count")
                log(mold.qualityCounts[i+1])
                log("expected count")
                log(counts[i])
                return false
            }
            i = i + 1
        }
        return true

    } else {
        log("Incorrect mold ID")
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
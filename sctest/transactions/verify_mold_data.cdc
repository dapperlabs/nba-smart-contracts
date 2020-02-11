import TopShot from 0x02

// This script is meant to be run whenever you want to verify that certain
// mold-related values are correct in the TopShot contract.
// You'll need to change the function arguments to match what you think the 
// state of your molds are after the minting and casting that you have done.
pub fun main() {
    if verifyIDs(supply: 0, moldID: 2) { log("PASS") 
    } else { log("FAIL") }

    if verifyMoldLen(2) { log("PASS") 
    } else { log("FAIL") }

    if numMomentsLeft(id: 0, quality: 1, expected: 1) { log("PASS") 
    } else { log("FAIL") }

    if numMinted(id: 0, quality: 1, expected: 0) { log("PASS") 
    } else { log("FAIL") }

    if verifyMoldMetaData(id: 0, key: "Name", value: "Lebron") { log("PASS") 
    } else { log("FAIL") }

    if verifyMoldQualityCounts(id: 0, counts: [UInt32(1),UInt32(2),UInt32(3),UInt32(4),UInt32(5)]) { log("PASS")
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

pub fun verifyMoldIDs(ids: [UInt32]): Bool  {
    log("verifyMoldIDs")

    let i = 0

    while i < ids.length {
        if TopShot.molds.ids[i] == nil {
            log("mold Id doesn't exist")
            log(TopShot.molds.ids[i])
            return false
        }
    }

    return true
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

pub fun verifyMoldMetaData(id: UInt32, key: String, value: String): Bool {
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
    if let mold = TopShot.molds[id] {
        var i = 1

        while i < 6 {
            if mold.qualityCounts[i] != counts[i] {
                log("Quality count is not what was expected!")
                log("Quality:")
                log(i)
                log("Count")
                log(mold.qualityCounts[i])
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




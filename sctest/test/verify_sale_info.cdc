import TopShot from 0x02
import Market from 0x03

// This script is meant to be run after initialization of the TopShot
// contract.  it verifies that everything was initialized correctly.
pub fun main() {
    if verifyCut(account: 0x02, expected: 5) { log("PASS") 
    } else { log("FAIL") }

    if verifyPrice(account: 0x02, id: 1, expected: 10) { log("PASS") 
    } else { log("FAIL") }
}

pub fun verifyCut(account: Address, expected: UInt64): Bool {
    log("verifyCut")

    let acct = getAccount(account)

    if let salePublic = acct.published[&Market.SalePublic] {
        if salePublic.cutPercentage != expected {
            log("Incorrect cut percentage!")
            return false
        } else {
            return true
        }
    } else {
        log("No public sale reference!")
        return false
    }
}

pub fun verifyPrice(account: Address, id: UInt64, expected: UInt64): Bool {
    log("verifyPrice")

    let acct = getAccount(account)

    if let salePublic = acct.published[&Market.SalePublic] {
        if salePublic.idPrice(tokenID: id) != expected {
            log("Incorrect price for this ID")
            log(id)
            return false
        } else {
            return true
        }
    } else {
        log("No public sale reference!")
        return false
    }
}




import TopShot from 0x02
import Market from 0x03

// This script is meant to be run after initialization of the TopShot
// contract.  it verifies that everything was initialized correctly.
pub fun main() {

    if verifyNumSales(0) { log("PASS") 
    } else { log("FAIL") }

    if verifySaleLength(0) { log("PASS") 
    } else { log("FAIL") }

    if verifyCut(5) { log("PASS") 
    } else { log("FAIL") }

}

pub fun verifyNumSales(_ expected: Int): Bool  {
    log("verifyNumSales")

    if Market.numSales != expected {
        log("Incorrect nuber of sales!")
        return false
    } else {
        return true
    }
}

pub fun verifySaleLength(_ expected: Int): Bool {
    log("verifySaleLength")

    if Market.saleReferences.length != expected {
        log("Incorrect length of sale dictionary!")
        return false
    } else {
        return true
    }
}

pub fun verifyCut(_ expected: UInt64): Bool {
    log("verifyCut")

    if Market.cutPercentage != expected {
        log("Incorrect cut percentage!")
        return false
    } else {
        return true
    }
}




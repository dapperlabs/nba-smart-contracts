import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

// This script determines how many NFTs are locked in the Top Shot Locking contract

// Returns: Int
// The number of locked NFTs

pub fun main(): Int {
    return TopShotLocking.getLockedNFTsLength()
}

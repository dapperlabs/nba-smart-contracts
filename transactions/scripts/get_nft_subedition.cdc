import TopShotRemix from 0xTOPSHOTREMIXADDRESS

pub fun main(nftID: UInt32): Bool {

    let subedition = TopShotRemix.getMomentsSubedition(momentID: nftID)
        ?? panic("Could not find the specified moment")

    return subedition
}
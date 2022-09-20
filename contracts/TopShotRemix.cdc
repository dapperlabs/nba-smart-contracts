import NonFungibleToken from 0xNFTADDRESS

pub contract TopShotRemix {

    access(self) var numberMintedPerSubedition: {UInt32:{UInt32:{UInt32:UInt32}}}

    pub fun getNumberMintedPerSubedition(setID: UInt32, playID: UInt32, subeditionID: UInt32): UInt32 {
        let numberMintedPerSubedition = self.numberMintedPerSubedition[setID][playID][subeditionID]!
        return <- numberMintedPerSubedition
    }

    pub fun addToNumberMintedPerSubedition(setID: UInt32, playID: UInt32, subeditionID: UInt32) {
        self.numberMintedPerSubedition[setID][playID][subeditionID]!= self.numberMintedPerSubedition[setID][playID][subeditionID] + UInt32(1)
    }

     init() {
            self.numberMintedPerSubedition = {}
     }
}
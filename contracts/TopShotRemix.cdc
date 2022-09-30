import NonFungibleToken from 0xNFTADDRESS

pub contract TopShotRemix {

    access(self) var numberMintedPerSubedition: {UInt32:UInt32}

    access(self) var momentsSubedition: {UInt64:UInt32}

    pub fun getNumberMintedPerSubedition(subeditionID: UInt32): UInt32 {
        let numberMintedPerSubedition = self.numberMintedPerSubedition[subeditionID]!
        return numberMintedPerSubedition
    }

    pub fun getMomentsSubedition( nftID: UInt64):UInt32? {
        return self.momentsSubedition[nftID]
    }


    pub fun addToNumberMintedPerSubedition( subeditionID: UInt32) {
        if TopShotRemix.numberMintedPerSubedition.containsKey(subeditionID) {
            TopShotRemix.numberMintedPerSubedition[subeditionID]!= TopShotRemix.numberMintedPerSubedition[subeditionID]! + UInt32(1)
        } else {
            TopShotRemix.numberMintedPerSubedition[subeditionID]=UInt32(1)
        }
    }

    pub fun setMomentsSubedition(nftID: UInt64, subeditionID: UInt32){
        TopShotRemix.momentsSubedition[nftID] != subeditionID
    }

    init() {
        self.numberMintedPerSubedition = {}
        self.momentsSubedition = {}
    }
}
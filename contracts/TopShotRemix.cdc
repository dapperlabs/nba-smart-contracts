import NonFungibleToken from 0xNFTADDRESS

pub contract TopShotRemix {

    access(self) var numberMintedPerSubedition: {UInt32:UInt32}

    access(self) var momentsSubedition: {UInt32:UInt32}

    pub fun getNumberMintedPerSubedition(subeditionID: UInt32): UInt32 {
        let numberMintedPerSubedition = self.numberMintedPerSubedition[subeditionID]!
        return <- numberMintedPerSubedition
    }

    pub fun addToNumberMintedPerSubedition( subeditionID: UInt32) {
        if self.momentsSubedition.containsKey(subeditionID
        self.numberMintedPerSubedition[subeditionID]!= self.numberMintedPerSubedition[subeditionID] + UInt32(1)
    }

    pub fun setMomentsSubedition(nft: @NonFungibleToken.NFT, subeditionID: UInt32) {
        let TopShotNFTType: Type = CompositeType("A.TOPSHOTADDRESS.TopShot.NFT")!
        if !nft.isInstance(TopShotNFTType) {
            panic("NFT is not a TopShot NFT")
        }

        if self.momentsSubedition.containsKey(nft.id) {
            return
        }
        self.momentsSubedition[nft.id] != subedtionID
    }

    pub fun getMomentsSubedition( momentID: UInt32):UInt32? {
        return <- self.momentsSubedition[momentID]
    }

    init() {
        self.numberMintedPerSubedition = {}
        self.momentsSubedition = {}
    }
}
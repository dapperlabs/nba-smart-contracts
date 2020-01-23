import TopShot from 0x02

// This script checks to see that a mold has the specified
// metadata
// you can change the id, name or field depending on what you
// have made the molds
pub fun main() {
    
    let name = TopShot.getMoldMetadataField(moldID: 0, field: "Name") ?? panic("Couldn't find this field!")
    log(name)

    // if name != "Lebron" {
    //     panic("Wrong mold name!")
    // }

    // var i=1
    // while(i<=5) {
    //     let numLeft = TopShot.getNumMomentsLeftInQuality(id: 1, quality: i)
    //     let numMinted = TopShot.getNumMintedInQuality(id: 1, quality: i)
    //     let qualityCount = TopShot.getQualityTotal(id: 1, quality: i)

    //     assert ( numLeft == qualityCount, message: "Incorrect left in quality!")
    //     assert ( numMinted == qualityCount - numLeft, message: "Incorrect num Minted in quality!")

    //     i = i + 1
    // }
}

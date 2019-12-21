import TopShot from 0x01

transaction {
    prepare(acct: Account) {
        let id1 = acct.storage[TopShot.MoldCaster]?.castMold(metadata: {"Name": "Lebron"}, qualityCounts: {1: 1, 2: 2, 3: 3, 4: 4, 5: 5}) ?? panic("No caster!")
        let id2 = acct.storage[TopShot.MoldCaster]?.castMold(metadata: {"Name": "Oladipo"}, qualityCounts: {1: 0, 2: 1, 3: 10, 4: 20, 5: 40}) ?? panic("No caster!")

        if let mold1 = TopShot.getMoldMetadataField(moldID: id1, field: "Name") {
            log("First name and number of uncommon moments")
            log(mold1)
            log(TopShot.getNumMomentsLeftInQuality(id: id1, quality: 1))
        }

        if let mold2 = TopShot.getMoldMetadataField(moldID: id2, field: "Name") {
            log("Second name and number of uncommon moments")
            log(mold2)
            log(TopShot.getNumMomentsLeftInQuality(id: id2, quality: 1))
        }
    }
}
import TopShotIPFSResolver from "TopShotIPFSResolver"

transaction(
    setIDs: [UInt32],
    playIDs: [UInt32],
    subeditionIDs: [UInt32],
    mediaTypes: [String],
    cids: [String]
) {

    let admin: &TopShotIPFSResolver.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        self.admin = signer.storage.borrow<&TopShotIPFSResolver.Admin>(
            from: TopShotIPFSResolver.AdminStoragePath
        ) ?? panic("Could not borrow Admin resource from storage")
    }

    pre {
        setIDs.length == playIDs.length &&
        playIDs.length == subeditionIDs.length &&
        subeditionIDs.length == mediaTypes.length &&
        mediaTypes.length == cids.length:
            "All input arrays must be the same length"
    }

    execute {
        var i = 0
        while i < setIDs.length {
            self.admin.setCID(
                setID: setIDs[i],
                playID: playIDs[i],
                subeditionID: subeditionIDs[i],
                mediaType: mediaTypes[i],
                cid: cids[i]
            )
            i = i + 1
        }
    }
}

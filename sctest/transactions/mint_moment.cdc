import TopShot from 0x01

transaction {

    prepare(acct: Account) {
        let receiverRef = acct.published[&TopShot.MomentReceiver] ?? panic("no ref!")

        let moment1 = acct.storage[TopShot.MomentMinter]?.mintMoment(moldID: 1, quality: 1, recipient: receiverRef) ?? panic("No minter!")
        let moment2 = acct.storage[TopShot.MomentMinter]?.mintMoment(moldID: 2, quality: 2, recipient: receiverRef) ?? panic("No minter!")

        if (moment1 == 1 && moment2 == 2) {
            log("Minted Moments successfully!")
            log("You own these moments!")
            log(receiverRef.getIDs())
        } else {
            panic("Moment minting failed!")
        }
    }
}
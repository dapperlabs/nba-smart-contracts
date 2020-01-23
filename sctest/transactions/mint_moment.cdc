import TopShot from 0x02

transaction {

    prepare(acct: Account) {
        let receiverRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no ref!")

        let moment1 <- acct.storage[TopShot.MomentMinter]?.mintMoment(moldID: 1, quality: 1) ?? panic("No minter!")
        let moment2 <- acct.storage[TopShot.MomentMinter]?.mintMoment(moldID: 2, quality: 2) ?? panic("No minter!")

        receiverRef.deposit(token: <-moment1)
        receiverRef.deposit(token: <-moment2)

        let ids = receiverRef.getIDs()

        // if (ids[0] == UInt64(0) && ids[1] == UInt64(1)) {
        //     log("Minted Moments successfully!")
        //     log("You own these moments!")
        //     log(ids)
        // } else {
        //     panic("Moment minting failed!")
        // }
    }
}
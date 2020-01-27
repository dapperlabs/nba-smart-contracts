import TopShot from 0x02

transaction {

    prepare(acct: Account) {
        let receiverRef = acct.published[&TopShot.MomentCollectionPublic] ?? panic("no ref!")

        let moment1 <- acct.storage[&TopShot.Admin]?.mintMoment(moldID: 0, quality: 1) ?? panic("No minter!")
        let moment2 <- acct.storage[&TopShot.Admin]?.mintMoment(moldID: 1, quality: 2) ?? panic("No minter!")

        receiverRef.deposit(token: <-moment1)
        receiverRef.deposit(token: <-moment2)

        log("Minted Moments successfully!")
        log("You own these moments!")
        log(receiverRef.getIDs())
    }
}
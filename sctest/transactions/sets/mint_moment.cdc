import TopShot from 0x03

transaction {

    // Reference for the collection who will own the minted NFT
    let receiverRef: &AnyResource{TopShot.MomentCollectionPublic}

    // Reference to the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // Get the two references from storage
        self.receiverRef = acct.published[&AnyResource{TopShot.MomentCollectionPublic}] ?? panic("no ref!")
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin

    }

    execute {

        let setRef = self.adminRef.getSetRef(setID: 1)

        // Mint two new NFTs from different mold IDs
        let moment1 <- setRef.mintMoment(playID: 1)
        let moment2 <- setRef.mintMoment(playID: 1)

        // deposit them into the owner's account
        self.receiverRef.deposit(token: <-moment1)
        self.receiverRef.deposit(token: <-moment2)

        log("Minted Moments successfully!")
        log("You own these moments!")
        log(self.receiverRef.getIDs())
    }
}
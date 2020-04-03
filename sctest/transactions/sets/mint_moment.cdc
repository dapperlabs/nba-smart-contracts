import TopShot from 0x03

// This transaction allows an admin to mint a new moment and
// deposit it into an NFT Collection

transaction {

    // Reference for the collection who will own the minted NFT
    let receiverRef: &AnyResource{TopShot.MomentCollectionPublic}

    // Reference to the admin resource
    let adminRef: &TopShot.Admin

    prepare(acct: AuthAccount) {
        // Get the two references from storage
        self.receiverRef = acct.published[&TopShot.Collection{TopShot.MomentCollectionPublic}] ?? panic("no ref!")
        self.adminRef = &acct.storage[TopShot.Admin] as &TopShot.Admin

    }

    execute {

        let setRef = self.adminRef.borrowSet(setID: 0)

        // Mint a new NFT
        let moment1 <- setRef.mintMoment(playID: 0)

        // deposit them into the owner's account
        self.receiverRef.deposit(token: <-moment1)

        log("Minted Moment successfully!")
        log("You own these moments!")
        log(self.receiverRef.getIDs())
    }
}
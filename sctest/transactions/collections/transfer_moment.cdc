import TopShot from 0x03


transaction {

    // The field that will hold the NFT as it is being
    // transfered to the other account
    let transferToken: @TopShot.NFT
	
    prepare(acct: AuthAccount) {

        // call the withdraw function on the sender's Collection
        // to move the NFT out of the collection
        self.transferToken <- acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!.withdraw(withdrawID: 0)
    }

    execute {
        // get the recipient's public account object
        let recipient = getAccount(0x01)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!

        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-self.transferToken)
    }
}

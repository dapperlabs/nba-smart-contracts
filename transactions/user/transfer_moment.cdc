import TopShot from 0xTOPSHOTADDRESS

transaction {

    let transferToken: @NonFungibleToken.NFT
    
    prepare(acct: AuthAccount) {

        self.transferToken <- acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!.withdraw(withdrawID: %d)
    }

    execute {
        // get the recipient's public account object
        let recipient = getAccount(0x%s)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!

        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-self.transferToken)
    }
}
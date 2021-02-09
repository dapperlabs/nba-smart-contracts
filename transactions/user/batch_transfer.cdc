import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS

transaction(recipientAddress: Address, momentIDs: [UInt64]) {

    let transferTokens: @NonFungibleToken.Collection
    
    prepare(acct: AuthAccount) {

        self.transferTokens <- acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!.batchWithdraw(ids: momentIDs)
    }

    execute {
        // get the recipient's public account object
        let recipient = getAccount(recipientAddress)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Could not borrow a reference to the recipients moment receiver")

        // deposit the NFT in the receivers collection
        receiverRef.batchDeposit(token: <-self.transferTokens)
    }
}
import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

transaction(recipient: Address, momentID: UInt64) {

    let transferToken: @NonFungibleToken.NFT
    
    prepare(acct: AuthAccount) {

        self.transferToken <- acct.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection)!.withdraw(withdrawID: momentID)
    }

    execute {
        // get the recipient's public account object
        let recipient = getAccount(recipient)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()!

        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-self.transferToken)
    }
}
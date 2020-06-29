import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

transaction(recipientAddr: Address, momentIDs: [UInt64]) {
    prepare(acct: AuthAccount) {
        let recipient = getAccount(recipientAddr)
        let receiverRef = recipient.getCapability(/public/MomentCollection)!
            .borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Could not borrow reference to receiver's collection")

        let momentIDs = [momentIDs]

        let collection <- acct.borrow<&TopShotShardedCollection.ShardedCollection>
            (from: /storage/ShardedMomentCollection)!
            .batchWithdraw(ids: momentIDs)
            
        receiverRef.batchDeposit(tokens: <-collection)
    }
}
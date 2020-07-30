import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

// This transaction is what Top Shot uses to send the moments in a "pack" to
// a user's collection

transaction(recipientAddr: Address, momentIDs: [UInt64]) {

    prepare(acct: AuthAccount) {
        
        // get the recipient's public account object
        let recipient = getAccount(recipientAddr)

        // borrow a reference to the recipient's moment collection
        let receiverRef = recipient.getCapability(/public/MomentCollection)!
            .borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Could not borrow reference to receiver's collection")

        // borrow a reference to the owner's moment collection
        let collection <- acct.borrow<&TopShotShardedCollection.ShardedCollection>
            (from: /storage/ShardedMomentCollection)!
            .batchWithdraw(ids: momentIDs)
            
        // Deposit the pack of moments to the recipient's collection
        receiverRef.batchDeposit(tokens: <-collection)
    }
}
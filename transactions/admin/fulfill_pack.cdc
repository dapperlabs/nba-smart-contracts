import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS

// This transaction is what Top Shot uses to send the moments in a "pack" to
// a user's collection

// Parameters:
//
// recipientAddr: the Flow address of the account receiving a pack of moments
// momentsIDs: an array of moment IDs to be withdrawn from the owner's moment collection

transaction(recipientAddr: Address, momentIDs: [UInt64]) {

    prepare(acct: auth(BorrowValue) &Account) {
        
        // get the recipient's public account object
        let recipient = getAccount(recipientAddr)

        // borrow a reference to the recipient's moment collection
        let receiverRef = recipient.capabilities.borrow<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection)
            ?? panic("Cannot borrow a reference to the recipient's collection")

        

        // borrow a reference to the owner's moment collection
        if let collection = acct.storage.borrow<auth(NonFungibleToken.Withdraw) &TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection) {
            
            receiverRef.batchDeposit(tokens: <-collection.batchWithdraw(ids: momentIDs))
        } else {

            let collection = acct.storage.borrow<auth(NonFungibleToken.Withdraw) &TopShot.Collection>(from: /storage/MomentCollection)!

            // Deposit the pack of moments to the recipient's collection
            receiverRef.batchDeposit(tokens: <-collection.batchWithdraw(ids: momentIDs))

        }
    }
}
package templates

import (
	"fmt"
	"strconv"

	"github.com/onflow/flow-go-sdk"
)

// GenerateSetupShardedCollectionScript creates a script that sets up an account to use topshot
func GenerateSetupShardedCollectionScript(nftAddr, tokenCodeAddr flow.Address, numBuckets int) []byte {
	template := `
	import TopShot from 0x%s
	import TopShotShardedCollection from 0x%s
	
	transaction {
	
		prepare(acct: AuthAccount) {
	
			if acct.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection) == nil {
				let collection <- TopShotShardedCollection.createEmptyCollection(numBuckets: %d)
				// Put a new Collection in storage
				acct.save(<-collection, to: /storage/ShardedMomentCollection)
	
				// create a public capability for the collection
				if acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/ShardedMomentCollection) == nil {
					acct.unlink(/public/MomentCollection)

					acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/ShardedMomentCollection)
				}
			} else {
				panic("Sharded Collection already exists!")
			}
		}
	}
	 
	`
	return []byte(fmt.Sprintf(template, nftAddr, tokenCodeAddr.String(), numBuckets))
}

// GenerateTransferMomentfromShardedCollectionScript creates a script that transfers a moment
func GenerateTransferMomentfromShardedCollectionScript(nftAddr, tokenCodeAddr, shardedAddr, recipientAddr flow.Address, tokenID int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s
		import TopShotShardedCollection from 0x%s

		transaction {

			let transferToken: @NonFungibleToken.NFT
			
			prepare(acct: AuthAccount) {
		
				self.transferToken <- acct.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection)!.withdraw(withdrawID: %d)
			}
		
			execute {
				// get the recipient's public account object
				let recipient = getAccount(0x%s)
		
				// get the Collection reference for the receiver
				let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!
		
				// deposit the NFT in the receivers collection
				receiverRef.deposit(token: <-self.transferToken)
			}
		}`
	return []byte(fmt.Sprintf(template, nftAddr, tokenCodeAddr.String(), shardedAddr, tokenID, recipientAddr))
}

// GenerateBatchTransferMomentfromShardedCollectionScript creates a script that transfers a moment
func GenerateBatchTransferMomentfromShardedCollectionScript(nftAddr, tokenCodeAddr, shardedAddr, recipientAddr flow.Address, momentIDs []uint64) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s
		import TopShotShardedCollection from 0x%s

		transaction {

			let transferTokens: @NonFungibleToken.Collection
			
			prepare(acct: AuthAccount) {
				let momentIDs = [%s]
		
				self.transferTokens <- acct.borrow<&TopShotShardedCollection.ShardedCollection>(from: /storage/ShardedMomentCollection)!.batchWithdraw(ids: momentIDs)
			}
		
			execute {
				// get the recipient's public account object
				let recipient = getAccount(0x%s)
		
				// get the Collection reference for the receiver
				let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!
		
				// deposit the NFT in the receivers collection
				receiverRef.batchDeposit(tokens: <-self.transferTokens)
			}
		}`

	// Stringify moment IDs
	momentIDList := ""
	for _, momentID := range momentIDs {
		id := strconv.Itoa(int(momentID))
		momentIDList = momentIDList + `UInt64(` + id + `), `
	}
	// Remove comma and space from last entry
	if idListLen := len(momentIDList); idListLen > 2 {
		momentIDList = momentIDList[:len(momentIDList)-2]
	}
	return []byte(fmt.Sprintf(template, nftAddr, tokenCodeAddr.String(), shardedAddr, momentIDList, recipientAddr))
}

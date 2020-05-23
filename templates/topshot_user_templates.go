package templates

import (
	"fmt"
	"strconv"

	"github.com/onflow/flow-go-sdk"
)

// GenerateSetupAccountScript creates a script that sets up an account to use topshot
func GenerateSetupAccountScript(nftAddr, tokenCodeAddr flow.Address) ([]byte, error) {
	template := `
	import NonFungibleToken from 0x%s
	import TopShot from 0x%s
	
	transaction {
	
		prepare(acct: AuthAccount) {
	
			if acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection) == nil {
				let collection <- TopShot.createEmptyCollection() as! @TopShot.Collection
				// Put a new Collection in storage
				acct.save(<-collection, to: /storage/MomentCollection)
	
				// create a public capability for the collection
				acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/MomentCollection)
			}
		}
	}
	 
	`
	return []byte(fmt.Sprintf(template, nftAddr, tokenCodeAddr.String())), nil
}

// GenerateTransferMomentScript creates a script that transfers a moment
func GenerateTransferMomentScript(nftAddr, tokenCodeAddr, recipientAddr flow.Address, tokenID int) ([]byte, error) {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

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
		}`
	return []byte(fmt.Sprintf(template, nftAddr, tokenCodeAddr.String(), tokenID, recipientAddr)), nil
}

// GenerateBatchTransferMomentScript creates a script that transfers a moment
func GenerateBatchTransferMomentScript(nftAddr, tokenCodeAddr, recipientAddr flow.Address, momentIDs []uint64) ([]byte, error) {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s

		transaction {

			let transferTokens: @NonFungibleToken.Collection
			
			prepare(acct: AuthAccount) {
				let momentIDs = [%s]
		
				self.transferTokens <- acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!.batchWithdraw(ids: momentIDs)
			}
		
			execute {
				// get the recipient's public account object
				let recipient = getAccount(0x%s)
		
				// get the Collection reference for the receiver
				let receiverRef = recipient.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()!
		
				// deposit the NFT in the receivers collection
				receiverRef.batchDeposit(token: <-self.transferTokens)
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
	return []byte(fmt.Sprintf(template, nftAddr, tokenCodeAddr.String(), momentIDList, recipientAddr)), nil
}

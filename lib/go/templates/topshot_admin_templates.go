package templates

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/onflow/flow-go-sdk"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/data"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	transactionsPath        = "../../../transactions/admin/"
	createPlayFilename      = "create_play.cdc"
	createSetFilename       = "create_set.cdc"
	addPlayFilename         = "add_play_to_set.cdc"
	addPlaysFilename        = "add_plays_to_set.cdc"
	lockSetFilename         = "lock_set.cdc"
	retirePlayFilename      = "retire_play.cdc"
	retireAllFilename       = "retire_all.cdc"
	newSeriesFilename       = "start_new_series.cdc"
	mintMomentFilename      = "mint_moment.cdc"
	batchMintMomentFilename = "batch_mint_moment.cdc"
	fulfillPackFilname      = "fulfill_pack.cdc"
)

// GenerateMintPlayScript creates a new play data struct
// and initializes it with metadata
func GenerateMintPlayScript(topShotAddr flow.Address, metadata data.PlayMetadata) []byte {
	md, err := json.Marshal(metadata)
	if err != nil {
		return nil
	}
	code := assets.MustAssetString(transactionsPath + createPlayFilename)

	code = replaceAddresses(code, topShotAddr.String(), "", "", "")

	return []byte(fmt.Sprintf(code, string(md)))
}

// GenerateMintSetScript creates a new Set struct and initializes its metadata
func GenerateMintSetScript(topShotAddr flow.Address, name string) []byte {
	template := `
		import TopShot from 0x%s

		transaction {
			prepare(acct: AuthAccount) {
				let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
				admin.createSet(name: "%s")
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), name))
}

// GenerateAddPlayToSetScript adds a play to a set
// so that moments can be minted from the combo
func GenerateAddPlayToSetScript(topShotAddr flow.Address, setID, playID uint32) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {

			prepare(acct: AuthAccount) {
				let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
				let setRef = admin.borrowSet(setID: %d)
				setRef.addPlay(playID: %d)
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID, playID))
}

// GenerateAddPlaysToSetScript adds multiple plays to a set
func GenerateAddPlaysToSetScript(topShotAddr flow.Address, setID uint32, playIDs []uint32) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {

			prepare(acct: AuthAccount) {
				let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
				let setRef = admin.borrowSet(setID: %d)
				setRef.addPlays(playIDs: %s)
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID, uint32ToCadenceArr(playIDs)))
}

// GenerateMintMomentScript generates a script to mint a new moment
// from a play-set combination
func GenerateMintMomentScript(topShotAddr, recipientAddress flow.Address, setID, playID uint32) []byte {
	template := `
		import TopShot from 0x%s

		transaction {
			let adminRef: &TopShot.Admin
		
			prepare(acct: AuthAccount) {
				self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
			}
		
			execute {
				let setRef = self.adminRef.borrowSet(setID: %d)

				// Mint a new NFT
				let moment1 <- setRef.mintMoment(playID: %d)
				let recipient = getAccount(0x%s)
				// get the Collection reference for the receiver
				let receiverRef = recipient.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()!
				// deposit the NFT in the receivers collection
				receiverRef.deposit(token: <-moment1)
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID, playID, recipientAddress.String()))
}

// GenerateBatchMintMomentScript mints multiple moments of the same play-set combination
func GenerateBatchMintMomentScript(topShotAddr flow.Address, destinationAccount flow.Address, setID, playID uint32, quantity uint64) []byte {
	template := `
		import TopShot from 0x%s

		transaction {
			let adminRef: &TopShot.Admin
		
			prepare(acct: AuthAccount) {
				self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)!
			}
		
			execute {
				let setRef = self.adminRef.borrowSet(setID: %d)

				// Mint a new NFT
				let collection <- setRef.batchMintMoment(playID: %d, quantity: %d)
				let recipient = getAccount(0x%s)
				// get the Collection reference for the receiver
				let receiverRef = recipient.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()!
				// deposit the NFT in the receivers collection
				receiverRef.batchDeposit(tokens: <-collection)
			}
		}`

	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID, playID, quantity, destinationAccount))
}

// GenerateRetirePlayScript retires a play from a set
func GenerateRetirePlayScript(topShotAddr flow.Address, setID, playID int) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {
			let adminRef: &TopShot.Admin

			prepare(acct: AuthAccount) {
				self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
					?? panic("No admin resource in storage")
			}

			execute {
				let setRef = self.adminRef.borrowSet(setID: %d)

				setRef.retirePlay(playID: UInt32(%d))
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID, playID))
}

// GenerateRetireAllPlaysScript retires all plays from a set
func GenerateRetireAllPlaysScript(topShotAddr flow.Address, setID int) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {
			let adminRef: &TopShot.Admin

			prepare(acct: AuthAccount) {
				self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
					?? panic("No admin resource in storage")
			}

			execute {
				let setRef = self.adminRef.borrowSet(setID: %d)

				setRef.retireAll()
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID))
}

// GenerateLockSetScript locks a set
func GenerateLockSetScript(topShotAddr flow.Address, setID int) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {
			let adminRef: &TopShot.Admin

			prepare(acct: AuthAccount) {
				self.adminRef = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
					?? panic("No admin resource in storage")
			}

			execute {
				let setRef = self.adminRef.borrowSet(setID: %d)

				setRef.lock()
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String(), setID))
}

// GenerateFulfillPackScript creates a script that fulfulls a pack
func GenerateFulfillPackScript(topShotAddr, shardedAddr, destinationAccount flow.Address, momentIDs []uint64) []byte {
	template := `
		import TopShot from 0x%s
		import TopShotShardedCollection from 0x%s

		transaction {
			prepare(acct: AuthAccount) {
				let recipient = getAccount(0x%s)
				let receiverRef = recipient.getCapability(/public/MomentCollection)
					.borrow<&{TopShot.MomentCollectionPublic}>()
					?? panic("Could not borrow reference to receiver's collection")

				let momentIDs = [%s]

				let collection <- acct.borrow<&TopShotShardedCollection.ShardedCollection>
					(from: /storage/ShardedMomentCollection)!
					.batchWithdraw(ids: momentIDs)
					
				receiverRef.batchDeposit(tokens: <-collection)
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

	return []byte(fmt.Sprintf(template, topShotAddr.String(), shardedAddr.String(), destinationAccount.String(), momentIDList))
}

// GenerateTransferAdminScript generates a script to create and admin capability
// and transfer it to another account's admin receiver
func GenerateTransferAdminScript(topshotAddr, adminReceiverAddr flow.Address) []byte {
	template := `
		import TopShot from 0x%s
		import TopshotAdminReceiver from 0x%s
		
		transaction {
		
			prepare(acct: AuthAccount) {
				let admin <- acct.load<@TopShot.Admin>(from: /storage/TopShotAdmin)
					?? panic("No topshot admin in storage")

				TopshotAdminReceiver.storeAdmin(newAdmin: <-admin)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr.String(), adminReceiverAddr.String()))
}

// GenerateChangeSeriesScript uses the admin to update the current series
func GenerateChangeSeriesScript(topShotAddr flow.Address) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {
			prepare(acct: AuthAccount) {
				let admin = acct.borrow<&TopShot.Admin>(from: /storage/TopShotAdmin)
					?? panic("No admin resource in storage")
				admin.startNewSeries()
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String()))
}

// GenerateInvalidChangePlaysScript tries to modify the playDatas dictionary
// which should be invalid
func GenerateInvalidChangePlaysScript(topShotAddr flow.Address) []byte {
	template := `
		import TopShot from 0x%s
		
		transaction {
			prepare(acct: AuthAccount) {
				TopShot.playDatas[UInt32(1)] = nil
			}
		}`
	return []byte(fmt.Sprintf(template, topShotAddr.String()))
}

// GenerateUnsafeNotInitializingSetCodeScript generates a script to upgrade the topshot
// contract
func GenerateUnsafeNotInitializingSetCodeScript(newCode []byte) []byte {
	template := `
		
		transaction {
			prepare(acct: AuthAccount, admin: AuthAccount) {
				acct.unsafeNotInitializingSetCode("%s".decodeHex())
			}
		}`
	return []byte(fmt.Sprintf(template, hex.EncodeToString(newCode)))
}

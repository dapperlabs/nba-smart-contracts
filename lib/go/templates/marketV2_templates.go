package templates

import (
	"fmt"

	"github.com/onflow/flow-go-sdk"
)

// This contains template transactions for the second version of the Top Shot
// marketplace, which uses a capability to access the owner's moment collection
// instead of storing the moments in the sale collection directly

// GenerateCreateSaleV2Script creates a cadence transaction that creates a Sale collection
// and stores in in the callers account published
func GenerateCreateSaleV2Script(ftAddr, topshotAddr, marketAddr, beneficiaryAddr flow.Address, tokenStorageName string, cutPercentage float64) []byte {

	template := `
		import FungibleToken from 0x%[5]s
		import TopShot from 0x%[6]s
		import Market from 0x%[1]s

		transaction {
			prepare(acct: AuthAccount) {
				let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!
				let beneficiaryCapability = getAccount(0x%[2]s).getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!

				let ownerCollection: Capability<&TopShot.Collection> = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

				let collection <- Market.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: %[4]f)
				
				acct.save(<-collection, to: /storage/topshotSaleCollection)
				
				acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
			}
		}`
	return []byte(fmt.Sprintf(template, marketAddr, beneficiaryAddr, tokenStorageName, cutPercentage, ftAddr, topshotAddr))
}

// GenerateStartSaleV2Script creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleV2Script(topshotAddr, marketAddr flow.Address, id, price int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.listForSale(tokenID: %[3]d, price: %[4]d.0)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id, price))
}

// GenerateCreateAndStartSaleV2Script creates a cadence transaction that creates a Sale collection
// and stores in in the callers account, and also puts an NFT up for sale in it
func GenerateCreateAndStartSaleV2Script(ftAddr, topshotAddr, marketAddr, beneficiaryAddr flow.Address, tokenStorageName string, cutPercentage, price float64, tokenID int) []byte {

	template := `
		import FungibleToken from 0x%[8]s
		import Market from 0x%[1]s
		import TopShot from 0x%[7]s

		transaction {
			prepare(acct: AuthAccount) {
				// check to see if a sale collection already exists
				if acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection) == nil {
					// get the fungible token capabilities for the owner and beneficiary
					let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!
					let beneficiaryCapability = getAccount(0x%[2]s).getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!

					let ownerCollection = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

					// create a new sale collection
					let topshotSaleCollection <- Market.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: %[4]f)
					
					// save it to storage
					acct.save(<-topshotSaleCollection, to: /storage/topshotSaleCollection)
				
					// create a public link to the sale collection
					acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
				}

				// borrow a reference to the sale
				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				// set the new cut percentage
				topshotSaleCollection.changePercentage(%[4]f)
				
				// put the moment up for sale
				topshotSaleCollection.listForSale(tokenID: %[5]d, price: %[6]f)
				
			}
		}`
	return []byte(fmt.Sprintf(template, marketAddr, beneficiaryAddr, tokenStorageName, cutPercentage, tokenID, price, topshotAddr, ftAddr))
}

// GenerateCancelSaleV2Script creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateCancelSaleV2Script(topshotAddr, marketAddr flow.Address, id int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				// cancel the moment from the sale, thereby de-listing it
				topshotSaleCollection.cancelSale(tokenID: %[3]d)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id))
}

// GenerateChangePriceV2Script creates a cadence transaction that changes the price on an existing sale
func GenerateChangePriceV2Script(topshotAddr, marketAddr flow.Address, id, price int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				// Change the price of the moment
				topshotSaleCollection.listForSale(tokenID: %[3]d, price: %[4]d.0)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id, price))
}

// GenerateChangeOwnerReceiverV2Script creates a cadence transaction
// that changes the sellers receiver capability
func GenerateChangeOwnerReceiverV2Script(fungibleTokenAddr, topshotAddr, marketAddr flow.Address, receiverName string) []byte {
	template := `
		import FungibleToken from 0x%[4]s
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.changeOwnerReceiver(acct.getCapability<&{FungibleToken.Receiver}>(/public/%[3]s)!)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, receiverName, fungibleTokenAddr))
}

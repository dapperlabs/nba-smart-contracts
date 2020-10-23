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
		import TopShotMarketV2 from 0x%[1]s

		transaction {
			prepare(acct: AuthAccount) {
				let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!
				let beneficiaryCapability = getAccount(0x%[2]s).getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!

				let ownerCollection: Capability<&TopShot.Collection> = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

				let collection <- TopShotMarketV2.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: %[4]f)
				
				acct.save(<-collection, to: TopShotMarketV2.marketStoragePath)
				
				acct.link<&TopShotMarketV2.SaleCollection{TopShotMarketV2.SalePublic}>(TopShotMarketV2.marketPublicPath, target: TopShotMarketV2.marketStoragePath)
			}
		}`
	return []byte(fmt.Sprintf(template, marketAddr, beneficiaryAddr, tokenStorageName, cutPercentage, ftAddr, topshotAddr))
}

// GenerateStartSaleV2Script creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleV2Script(topshotAddr, marketAddr flow.Address, id, price int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import TopShotMarketV2 from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
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
		import TopShotMarketV2 from 0x%[1]s
		import TopShot from 0x%[7]s

		transaction {
			prepare(acct: AuthAccount) {
				// check to see if a sale collection already exists
				if acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath) == nil {
					// get the fungible token capabilities for the owner and beneficiary
					let ownerCapability = acct.getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!
					let beneficiaryCapability = getAccount(0x%[2]s).getCapability<&{FungibleToken.Receiver}>(/public/%[3]sReceiver)!

					let ownerCollection = acct.link<&TopShot.Collection>(/private/MomentCollection, target: /storage/MomentCollection)!

					// create a new sale collection
					let topshotSaleCollection <- TopShotMarketV2.createSaleCollection(ownerCollection: ownerCollection, ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: %[4]f)
					
					// save it to storage
					acct.save(<-topshotSaleCollection, to: TopShotMarketV2.marketStoragePath)
				
					// create a public link to the sale collection
					acct.link<&TopShotMarketV2.SaleCollection{TopShotMarketV2.SalePublic}>(TopShotMarketV2.marketPublicPath, target: TopShotMarketV2.marketStoragePath)
				}

				// borrow a reference to the sale
				let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
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
		import TopShotMarketV2 from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
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
		import TopShotMarketV2 from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
					?? panic("Could not borrow from sale in storage")

				// Change the price of the moment
				topshotSaleCollection.listForSale(tokenID: %[3]d, price: %[4]d.0)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id, price))
}

// GenerateChangePercentageV2Script creates a cadence transaction that changes the cut percentage of an existing sale
func GenerateChangePercentageV2Script(topshotAddr, marketAddr flow.Address, percentage float64) []byte {
	template := `
		import TopShot from 0x%[1]s
		import TopShotMarketV2 from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.changePercentage(%[3]f)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, percentage))
}

// GenerateChangeOwnerReceiverV2Script creates a cadence transaction
// that changes the sellers receiver capability
func GenerateChangeOwnerReceiverV2Script(fungibleTokenAddr, topshotAddr, marketAddr flow.Address, receiverName string) []byte {
	template := `
		import FungibleToken from 0x%[4]s
		import TopShot from 0x%[1]s
		import TopShotMarketV2 from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&TopShotMarketV2.SaleCollection>(from: TopShotMarketV2.marketStoragePath)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.changeOwnerReceiver(acct.getCapability<&{FungibleToken.Receiver}>(/public/%[3]s)!)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, receiverName, fungibleTokenAddr))
}

// GenerateBuySaleV2Script creates a cadence transaction that makes a purchase of
// an existing sale
func GenerateBuySaleV2Script(fungibleTokenAddr, tokenAddr, topshotAddr, marketAddr, sellerAddr flow.Address, tokenName, tokenStorageName string, id, amount int) []byte {
	template := `
		import FungibleToken from 0x%[9]s
		import %[1]s from 0x%[2]s
		import TopShot from 0x%[8]s
		import TopShotMarketV2 from 0x%[3]s

		transaction {
			prepare(acct: AuthAccount) {
				let seller = getAccount(0x%[4]s)

				let collection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
					?? panic("Could not borrow reference to the Moment Collection")

				let provider = acct.borrow<&%[1]s.Vault{FungibleToken.Provider}>(from: /storage/%[5]sVault)!

				let tokens <- provider.withdraw(amount: %[6]d.0) as! @%[1]s.Vault

				let topshotSaleCollection = seller.getCapability(TopShotMarketV2.marketPublicPath)!
					.borrow<&{TopShotMarketV2.SalePublic}>()
					?? panic("Could not borrow public sale reference")

				let purchasedToken <- topshotSaleCollection.purchase(tokenID: %[7]d, buyTokens: <-tokens)

				collection.deposit(token: <-purchasedToken)

			}
		}`
	return []byte(fmt.Sprintf(template, tokenName, tokenAddr, marketAddr, sellerAddr, tokenStorageName, amount, id, topshotAddr, fungibleTokenAddr))
}

// GenerateMintTokensAndBuyV2Script creates a script that uses the admin resource
// from the admin accountto mint new tokens and use them to purchase a topshot
// moment from a market collection
func GenerateMintTokensAndBuyV2Script(fungibleTokenAddr, tokenAddr, topshotAddr, marketAddr, sellerAddr, receiverAddr flow.Address, tokenName, storageName string, tokenID, amount int) []byte {
	template := `
		import FungibleToken from 0x%[10]s
		import %[1]s from 0x%[2]s
		import TopShot from 0x%[3]s
		import TopShotMarketV2 from 0x%[4]s

		transaction {

			prepare(signer: AuthAccount) {

			  	let tokenAdmin = signer
					.borrow<&%[1]s.Administrator>(from: /storage/%[5]sAdmin)
					?? panic("Signer is not the token admin")

				let minter <- tokenAdmin.createNewMinter(allowedAmount: UFix64(%[6]d))
				let mintedVault <- minter.mintTokens(amount: UFix64(%[6]d)) as! @%[1]s.Vault

				destroy minter

				let seller = getAccount(0x%[7]s)
				let topshotSaleCollection = seller.getCapability(TopShotMarketV2.marketPublicPath)!
					.borrow<&{TopShotMarketV2.SalePublic}>()
					?? panic("Could not borrow public sale reference")

			  	let boughtToken <- topshotSaleCollection.purchase(tokenID: %[8]d, buyTokens: <-mintedVault)

			  	// get the recipient's public account object and borrow a reference to their moment receiver
			  	let recipient = getAccount(0x%[9]s)
			  		.getCapability(/public/MomentCollection)!.borrow<&{TopShot.MomentCollectionPublic}>()
					?? panic("Could not borrow a reference to the moment collection")

			  	// deposit the NFT in the receivers collection
			  	recipient.deposit(token: <-boughtToken)
			}
		}
	`

	return []byte(fmt.Sprintf(template, tokenName, tokenAddr, topshotAddr, marketAddr, storageName, amount, sellerAddr, tokenID, receiverAddr, fungibleTokenAddr))
}

// GenerateInspectSaleV2Script creates a script that retrieves a sale collection
// from storage and checks that the price is correct
func GenerateInspectSaleV2Script(saleCodeAddr, userAddr flow.Address, nftID int, expectedPrice int) []byte {
	template := `
		import TopShotMarketV2 from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(TopShotMarketV2.marketPublicPath)!.borrow<&{TopShotMarketV2.SalePublic}>()
				?? panic("Could not borrow capability from public collection")

			if collectionRef.getPrice(tokenID: UInt64(%d))! != UFix64(%d) {
				panic("Price for token ID is not correct")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, nftID, expectedPrice))
}

// GenerateInspectSalePercentageV2Script creates a script that retrieves a sale collection
// from storage and checks that the cut percentage is correct
func GenerateInspectSalePercentageV2Script(saleCodeAddr, userAddr flow.Address, percentage float64) []byte {
	template := `
		import TopShotMarketV2 from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(TopShotMarketV2.marketPublicPath)!.borrow<&{TopShotMarketV2.SalePublic}>()
				?? panic("Could not borrow capability from public collection")

			if collectionRef.cutPercentage != UFix64(%f) {
				panic("Cut percentage is incorrect")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, percentage))
}

// GenerateInspectSaleLenV2Script creates a script that retrieves an NFT collection
// from storage and tries to borrow a reference for an NFT that it owns.
// If it owns it, it will not fail.
func GenerateInspectSaleLenV2Script(saleCodeAddr, userAddr flow.Address, length int) []byte {
	template := `
		import TopShotMarketV2 from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(TopShotMarketV2.marketPublicPath)!
				.borrow<&{TopShotMarketV2.SalePublic}>()
				?? panic("Could not borrow capability from public collection")

			if %d != collectionRef.getIDs().length {
				panic("Collection Length is not correct")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, length))
}

// GenerateInspectSaleMomentDataV2Script creates a script that checks
// a sale for a certain ID and makes sure it has the right set
func GenerateInspectSaleMomentDataV2Script(nftAddr, tokenAddr, marketAddr, ownerAddr flow.Address, expectedID, expectedSet int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s
		import TopShotMarketV2 from 0x%s

		pub fun main() {
			let saleRef = getAccount(0x%s).getCapability(TopShotMarketV2.marketPublicPath)!
				.borrow<&{TopShotMarketV2.SalePublic}>()
				?? panic("Could not get public sale reference")

			let token = saleRef.borrowMoment(id: %d)
				?? panic("Could not borrow a reference to the specified moment")

			let data = token.data

			assert(
                data.setID == UInt32(%d),
                message: "ID %d does not have the expected Set ID %d"
            )
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, tokenAddr, marketAddr, ownerAddr, expectedID, expectedSet, expectedID, expectedSet))
}

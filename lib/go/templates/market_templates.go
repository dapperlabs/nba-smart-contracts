package templates

import (
	"fmt"

	"github.com/onflow/flow-go-sdk"
)

// GenerateCreateSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account published
func GenerateCreateSaleScript(marketAddr, beneficiaryAddr flow.Address, tokenStorageName string, cutPercentage float64) []byte {

	template := `
		import Market from 0x%[1]s

		transaction {
			prepare(acct: AuthAccount) {
				let ownerCapability = acct.getCapability(/public/%[3]sReceiver)
				let beneficiaryCapability = getAccount(0x%[2]s).getCapability(/public/%[3]sReceiver)

				let collection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: %[4]f)
				
				acct.save(<-collection, to: /storage/topshotSaleCollection)
				
				acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
			}
		}`
	return []byte(fmt.Sprintf(template, marketAddr, beneficiaryAddr, tokenStorageName, cutPercentage))
}

// GenerateStartSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleScript(topshotAddr, marketAddr flow.Address, id, price int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {
				let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
					?? panic("Could not borrow from MomentCollection in storage")

                let token <- nftCollection.withdraw(withdrawID: %[3]d) as! @TopShot.NFT

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.listForSale(token: <-token, price: %[4]d.0)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id, price))
}

// GenerateCreateAndStartSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account, and also puts an NFT up for sale in it
func GenerateCreateAndStartSaleScript(topshotAddr, marketAddr, beneficiaryAddr flow.Address, tokenStorageName string, cutPercentage float64, tokenID, price int) []byte {

	template := `
		import Market from 0x%[1]s
		import TopShot from 0x%[7]s

		transaction {
			prepare(acct: AuthAccount) {
				// check to see if a sale collection already exists
				if acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection) == nil {
					// get the fungible token capabilities for the owner and beneficiary
					let ownerCapability = acct.getCapability(/public/%[3]sReceiver)
					let beneficiaryCapability = getAccount(0x%[2]s).getCapability(/public/%[3]sReceiver)

					// create a new sale collection
					let topshotSaleCollection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: %[4]f)
					
					// save it to storage
					acct.save(<-topshotSaleCollection, to: /storage/topshotSaleCollection)
				
					// create a public link to the sale collection
					acct.link<&Market.SaleCollection{Market.SalePublic}>(/public/topshotSaleCollection, target: /storage/topshotSaleCollection)
				}
				
				// borrow a reference to the seller's moment collection
				let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
					?? panic("Could not borrow from MomentCollection in storage")

				// withdraw the moment to put up for sale
				let token <- nftCollection.withdraw(withdrawID: %[5]d) as! @TopShot.NFT

				// borrow a reference to the sale
				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				// set the new cut percentage
				topshotSaleCollection.changePercentage(%[4]f)
				
				// the the moment for sale
				topshotSaleCollection.listForSale(token: <-token, price: UFix64(%[6]d))
				
			}
		}`
	return []byte(fmt.Sprintf(template, marketAddr, beneficiaryAddr, tokenStorageName, cutPercentage, tokenID, price, topshotAddr))
}

// GenerateWithdrawFromSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateWithdrawFromSaleScript(topshotAddr, marketAddr flow.Address, id int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {
				let nftCollection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
					?? panic("Could not borrow from MomentCollection in storage")

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				let token <- topshotSaleCollection.withdraw(tokenID: %[3]d)

				nftCollection.deposit(token: <-token)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id))
}

// GenerateChangePriceScript creates a cadence transaction that changes the price on an existing sale
func GenerateChangePriceScript(topshotAddr, marketAddr flow.Address, id, price int) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.changePrice(tokenID: %[3]d, newPrice: %[4]d.0)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, id, price))
}

// GenerateChangePercentageScript creates a cadence transaction that changes the cut percentage of an existing sale
func GenerateChangePercentageScript(topshotAddr, marketAddr flow.Address, percentage float64) []byte {
	template := `
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.changePercentage(%[3]f)
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, percentage))
}

// GenerateChangeOwnerReceiverScript creates a cadence transaction
// that changes the sellers receiver capability
func GenerateChangeOwnerReceiverScript(fungibleTokenAddr, topshotAddr, marketAddr flow.Address, receiverName string) []byte {
	template := `
		import FungibleToken from 0x%[4]s
		import TopShot from 0x%[1]s
		import Market from 0x%[2]s

		transaction {
			prepare(acct: AuthAccount) {

				let topshotSaleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/topshotSaleCollection)
					?? panic("Could not borrow from sale in storage")

				topshotSaleCollection.changeOwnerReceiver(acct.getCapability(/public/%[3]s))
			}
		}`
	return []byte(fmt.Sprintf(template, topshotAddr, marketAddr, receiverName, fungibleTokenAddr))
}

// GenerateBuySaleScript creates a cadence transaction that makes a purchase of
// an existing sale
func GenerateBuySaleScript(fungibleTokenAddr, tokenAddr, topshotAddr, marketAddr, sellerAddr flow.Address, tokenName, tokenStorageName string, id, amount int) []byte {
	template := `
		import FungibleToken from 0x%[9]s
		import %[1]s from 0x%[2]s
		import TopShot from 0x%[8]s
		import Market from 0x%[3]s

		transaction {
			prepare(acct: AuthAccount) {
				let seller = getAccount(0x%[4]s)

				let collection = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
					?? panic("Could not borrow reference to the Moment Collection")

				let provider = acct.borrow<&%[1]s.Vault{FungibleToken.Provider}>(from: /storage/%[5]sVault)!
				
				let tokens <- provider.withdraw(amount: %[6]d.0) as! @%[1]s.Vault

				let topshotSaleCollection = seller.getCapability(/public/topshotSaleCollection)
					.borrow<&{Market.SalePublic}>()
					?? panic("Could not borrow public sale reference")
			
				let purchasedToken <- topshotSaleCollection.purchase(tokenID: %[7]d, buyTokens: <-tokens)

				collection.deposit(token: <-purchasedToken)

			}
		}`
	return []byte(fmt.Sprintf(template, tokenName, tokenAddr, marketAddr, sellerAddr, tokenStorageName, amount, id, topshotAddr, fungibleTokenAddr))
}

// GenerateMintTokensAndBuyScript creates a script that uses the admin resource
// from the admin accountto mint new tokens and use them to purchase a topshot
// moment from a market collection
func GenerateMintTokensAndBuyScript(fungibleTokenAddr, tokenAddr, topshotAddr, marketAddr, sellerAddr, receiverAddr flow.Address, tokenName, storageName string, tokenID, amount int) []byte {
	template := `
		import FungibleToken from 0x%[10]s
		import %[1]s from 0x%[2]s
		import TopShot from 0x%[3]s
		import Market from 0x%[4]s
	
		transaction {
	
			prepare(signer: AuthAccount) {

			  	let tokenAdmin = signer
					.borrow<&%[1]s.Administrator>(from: /storage/%[5]sAdmin) 
					?? panic("Signer is not the token admin")

				let minter <- tokenAdmin.createNewMinter(allowedAmount: UFix64(%[6]d))
				let mintedVault <- minter.mintTokens(amount: UFix64(%[6]d)) as! @%[1]s.Vault

				destroy minter
	
				let seller = getAccount(0x%[7]s)
				let topshotSaleCollection = seller.getCapability(/public/topshotSaleCollection)
					.borrow<&{Market.SalePublic}>()
					?? panic("Could not borrow public sale reference")
	
			  	let boughtToken <- topshotSaleCollection.purchase(tokenID: %[8]d, buyTokens: <-mintedVault)

			  	// get the recipient's public account object and borrow a reference to their moment receiver
			  	let recipient = getAccount(0x%[9]s)
			  		.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
					?? panic("Could not borrow a reference to the moment collection")
	
			  	// deposit the NFT in the receivers collection
			  	recipient.deposit(token: <-boughtToken)
			}
		}
	`

	return []byte(fmt.Sprintf(template, tokenName, tokenAddr, topshotAddr, marketAddr, storageName, amount, sellerAddr, tokenID, receiverAddr, fungibleTokenAddr))
}

// GenerateInspectSaleScript creates a script that retrieves a sale collection
// from storage and checks that the price is correct
func GenerateInspectSaleScript(saleCodeAddr, userAddr flow.Address, nftID int, expectedPrice int) []byte {
	template := `
		import Market from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(/public/topshotSaleCollection).borrow<&{Market.SalePublic}>()
				?? panic("Could not borrow capability from public collection")
			
			if collectionRef.getPrice(tokenID: UInt64(%d))! != UFix64(%d) {
				panic("Price for token ID is not correct")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, nftID, expectedPrice))
}

// GenerateInspectSalePercentageScript creates a script that retrieves a sale collection
// from storage and checks that the cut percentage is correct
func GenerateInspectSalePercentageScript(saleCodeAddr, userAddr flow.Address, percentage float64) []byte {
	template := `
		import Market from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(/public/topshotSaleCollection).borrow<&{Market.SalePublic}>()
				?? panic("Could not borrow capability from public collection")
			
			if collectionRef.cutPercentage != UFix64(%f) {
				panic("Cut percentage is incorrect")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, percentage))
}

// GenerateInspectSaleLenScript creates a script that retrieves an NFT collection
// from storage and tries to borrow a reference for an NFT that it owns.
// If it owns it, it will not fail.
func GenerateInspectSaleLenScript(saleCodeAddr, userAddr flow.Address, length int) []byte {
	template := `
		import Market from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(/public/topshotSaleCollection)
				.borrow<&{Market.SalePublic}>()
				?? panic("Could not borrow capability from public collection")
			
			if %d != collectionRef.getIDs().length {
				panic("Collection Length is not correct")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, length))
}

// GenerateInspectSaleMomentDataScript creates a script that checks
// a sale for a certain ID and makes sure it has the right set
func GenerateInspectSaleMomentDataScript(nftAddr, tokenAddr, marketAddr, ownerAddr flow.Address, expectedID, expectedSet int) []byte {
	template := `
		import NonFungibleToken from 0x%s
		import TopShot from 0x%s
		import Market from 0x%s

		pub fun main() {
			let saleRef = getAccount(0x%s).getCapability(/public/topshotSaleCollection)
				.borrow<&{Market.SalePublic}>()
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

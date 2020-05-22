package templates

import (
	"fmt"

	"github.com/onflow/flow-go-sdk"
)

// GenerateCreateSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account published
func GenerateCreateSaleScript(ftAddr, topshotAddr, marketAddr flow.Address, cutPercentage float64) []byte {

	template := `
		import FungibleToken from 0x%s
		import TopShot from 0x%s
		import Market from 0x%s

		transaction {
			prepare(acct: AuthAccount) {
				let ownerCapability = acct.getCapability(/public/flowTokenReceiver)!
				let topshotCapability = getAccount(0x%s).getCapability(/public/flowTokenReceiver)!

				let collection <- Market.createSaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: topshotCapability, cutPercentage: %f)
				
				acct.save(<-collection, to: /storage/saleCollection)
				
				acct.link<&{Market.SalePublic}>(/public/saleCollection, target: /storage/saleCollection)
			}
		}`
	return []byte(fmt.Sprintf(template, ftAddr, topshotAddr, marketAddr, topshotAddr, cutPercentage))
}

// GenerateStartSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleScript(nftAddr flow.Address, marketAddr flow.Address, id, price int) []byte {
	template := `
		import NonFungibleToken, ExampleNFT from 0x%s
		import Market from 0x%s

		transaction {
			prepare(acct: AuthAccount) {
				let nftCollection = acct.borrow<&ExampleNFT.Collection>(from: /storage/NFTCollection)
					?? panic("Could not borrow from NFTCollection in storage")

                let token <- nftCollection.withdraw(withdrawID: %d)

				let saleCollection = acct.borrow<&Market.SaleCollection>(from: /storage/saleCollection)
					?? panic("Could not borrow from sale in storage")

				saleCollection.listForSale(token: <-token, price: %d.0)
			}
		}`
	return []byte(fmt.Sprintf(template, nftAddr, marketAddr, id, price))
}

// GenerateBuySaleScript creates a cadence transaction that makes a purchase of
// an existing sale
func GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, userAddr flow.Address, id, amount int) []byte {
	template := `
		import FungibleToken, FlowToken from 0x%s
		import NonFungibleToken, ExampleNFT from 0x%s
		import Market from 0x%s

		transaction {
			prepare(acct: AuthAccount) {
				let seller = getAccount(0x%s)

				let collection = acct.borrow<&ExampleNFT.Collection>(from: /storage/NFTCollection)
					?? panic("Could not borrow public reference to NFT Collection")

				let provider = acct.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!
				
				let tokens <- provider.withdraw(amount: %d.0)

				let saleCollection = seller.getCapability(/public/saleCollection)!
					.borrow<&{Market.SalePublic}>()
					?? panic("Could not borrow public sale reference")
			
				saleCollection.purchase(tokenID: %d, recipient: collection, buyTokens: <-tokens)

			}
		}`
	return []byte(fmt.Sprintf(template, tokenAddr, nftAddr, marketAddr, userAddr, amount, id))
}

// GenerateInspectSaleScript creates a script that retrieves a sale collection
// from storage and checks that the price is correct
func GenerateInspectSaleScript(saleCodeAddr, userAddr flow.Address, nftID int, price int) []byte {
	template := `
		import Market from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(/public/saleCollection)!.borrow<&{Market.SalePublic}>()
				?? panic("Could not borrow capability from public collection")
			
			if collectionRef.prices[UInt64(%d)] != UFix64(%d) {
				panic("Price for token ID is not correct")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, nftID, price))
}

// GenerateInspectSaleLenScript creates a script that retrieves an NFT collection
// from storage and tries to borrow a reference for an NFT that it owns.
// If it owns it, it will not fail.
func GenerateInspectSaleLenScript(saleCodeAddr, userAddr flow.Address, length int) []byte {
	template := `
		import Market from 0x%s

		pub fun main() {
			let acct = getAccount(0x%s)
			let collectionRef = acct.getCapability(/public/saleCollection)!
				.borrow<&{Market.SalePublic}>()
				?? panic("Could not borrow capability from public collection")
			
			if %d != collectionRef.getIDs().length {
				panic("Collection Length is not correct")
			}
		}
	`

	return []byte(fmt.Sprintf(template, saleCodeAddr, userAddr, length))
}

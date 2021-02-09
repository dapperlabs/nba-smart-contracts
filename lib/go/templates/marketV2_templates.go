package templates

import (
	"fmt"
	"strings"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
	"github.com/onflow/flow-go-sdk"
)

const (
	createSaleV2Filename          = "marketV2/create_sale.cdc"
	startSaleV2Filename           = "marketV2/start_sale.cdc"
	createAndStartSaleV2Filename  = "marketV2/create_start_sale.cdc"
	withdrawSaleV2Filename        = "marketV2/stop_sale.cdc"
	changePriceV2Filename         = "marketV2/change_price.cdc"
	changePercentageV2Filename    = "marketV2/change_percentage.cdc"
	changeOwnerReceiverV2Filename = "marketV2/change_receiver.cdc"
	purchaseV2Filename            = "marketV2/purchase_moment.cdc"
	mintAndPurchaseV2Filename     = "marketV2/mint_and_purchase.cdc"
	upgradeSaleFilename           = "marketV2/upgrade_sale.cdc"

	// scripts
	getSalePriceV2Filename      = "marketV2/scripts/get_sale_price.cdc"
	getSalePercentageV2Filename = "marketV2/scripts/get_sale_percentage.cdc"
	getSaleLengthV2Filename     = "marketV2/scripts/get_sale_len.cdc"
	getSaleSetIDV2Filename      = "marketV2/scripts/get_sale_set_id.cdc"
)

// This contains template transactions for the second version of the Top Shot
// marketplace, which uses a capability to access the owner's moment collection
// instead of storing the moments in the sale collection directly

// GenerateCreateSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account published
func GenerateCreateSaleV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createSaleV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateStartSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + startSaleV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateCreateAndStartSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account, and also puts an NFT up for sale in it
func GenerateCreateAndStartSaleV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createAndStartSaleV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateCancelSaleV2Script creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateCancelSaleV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + withdrawSaleV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangePriceScript creates a cadence transaction that changes the price on an existing sale
func GenerateChangePriceV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changePriceV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangePercentageScript creates a cadence transaction that changes the cut percentage of an existing sale
func GenerateChangePercentageV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changePercentageV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangeOwnerReceiverScript creates a cadence transaction
// that changes the sellers receiver capability
func GenerateChangeOwnerReceiverV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changeOwnerReceiverV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBuySaleScript creates a cadence transaction that makes a purchase of
// an existing sale
func GenerateBuySaleV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + purchaseV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateMintTokensAndBuyScript creates a script that uses the admin resource
// from the admin accountto mint new tokens and use them to purchase a topshot
// moment from a market collection
func GenerateMintTokensAndBuyV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + mintAndPurchaseV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSalePriceScript creates a script that retrieves a sale collection
// and returns the price of the specified moment
func GenerateGetSalePriceV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSalePriceV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSalePercentageScript creates a script that retrieves a sale collection
// from storage and returns the cut percentage
func GenerateGetSalePercentageV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSalePercentageV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSaleLenScript creates a script that retrieves an NFT collection
// reference and returns its length
func GenerateGetSaleLenV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSaleLengthV2Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSaleSetIDScript creates a script that checks
// a sale for a certain ID and returns its set ID
func GenerateGetSaleSetIDV2Script(env Environment) []byte {

	code := assets.MustAssetString(transactionsPath + getSaleSetIDV2Filename)

	return []byte(replaceAddresses(code, env))
}

func GenerateUpgradeSaleV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + upgradeSaleFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateMultiContractP2PPurchaseScript(ftInterfaceAddr, topshotAddr, marketV2Addr, marketAddr, sellerAccount, p2pTokenAddr flow.Address, momentID uint64, tokenName, storageName string) []byte {
	template := `
		import FungibleToken from 0x{{.FTInterfaceAddress}}
		import {{.P2PTokenName}} from 0x{{.P2PTokenAddress}}
		import TopShot from 0x{{.TopShotAddress}}
		import Market from 0x{{.P2PMarketAddress}}
		import MarketV2 from 0x{{.P2PMarketV2Address}}

		// This transaction purchases a moment by first checking if it is in the first version of the market collecion
		// If it isn't in the first version, it checks if it is in the second and purchases it there

		transaction(seller: Address, recipient: Address, momentID: UInt64, purchaseAmount: UFix64) {

			let purchaseTokens: @{{.P2PTokenName}}.Vault

			prepare(acct: AuthAccount) {

				// Borrow a provider reference to the buyers vault
				let provider = acct.borrow<&{{.P2PTokenName}}.Vault{FungibleToken.Provider}>(from: /storage/{{.P2PTokenName}}Vault)
					?? panic("Could not borrow a reference to the buyers {{.P2PTokenName}} Vault")

				// withdraw the purchase tokens from the vault
				self.purchaseTokens <- provider.withdraw(amount: purchaseAmount) as! @{{.P2PTokenName}}.Vault
			}

			execute {

				// get the accounts for the seller and recipient
				let seller = getAccount(0x{{.SellerAccountAddress}})
				let recipient = getAccount(0x{{.Recipient}})

				// Get the reference for the recipient's nft receiver
				let receiverRef = recipient.getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
					?? panic("Could not borrow a reference to the recipients moment collection")

				// Check if the V1 market collection exists
				if let marketCollection = seller.getCapability(/public/topshotSaleCollection)
					.borrow<&{Market.SalePublic}>() {

					// Check if the V1 market has the moment for sale
					if marketCollection.borrowMoment(id: momentID) != nil {
						// purchase from the V1 market
						let purchasedToken <- marketCollection.purchase(tokenID: momentID, buyTokens: <-tokens)
						receiverRef.deposit(token: <-purchasedToken)
					}

				// Check if the V2 market collection exists
				} else if let marketV2Collection = seller.getCapability(/public/topshotSalev2Collection)
						.borrow<&{MarketV2.SalePublic}>() {

					// Check if the V2 market has the moment for sale
					if marketV2Collection.borrowMoment(id: momentID) != nil {

						// Purchase from the V2 market
						let purchasedToken <- marketV2Collection.purchase(tokenID: momentID, buyTokens: <-tokens)
						receiverRef.deposit(token: <-purchasedToken)

					} else {
						panic("Could not find the moment sale in either collection")
					}

				} else {
					panic("Could not find either sale collection")
				}
			}
		}`
	oldNew := []string{
		"{{.FTInterfaceAddress}}", ftInterfaceAddr.String(),
		"{{.P2PTokenName}}", tokenName,
		"{{.P2PTokenAddress}}", p2pTokenAddr.String(),
		"{{.TopShotAddress}}", topshotAddr.String(),
		"{{.P2PMarketAddress}}", marketAddr.String(),
		"{{.P2PMarketV2Address}}", marketV2Addr.String(),
		"{{.SellerAccountAddress}}", sellerAccount.String(),
		"{{.P2PTokenStorageName}}", storageName,
		"{{.MomentFlowID}}", fmt.Sprintf("%d", momentID),
	}
	replacer := strings.NewReplacer(oldNew...)
	return []byte(replacer.Replace(template))
}

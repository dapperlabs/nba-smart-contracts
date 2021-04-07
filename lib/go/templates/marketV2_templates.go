package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
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

	purchaseBothMarketsFilename = "marketV2/purchase_both_markets.cdc"

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

func GenerateMultiContractP2PPurchaseScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + purchaseBothMarketsFilename)

	return []byte(replaceAddresses(code, env))
}

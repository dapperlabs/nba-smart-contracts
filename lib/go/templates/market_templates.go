package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	createSaleFilename          = "market/create_sale.cdc"
	startSaleFilename           = "market/start_sale.cdc"
	createAndStartSaleFilename  = "market/create_start_sale.cdc"
	withdrawSaleFilename        = "market/stop_sale.cdc"
	changePriceFilename         = "market/change_price.cdc"
	changePercentageFilename    = "market/change_percentage.cdc"
	changeOwnerReceiverFilename = "market/change_receiver.cdc"
	purchaseFilename            = "market/purchase_moment.cdc"
	mintAndPurchaseFilename     = "market/mint_and_purchase.cdc"

	// scripts
	getSalePriceFilename      = "market/scripts/get_sale_price.cdc"
	getSalePercentageFilename = "market/scripts/get_sale_percentage.cdc"
	getSaleLengthFilename     = "market/scripts/get_sale_len.cdc"
	getSaleSetIDFilename      = "market/scripts/get_sale_set_id.cdc"
)

// These templates are for the first version of the Top Shot marketplace
// which actually stored moments that were for sale in the sale collections
// in the seller's accounts

// GenerateCreateSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account published
func GenerateCreateSaleScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createSaleFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateStartSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + startSaleFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateCreateAndStartSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account, and also puts an NFT up for sale in it
func GenerateCreateAndStartSaleScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createAndStartSaleFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateWithdrawFromSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateWithdrawFromSaleScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + withdrawSaleFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangePriceScript creates a cadence transaction that changes the price on an existing sale
func GenerateChangePriceScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changePriceFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangePercentageScript creates a cadence transaction that changes the cut percentage of an existing sale
func GenerateChangePercentageScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changePercentageFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangeOwnerReceiverScript creates a cadence transaction
// that changes the sellers receiver capability
func GenerateChangeOwnerReceiverScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changeOwnerReceiverFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBuySaleScript creates a cadence transaction that makes a purchase of
// an existing sale
func GenerateBuySaleScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + purchaseFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateMintTokensAndBuyScript creates a script that uses the admin resource
// from the admin accountto mint new tokens and use them to purchase a topshot
// moment from a market collection
func GenerateMintTokensAndBuyScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + mintAndPurchaseFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSalePriceScript creates a script that retrieves a sale collection
// and returns the price of the specified moment
func GenerateGetSalePriceScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSalePriceFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSalePercentageScript creates a script that retrieves a sale collection
// from storage and returns the cut percentage
func GenerateGetSalePercentageScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSalePercentageFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSaleLenScript creates a script that retrieves an NFT collection
// reference and returns its length
func GenerateGetSaleLenScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSaleLengthFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSaleSetIDScript creates a script that checks
// a sale for a certain ID and returns its set ID
func GenerateGetSaleSetIDScript(env Environment) []byte {

	code := assets.MustAssetString(transactionsPath + getSaleSetIDFilename)

	return []byte(replaceAddresses(code, env))
}

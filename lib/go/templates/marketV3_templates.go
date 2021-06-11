package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	createSaleV3Filename          = "marketV3/create_sale.cdc"
	startSaleV3Filename           = "marketV3/start_sale.cdc"
	createAndStartSaleV3Filename  = "marketV3/create_start_sale.cdc"
	withdrawSaleV3Filename        = "marketV3/stop_sale.cdc"
	changePriceV3Filename         = "marketV3/change_price.cdc"
	changePercentageV3Filename    = "marketV3/change_percentage.cdc"
	changeOwnerReceiverV3Filename = "marketV3/change_receiver.cdc"
	purchaseV3Filename            = "marketV3/purchase_moment.cdc"
	mintAndPurchaseV3Filename     = "marketV3/mint_and_purchase.cdc"
	upgradeSaleFilename           = "marketV3/upgrade_sale.cdc"

	purchaseBothMarketsFilename = "marketV3/purchase_both_markets.cdc"

	// scripts
	getSalePriceV3Filename      = "marketV3/scripts/get_sale_price.cdc"
	getSalePercentageV3Filename = "marketV3/scripts/get_sale_percentage.cdc"
	getSaleLengthV3Filename     = "marketV3/scripts/get_sale_len.cdc"
	getSaleSetIDV3Filename      = "marketV3/scripts/get_sale_set_id.cdc"
)

// This contains template transactions for the third version of the Top Shot
// marketplace, which uses a capability to access the owner's moment collection
// instead of storing the moments in the sale collection directly

// GenerateCreateSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account published
func GenerateCreateSaleV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createSaleV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateStartSaleScript creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateStartSaleV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + startSaleV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateCreateAndStartSaleScript creates a cadence transaction that creates a Sale collection
// and stores in in the callers account, and also puts an NFT up for sale in it
func GenerateCreateAndStartSaleV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createAndStartSaleV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateCancelSaleV3Script creates a cadence transaction that starts a sale by depositing
// an NFT into the Sale Collection with an associated price
func GenerateCancelSaleV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + withdrawSaleV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangePriceScript creates a cadence transaction that changes the price on an existing sale
func GenerateChangePriceV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changePriceV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangePercentageScript creates a cadence transaction that changes the cut percentage of an existing sale
func GenerateChangePercentageV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changePercentageV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangeOwnerReceiverScript creates a cadence transaction
// that changes the sellers receiver capability
func GenerateChangeOwnerReceiverV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + changeOwnerReceiverV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBuySaleScript creates a cadence transaction that makes a purchase of
// an existing sale
func GenerateBuySaleV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + purchaseV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateMintTokensAndBuyScript creates a script that uses the admin resource
// from the admin accountto mint new tokens and use them to purchase a topshot
// moment from a market collection
func GenerateMintTokensAndBuyV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + mintAndPurchaseV3Filename)

	return []byte(replaceAddresses(code, env))
}

func GenerateUpgradeSaleV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + upgradeSaleFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateMultiContractP2PPurchaseScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + purchaseBothMarketsFilename)

	return []byte(replaceAddresses(code, env))
}

/*************** V3 SCRIPTS **************************/

// GenerateGetSalePriceScript creates a script that retrieves a sale collection
// and returns the price of the specified moment
func GenerateGetSalePriceV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSalePriceV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSalePercentageScript creates a script that retrieves a sale collection
// from storage and returns the cut percentage
func GenerateGetSalePercentageV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSalePercentageV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSaleLenScript creates a script that retrieves an NFT collection
// reference and returns its length
func GenerateGetSaleLenV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + getSaleLengthV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSaleSetIDScript creates a script that checks
// a sale for a certain ID and returns its set ID
func GenerateGetSaleSetIDV3Script(env Environment) []byte {

	code := assets.MustAssetString(transactionsPath + getSaleSetIDV3Filename)

	return []byte(replaceAddresses(code, env))
}

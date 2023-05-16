package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	setupAccountFilename   = "user/setup_account.cdc"
	transferMomentFilename = "user/transfer_moment.cdc"
	batchTransferFilename  = "user/batch_transfer.cdc"

	transferMomentV3Filename = "user/transfer_moment_v3_sale.cdc"
	destroyMomentsFilename   = "user/destroy_moments.cdc"
	destroyMomentsV2Filename = "user/destroy_moments_v2.cdc"
)

// GenerateSetupAccountScript creates a script that sets up an account to use topshot
func GenerateSetupAccountScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + setupAccountFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateTransferMomentScript creates a script that transfers a moment
func GenerateTransferMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + transferMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBatchTransferMomentScript creates a script that transfers multiple moments
func GenerateBatchTransferMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + batchTransferFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateTransferMomentScript creates a script that transfers a moment
// and cancels its sale if it is for sale
func GenerateTransferMomentV3Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + transferMomentV3Filename)

	return []byte(replaceAddresses(code, env))
}

// GenerateDestroyMomentsScript creates a script that destroyes select
// moments from a user's collection
func GenerateDestroyMomentsScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + destroyMomentsFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateDestroyMomentsV2Script creates a script that destroys select
// moments from a user's collection using the Top Shot contract destroyMoments function
func GenerateDestroyMomentsV2Script(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + destroyMomentsV2Filename)

	return []byte(replaceAddresses(code, env))
}

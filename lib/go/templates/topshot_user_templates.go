package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	setupAccountFilename   = "user/setup_account.cdc"
	transferMomentFilename = "user/transfer_moment.cdc"
	batchTransferFilename  = "user/batch_transfer.cdc"
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

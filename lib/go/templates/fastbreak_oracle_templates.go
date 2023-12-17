package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	createFastBreakRunFilename = "fastbreak/oracle/create_run.cdc"
)

func GenerateCreateRunScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createFastBreakRunFilename)

	return []byte(replaceAddresses(code, env))
}

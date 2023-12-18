package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	fastBreakScriptsPath = "../../../transactions/fastbreak/scripts/"

	getFastBreakByIdFilename       = "get_fast_break.cdc"
	getFastBreakTokenCountFilename = "get_token_count.cdc"
)

func GenerateGetFastBreakScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getFastBreakByIdFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetFastBreakTokenCountScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getFastBreakTokenCountFilename)

	return []byte(replaceAddresses(code, env))
}

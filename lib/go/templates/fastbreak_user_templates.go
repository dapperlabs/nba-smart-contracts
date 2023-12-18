package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	fastBreakSetupGameFilename = "fastbreak/user/setup_game.cdc"
	fastBreakPlayFilename      = "fastbreak/user/play.cdc"
)

func GenerateFastBreakCreateAccountScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + fastBreakSetupGameFilename)

	return []byte(replaceAddresses(code, env))
}

func GeneratePlayFastBreakScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + fastBreakPlayFilename)

	return []byte(replaceAddresses(code, env))
}

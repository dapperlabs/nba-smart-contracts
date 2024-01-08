package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	fastBreakSetupGameFilename = "fastbreak/player/create_player.cdc"
	fastBreakPlayFilename      = "fastbreak/player/play.cdc"
)

func GenerateFastBreakCreateAccountScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + fastBreakSetupGameFilename)

	return []byte(replaceAddresses(code, env))
}

func GeneratePlayFastBreakScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + fastBreakPlayFilename)

	return []byte(replaceAddresses(code, env))
}

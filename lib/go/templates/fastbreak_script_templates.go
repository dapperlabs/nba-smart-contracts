package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	fastBreakScriptsPath = "../../../transactions/fastbreak/scripts/"

	getFastBreakByIdFilename        = "get_fast_break.cdc"
	getFastBreakTokenCountFilename  = "get_token_count.cdc"
	getScoreByPlayerFilename        = "get_player_score.cdc"
	getFastBreakStatsFilename       = "get_fast_break_stats.cdc"
	fastBreakCurrentPlayer          = "get_current_player.cdc"
	getPlayerWinCountForRunFilename = "get_player_win_count_for_run.cdc"
)

func GenerateGetFastBreakScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getFastBreakByIdFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetFastBreakTokenCountScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getFastBreakTokenCountFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetPlayerScoreScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getScoreByPlayerFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetFastBreakStatsScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getFastBreakStatsFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateCurrentPlayerScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + fastBreakCurrentPlayer)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetPlayerWinCountForRunScript(env Environment) []byte {
	code := assets.MustAssetString(fastBreakScriptsPath + getPlayerWinCountForRunFilename)

	return []byte(replaceAddresses(code, env))
}

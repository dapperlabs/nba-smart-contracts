package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	createFastBreakRunFilename       = "fastbreak/oracle/create_run.cdc"
	createFastBreakGameFilename      = "fastbreak/oracle/create_game.cdc"
	addStatToFastBreakGameFilename   = "fastbreak/oracle/add_stat_to_game.cdc"
	updateFastBreakGameFilename      = "fastbreak/oracle/update_fast_break_game.cdc"
	scoreFastBreakSubmissionFilename = "fastbreak/oracle/score_fast_break_submission.cdc"
)

func GenerateCreateRunScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createFastBreakRunFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateCreateGameScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createFastBreakGameFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateAddStatToGameScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + addStatToFastBreakGameFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateUpdateFastBreakGameScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + updateFastBreakGameFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateScoreFastBreakSubmissionScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + scoreFastBreakSubmissionFilename)

	return []byte(replaceAddresses(code, env))
}

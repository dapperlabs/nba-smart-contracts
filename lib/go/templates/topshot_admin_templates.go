package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	transactionsPath         = "../../../transactions/"
	createPlayFilename       = "admin/create_play.cdc"
	createSetFilename        = "admin/create_set.cdc"
	addPlayFilename          = "admin/add_play_to_set.cdc"
	addPlaysFilename         = "admin/add_plays_to_set.cdc"
	lockSetFilename          = "admin/lock_set.cdc"
	retirePlayFilename       = "admin/retire_play_from_set.cdc"
	retireAllFilename        = "admin/retire_all.cdc"
	newSeriesFilename        = "admin/start_new_series.cdc"
	mintMomentFilename       = "admin/mint_moment.cdc"
	batchMintMomentFilename  = "admin/batch_mint_moment.cdc"
	fulfillPackFilename      = "admin/fulfill_pack.cdc"
	createSetAndPlayFilename = "admin/create_set_and_play_struct.cdc"

	transferAdminFilename = "admin/transfer_admin.cdc"

	mintMomentWithSubeditionFilename      = "admin/mint_moment_with_subedition.cdc"
	batchMintMomentWithSubeditionFilename = "admin/batch_mint_moment_with_subedition.cdc"
	createNewSubeditionResourceFilename   = "admin/create_new_showcase_resource.cdc"
)

// GenerateMintPlayScript creates a new play data struct
// and initializes it with metadata
func GenerateMintPlayScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createPlayFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateMintSetScript creates a new Set struct and initializes its metadata
func GenerateMintSetScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createSetFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateAddPlayToSetScript adds a play to a set
// so that moments can be minted from the combo
func GenerateAddPlayToSetScript(env Environment) []byte {

	code := assets.MustAssetString(transactionsPath + addPlayFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateAddPlaysToSetScript adds multiple plays to a set
func GenerateAddPlaysToSetScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + addPlaysFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateMintMomentScript generates a script to mint a new moment
// from a play-set combination
func GenerateMintMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + mintMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBatchMintMomentScript mints multiple moments of the same play-set combination
func GenerateBatchMintMomentScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + batchMintMomentFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateRetirePlayScript retires a play from a set
func GenerateRetirePlayScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + retirePlayFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateRetireAllPlaysScript retires all plays from a set
func GenerateRetireAllPlaysScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + retireAllFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateLockSetScript locks a set
func GenerateLockSetScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + lockSetFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateFulfillPackScript creates a script that fulfulls a pack
func GenerateFulfillPackScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + fulfillPackFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateTransferAdminScript generates a script to create and admin capability
// and transfer it to another account's admin receiver
func GenerateTransferAdminScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + transferAdminFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateChangeSeriesScript uses the admin to update the current series
func GenerateChangeSeriesScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + newSeriesFilename)

	return []byte(replaceAddresses(code, env))
}

// For testing purposes only
func GenerateCreateSetandPlayDataScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createSetAndPlayFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateMintMomentWithSubeditionScript generates a script to mint a new moment
// with Subedition from a play-set-subedition combination
func GenerateMintMomentWithSubeditionScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + mintMomentWithSubeditionFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBatchMintMomentWithSubeditionScript mints multiple moments with Subedition
// of the same play-set-subedition combination
func GenerateBatchMintMomentWithSubeditionScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + batchMintMomentWithSubeditionFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateCreateNewSubeditionResourceScript creates new Subedition resource
// for minting with Subeditions
func GenerateCreateNewSubeditionResourceScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + createNewSubeditionResourceFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateInvalidChangePlaysScript tries to modify the playDatas dictionary
// which should be invalid
func GenerateInvalidChangePlaysScript(env Environment) []byte {

	code := `
		import TopShot from 0xTOPSHOTADDRESS
		
		transaction {
			prepare(acct: AuthAccount) {
				TopShot.playDatas[UInt32(1)] = nil
			}
		}`
	return []byte(replaceAddresses(code, env))
}

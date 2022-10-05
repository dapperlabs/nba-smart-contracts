package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	scriptsPath = "../../../transactions/scripts/"

	// Topshot contract scripts
	currentSeriesFilename = "get_currentSeries.cdc"
	totalSupplyFilename   = "get_totalSupply.cdc"

	// Play related scripts
	getAllPlaysFilename = "plays/get_all_plays.cdc"
	nextPlayIDFilename  = "plays/get_nextPlayID.cdc"
	playMetadata        = "plays/get_play_metadata.cdc"
	playMetadataField   = "plays/get_play_metadata_field.cdc"

	// Set related scripts
	editionRetiredFilename      = "sets/get_edition_retired.cdc"
	numMomentsInEditionFilename = "sets/get_numMoments_in_edition.cdc"
	setIDsByNameFilename        = "sets/get_setIDs_by_name.cdc"
	setSeriesFilename           = "sets/get_setSeries.cdc"
	nextSetIDFilename           = "sets/get_nextSetID.cdc"
	playsInSetFilename          = "sets/get_plays_in_set.cdc"
	setNameFilename             = "sets/get_setName.cdc"
	setLockedFilename           = "sets/get_set_locked.cdc"
	getSetMetadataFilename      = "sets/get_set_data.cdc"

	// collections scripts
	collectionIDsFilename       = "collections/get_collection_ids.cdc"
	metadataFieldFilename       = "collections/get_metadata_field.cdc"
	momentSeriesFilename        = "collections/get_moment_series.cdc"
	idInCollectionFilename      = "collections/get_id_in_Collection.cdc"
	momentPlayIDFilename        = "collections/get_moment_playID.cdc"
	momentSetIDFilename         = "collections/get_moment_setID.cdc"
	metadataFilename            = "collections/get_metadata.cdc"
	momentSerialNumFilename     = "collections/get_moment_serialNum.cdc"
	momentSetNameFilename       = "collections/get_moment_setName.cdc"
	getSetPlaysAreOwnedFilename = "collections/get_setplays_are_owned.cdc"

	// metadata scripts
	getNFTMetadataFilename     = "get_nft_metadata.cdc"
	getTopShotMetadataFilename = "get_topshot_metadata.cdc"

	//subEdition scripts
	getNFTSubeditionFilename = "get_nft_subedition.cdc"
)

// Global Data Gettetrs

func GenerateGetSeriesScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + currentSeriesFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetSupplyScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + totalSupplyFilename)

	return []byte(replaceAddresses(code, env))
}

// Play Related Scripts

func GenerateGetAllPlaysScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getAllPlaysFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetNextPlayIDScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + nextPlayIDFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetPlayMetadataScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + playMetadata)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetPlayMetadataFieldScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + playMetadataField)

	return []byte(replaceAddresses(code, env))
}

// Set-related scripts

// GenerateGetIsEditionRetiredScript creates a script that indicates if an edition is retired
func GenerateGetIsEditionRetiredScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + editionRetiredFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetNumMomentsInEditionScript creates a script
// that returns the number of moments that have been minted in an edition
func GenerateGetNumMomentsInEditionScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + numMomentsInEditionFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSetIDsByNameScript creates a script that returns setIDs that share a name
func GenerateGetSetIDsByNameScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + setIDsByNameFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSetNameScript creates a script that returns the name of a set
func GenerateGetSetNameScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + setNameFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetSetSeriesScript creates a script that returns the metadata of a play
func GenerateGetSetSeriesScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + setSeriesFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetNextSetIDScript creates a script that returns next set ID that will be used
func GenerateGetNextSetIDScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + nextSetIDFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetPlaysInSetScript creates a script that returns an array of plays in a set
func GenerateGetPlaysInSetScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + playsInSetFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetIsSetLockedScript creates a script that indicates if a set is locked
func GenerateGetIsSetLockedScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + setLockedFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetSetMetadataScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getSetMetadataFilename)

	return []byte(replaceAddresses(code, env))
}

// Collection related scripts

func GenerateGetCollectionIDsScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + collectionIDsFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentMetadataScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + metadataFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentMetadataFieldScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + metadataFieldFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentSeriesScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + momentSeriesFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateIsIDInCollectionScript creates a script that checks
// a collection for a certain ID
func GenerateIsIDInCollectionScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + idInCollectionFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentPlayScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + momentPlayIDFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentSetScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + momentSetIDFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentSetNameScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + momentSetNameFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetMomentSerialNumScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + momentSerialNumFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateSetPlaysOwnedByAddressScript generates a script that returns true if each of the SetPlays corresponding to
// the passed Set and Play IDs are owned by the passed flow.Address.
//
// Set and Play IDs are matched up by index in the passed slices.
func GenerateSetPlaysOwnedByAddressScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getSetPlaysAreOwnedFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetNFTMetadataScript creates a script that returns the metadata for an NFT.
func GenerateGetNFTMetadataScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getNFTMetadataFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetTopShotMetadataScript creates a script that returns the metadata for an NFT.
func GenerateGetTopShotMetadataScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getTopShotMetadataFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateGetTopShotMetadataScript creates a script that returns the subEdition for an NFT.
func GenerateGetNFTSubEditionScript(env Environment) []byte {
	code := assets.MustAssetString(scriptsPath + getNFTSubeditionFilename)

	return []byte(replaceAddresses(code, env))
}

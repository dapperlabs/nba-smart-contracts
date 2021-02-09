package templates

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/internal/assets"
)

const (
	setupShardedCollectionFilename   = "shardedCollection/setup_sharded_collection.cdc"
	transferFromShardedFilename      = "shardedCollection/transfer_from_sharded.cdc"
	batchTransferFromShardedFilename = "shardedCollection/batch_from_sharded.cdc"
)

// GenerateSetupShardedCollectionScript creates a script that sets up an account to use topshot
func GenerateSetupShardedCollectionScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + setupShardedCollectionFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateTransferMomentfromShardedCollectionScript creates a script that transfers a moment
func GenerateTransferMomentfromShardedCollectionScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + transferFromShardedFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateBatchTransferMomentfromShardedCollectionScript creates a script that transfers a moment
func GenerateBatchTransferMomentfromShardedCollectionScript(env Environment) []byte {
	code := assets.MustAssetString(transactionsPath + batchTransferFromShardedFilename)

	return []byte(replaceAddresses(code, env))
}

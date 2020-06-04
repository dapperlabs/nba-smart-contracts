package contracts

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -prefix ../src/contracts -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../src/contracts

import (
	"strings"

	"github.com/dapperlabs/nba-smart-contracts/contracts/internal/assets"
)

const (
	topshotFile                    = "TopShot.cdc"
	marketFile                     = "MarketTopShot.cdc"
	shardedCollectionFile          = "TopShotShardedCollection.cdc"
	adminReceiverFile              = "TopshotAdminReceiver.cdc"
	defaultNonFungibleTokenAddress = "02"
	defaultFungibleTokenAddress    = "04"
)

// GenerateTopShotContract returns a copy
// of the topshot contract with the import addresses updated
func GenerateTopShotContract(nftAddr string) []byte {

	topShotCode := assets.MustAssetString(topshotFile)

	codeWithNFTAddr := strings.ReplaceAll(topShotCode, defaultNonFungibleTokenAddress, nftAddr)

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotShardedCollectionContract returns a copy
// of the TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr string) []byte {

	shardedCode := assets.MustAssetString(shardedCollectionFile)
	codeWithNFTAddr := strings.ReplaceAll(shardedCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, "03", topshotAddr)

	return []byte(codeWithTopshotAddr)
}

// GenerateTopshotAdminReceiverContract returns a copy
// of the TopshotAdminReceiver contract with the import addresses updated
func GenerateTopshotAdminReceiverContract(topshotAddr, shardedAddr string) []byte {

	adminReceiverCode := assets.MustAssetString(adminReceiverFile)
	codeWithTopshotAddr := strings.ReplaceAll(adminReceiverCode, "03", topshotAddr)
	codeWithShardedAddr := strings.ReplaceAll(codeWithTopshotAddr, "04", shardedAddr)

	return []byte(codeWithShardedAddr)
}

// GenerateTopShotMarketContract returns a copy
// of the TopShotMarketContract with the import addresses updated
func GenerateTopShotMarketContract(ftAddr, nftAddr, topshotAddr, flowTokenAddr string) []byte {

	marketCode := assets.MustAssetString(marketFile)
	codeWithNFTAddr := strings.ReplaceAll(marketCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, "03", topshotAddr)
	codeWithFTAddr := strings.ReplaceAll(codeWithTopshotAddr, defaultFungibleTokenAddress, ftAddr)
	codeWithFlowTokenAddr := strings.ReplaceAll(codeWithFTAddr, "05", flowTokenAddr)

	return []byte(codeWithFlowTokenAddr)
}

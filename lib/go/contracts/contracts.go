package contracts

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -prefix ../../../contracts -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../contracts

import (
	"strings"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts/internal/assets"
)

const (
	topshotFile                    = "TopShot.cdc"
	topshotV1File                  = "TopShotv1.cdc"
	marketFile                     = "MarketTopShot.cdc"
	shardedCollectionFile          = "TopShotShardedCollection.cdc"
	shardedCollectionV1File        = "TopShotShardedCollectionV1.cdc"
	adminReceiverFile              = "TopshotAdminReceiver.cdc"
	defaultNonFungibleTokenAddress = "NFTADDRESS"
	defaultFungibleTokenAddress    = "FUNGIBLETOKENADDRESS"
	defaultTopshotAddress          = "TOPSHOTADDRESS"
	defaultShardedAddress          = "SHARDEDADDRESS"
)

// GenerateTopShotContract returns a copy
// of the topshot contract with the import addresses updated
func GenerateTopShotContract(nftAddr string) []byte {

	topShotCode := assets.MustAssetString(topshotFile)

	codeWithNFTAddr := strings.ReplaceAll(topShotCode, defaultNonFungibleTokenAddress, nftAddr)

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotV1Contract returns a copy
// of the original topshot contract with the import addresses updated
func GenerateTopShotV1Contract(nftAddr string) []byte {

	topShotCode := assets.MustAssetString(topshotV1File)
	codeWithNFTAddr := strings.ReplaceAll(topShotCode, defaultNonFungibleTokenAddress, nftAddr)

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotShardedCollectionContract returns a copy
// of the TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr string) []byte {

	shardedCode := assets.MustAssetString(shardedCollectionFile)
	codeWithNFTAddr := strings.ReplaceAll(shardedCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)

	return []byte(codeWithTopshotAddr)
}

// GenerateTopShotShardedCollectionV1Contract returns a copy
// of the original TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionV1Contract(nftAddr, topshotAddr string) []byte {

	shardedCode := assets.MustAssetString(shardedCollectionV1File)
	codeWithNFTAddr := strings.ReplaceAll(shardedCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)

	return []byte(codeWithTopshotAddr)
}

// GenerateTopshotAdminReceiverContract returns a copy
// of the TopshotAdminReceiver contract with the import addresses updated
func GenerateTopshotAdminReceiverContract(topshotAddr, shardedAddr string) []byte {

	adminReceiverCode := assets.MustAssetString(adminReceiverFile)
	codeWithTopshotAddr := strings.ReplaceAll(adminReceiverCode, defaultTopshotAddress, topshotAddr)
	codeWithShardedAddr := strings.ReplaceAll(codeWithTopshotAddr, defaultShardedAddress, shardedAddr)

	return []byte(codeWithShardedAddr)
}

// GenerateTopShotMarketContract returns a copy
// of the TopShotMarketContract with the import addresses updated
func GenerateTopShotMarketContract(ftAddr, nftAddr, topshotAddr, ducTokenAddr string) []byte {

	marketCode := assets.MustAssetString(marketFile)
	codeWithNFTAddr := strings.ReplaceAll(marketCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)
	codeWithFTAddr := strings.ReplaceAll(codeWithTopshotAddr, defaultFungibleTokenAddress, ftAddr)
	codeWithTokenAddr := strings.ReplaceAll(codeWithFTAddr, "DUCADDRESS", ducTokenAddr)

	return []byte(codeWithTokenAddr)
}

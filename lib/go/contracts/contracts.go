package contracts

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -prefix ../../../contracts -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../contracts

import (
	"strings"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts/internal/assets"
	_ "github.com/kevinburke/go-bindata"
)

const (
	topshotFile  = "TopShot.cdc"
	marketV3File = "TopShotMarketV3.cdc"
	// There is a MarketTopShot.cdc contract which was updated to be token agnostic, however this was not backwards compatible.
	// MarketTopShotOldVersion.cdc is the current contract in production
	marketFile                     = "MarketTopShotOldVersion.cdc"
	shardedCollectionFile          = "TopShotShardedCollection.cdc"
	adminReceiverFile              = "TopshotAdminReceiver.cdc"
	defaultNonFungibleTokenAddress = "NFTADDRESS"
	defaultFungibleTokenAddress    = "FUNGIBLETOKENADDRESS"
	defaultTopshotAddress          = "TOPSHOTADDRESS"
	defaultShardedAddress          = "SHARDEDADDRESS"
	defaultMarketAddress           = "MARKETADDRESS"
	defaultMetadataviewsAddress    = "METADATAVIEWSADDRESS"
	defaultTopShotLockingAddress   = "TOPSHOTLOCKINGADDRESS"
)

// GenerateTopShotContract returns a copy
// of the topshot contract with the import addresses updated
func GenerateTopShotContract(nftAddr string, metadataViewsAddr string, topShotLockingAddr string) []byte {

	topShotCode := assets.MustAssetString(topshotFile)

	codeWithNFTAddr := strings.ReplaceAll(topShotCode, defaultNonFungibleTokenAddress, nftAddr)

	codeWithMetadataViewsAddr := strings.ReplaceAll(codeWithNFTAddr, defaultMetadataviewsAddress, metadataViewsAddr)

	codeWithTopShotLockingAddr := strings.ReplaceAll(codeWithMetadataViewsAddr, defaultTopShotLockingAddress, topShotLockingAddr)

	//log.Printf("code w/ locking %+v", codeWithTopShotLockingAddr)

	return []byte(codeWithTopShotLockingAddr)
}

// GenerateTopShotShardedCollectionContract returns a copy
// of the TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr string) []byte {

	shardedCode := assets.MustAssetString(shardedCollectionFile)
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

// GenerateTopShotMarketV3Contract returns a copy
// of the third version TopShotMarketContract with the import addresses updated
func GenerateTopShotMarketV3Contract(ftAddr, nftAddr, topshotAddr, marketAddr, ducTokenAddr string) []byte {

	marketCode := assets.MustAssetString(marketV3File)
	codeWithNFTAddr := strings.ReplaceAll(marketCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)
	codeWithFTAddr := strings.ReplaceAll(codeWithTopshotAddr, defaultFungibleTokenAddress, ftAddr)
	codeWithMarketV3Addr := strings.ReplaceAll(codeWithFTAddr, defaultMarketAddress, marketAddr)
	codeWithTokenAddr := strings.ReplaceAll(codeWithMarketV3Addr, "DUCADDRESS", ducTokenAddr)

	return []byte(codeWithTokenAddr)
}

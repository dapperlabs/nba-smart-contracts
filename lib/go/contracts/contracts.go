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
	marketFile                     = "Market.cdc"
	shardedCollectionFile          = "TopShotShardedCollection.cdc"
	adminReceiverFile              = "TopshotAdminReceiver.cdc"
	topShotLockingFile             = "TopShotLocking.cdc"
	defaultNonFungibleTokenAddress = "NFTADDRESS"
	defaultFungibleTokenAddress    = "FUNGIBLETOKENADDRESS"
	defaultTopshotAddress          = "TOPSHOTADDRESS"
	defaultShardedAddress          = "SHARDEDADDRESS"
	defaultMarketAddress           = "MARKETADDRESS"
	defaultMetadataviewsAddress    = "METADATAVIEWSADDRESS"
	defaultTopShotLockingAddress   = "TOPSHOTLOCKINGADDRESS"
	defaultTopShotRoyaltyAddress   = "TOPSHOTROYALTYADDRESS"
	defaultViewResolverAddress     = "VIEWRESOLVERADDRESS"
	defaultNetwork                 = "${NETWORK}"
	fastBreakFile                  = "FastBreakV1.cdc"
)

// GenerateTopShotContract returns a copy
// of the topshot contract with the import addresses updated
func GenerateTopShotContract(ftAddr string, nftAddr string, metadataViewsAddr string, viewResolverAddr string, topShotLockingAddr string, royaltyAddr string, network string) []byte {

	topShotCode := assets.MustAssetString(topshotFile)

	codeWithFTAddr := strings.ReplaceAll(topShotCode, defaultFungibleTokenAddress, ftAddr)

	codeWithNFTAddr := strings.ReplaceAll(codeWithFTAddr, defaultNonFungibleTokenAddress, nftAddr)

	codeWithMetadataViewsAddr := strings.ReplaceAll(codeWithNFTAddr, defaultViewResolverAddress, viewResolverAddr)

	codeWithViewResolverAddr := strings.ReplaceAll(codeWithMetadataViewsAddr, defaultMetadataviewsAddress, metadataViewsAddr)

	codeWithTopShotLockingAddr := strings.ReplaceAll(codeWithViewResolverAddr, defaultTopShotLockingAddress, topShotLockingAddr)

	codeWithTopShotRoyaltyAddr := strings.ReplaceAll(codeWithTopShotLockingAddr, defaultTopShotRoyaltyAddress, royaltyAddr)

	codeWithNetwork := strings.ReplaceAll(codeWithTopShotRoyaltyAddr, defaultNetwork, network)

	return []byte(codeWithNetwork)
}

// GenerateTopShotShardedCollectionContract returns a copy
// of the TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr string, viewResolverAddr string) []byte {

	shardedCode := assets.MustAssetString(shardedCollectionFile)
	codeWithNFTAddr := strings.ReplaceAll(shardedCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)
	codeWithViewResolverAddr := strings.ReplaceAll(codeWithTopshotAddr, defaultViewResolverAddress, viewResolverAddr)

	return []byte(codeWithViewResolverAddr)
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
func GenerateTopShotMarketV3Contract(ftAddr, nftAddr, topshotAddr, marketAddr, ducTokenAddr, topShotLockingAddr, metadataViewsAddr string) []byte {

	marketCode := assets.MustAssetString(marketV3File)
	codeWithNFTAddr := strings.ReplaceAll(marketCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopshotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)
	codeWithFTAddr := strings.ReplaceAll(codeWithTopshotAddr, defaultFungibleTokenAddress, ftAddr)
	codeWithMarketV3Addr := strings.ReplaceAll(codeWithFTAddr, defaultMarketAddress, marketAddr)
	codeWithTokenAddr := strings.ReplaceAll(codeWithMarketV3Addr, "DUCADDRESS", ducTokenAddr)
	codeWithTopShotLockingAddr := strings.ReplaceAll(codeWithTokenAddr, defaultTopShotLockingAddress, topShotLockingAddr)
	codeWithMetadataViewAddr := strings.ReplaceAll(codeWithTopShotLockingAddr, defaultMetadataviewsAddress, metadataViewsAddr)

	return []byte(codeWithMetadataViewAddr)
}

// GenerateTopShotLockingContract returns a copy
// of the TopShotLockingContract with the import addresses updated
func GenerateTopShotLockingContract(nftAddr string) []byte {
	lockingCode := assets.MustAssetString(topShotLockingFile)
	codeWithNFTAddr := strings.ReplaceAll(lockingCode, defaultNonFungibleTokenAddress, nftAddr)

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotLockingContractWithTopShotRuntimeAddr returns a copy
// of the TopShotLockingContractWithTopShotRuntimeAddr with the import addresses updated
// the contract includes a runtime type check relying on the topshotAddr
func GenerateTopShotLockingContractWithTopShotRuntimeAddr(nftAddr string, topshotAddr string) []byte {
	lockingCode := assets.MustAssetString(topShotLockingFile)
	codeWithNFTAddr := strings.ReplaceAll(lockingCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopShotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)

	return []byte(codeWithTopShotAddr)
}

// GenerateFastBreakContract returns a copy
// of the FastBreakContract with the import addresses updated
func GenerateFastBreakContract(nftAddr string, topshotAddr string, metadataViewsAddr string) []byte {
	fastBreakCode := assets.MustAssetString(fastBreakFile)
	codeWithNFTAddr := strings.ReplaceAll(fastBreakCode, defaultNonFungibleTokenAddress, nftAddr)
	codeWithTopShotAddr := strings.ReplaceAll(codeWithNFTAddr, defaultTopshotAddress, topshotAddr)
	codeWithMetadataViewsAddr := strings.ReplaceAll(codeWithTopShotAddr, defaultMetadataviewsAddress, metadataViewsAddr)

	return []byte(codeWithMetadataViewsAddr)
}

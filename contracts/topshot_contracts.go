package contracts

// This package defines functions to return byte arrays of the
// flow core contracts files for use in go testing and deployment

import (
	"strings"

	"github.com/onflow/flow-ft/fttest"

	"github.com/onflow/flow-go-sdk"
)

const (
	topshotFile           = "./TopShot.cdc"
	topshotV1File         = "./TopShotv1.cdc"
	marketFile            = "./MarketTopShot.cdc"
	shardedCollectionFile = "./TopShotShardedCollection.cdc"
	adminReceiverFile     = "./TopShotAdminReceiver.cdc"
)

// GenerateTopShotContract returns a copy
// of the topshot contract with the import addresses updated
func GenerateTopShotContract(nftAddr flow.Address) []byte {

	topShotCode := fttest.ReadFile(topshotFile)
	codeWithNFTAddr := strings.ReplaceAll(string(topShotCode), "0x02", nftAddr.String())

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotV1Contract returns a copy
// of the original topshot contract with the import addresses updated
func GenerateTopShotV1Contract(nftAddr flow.Address) []byte {

	topShotCode := fttest.ReadFile(topshotV1File)
	codeWithNFTAddr := strings.ReplaceAll(string(topShotCode), "0x02", nftAddr.String())

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotShardedCollectionContract returns a copy
// of the TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr flow.Address) []byte {

	shardedCode := fttest.ReadFile(marketFile)
	codeWithNFTAddr := strings.ReplaceAll(string(shardedCode), "0x02", nftAddr.String())
	codeWithTopshotAddr := strings.ReplaceAll(string(codeWithNFTAddr), "0x03", topshotAddr.String())

	return []byte(codeWithTopshotAddr)
}

// GenerateTopshotAdminReceiverContract returns a copy
// of the TopshotAdminReceiver contract with the import addresses updated
func GenerateTopshotAdminReceiverContract(topshotAddr, shardedAddr flow.Address) []byte {

	adminReceiverCode := fttest.ReadFile(adminReceiverFile)
	codeWithTopshotAddr := strings.ReplaceAll(string(adminReceiverCode), "0x03", topshotAddr.String())
	codeWithShardedAddr := strings.ReplaceAll(string(codeWithTopshotAddr), "0x04", shardedAddr.String())

	return []byte(codeWithShardedAddr)
}

// GenerateTopShotMarketContract returns a copy
// of the TopShotMarketContract with the import addresses updated
func GenerateTopShotMarketContract(fungibletokenAddr, flowtokenAddr, nftAddr, topshotAddr flow.Address) []byte {

	marketCode := fttest.ReadFile(marketFile)
	codeWithFunTAddr := strings.ReplaceAll(string(marketCode), "0x04", fungibletokenAddr.String())
	codeWithFlowTAddr := strings.ReplaceAll(string(codeWithFunTAddr), "0x05", flowtokenAddr.String())
	codeWithNFTAddr := strings.ReplaceAll(string(codeWithFlowTAddr), "0x02", nftAddr.String())
	codeWithTopshotAddr := strings.ReplaceAll(string(codeWithNFTAddr), "0x03", topshotAddr.String())

	return []byte(codeWithTopshotAddr)
}

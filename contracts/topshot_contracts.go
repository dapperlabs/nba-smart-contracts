package contracts

// This package defines functions to return byte arrays of the
// flow core contracts files for use in go testing and deployment

import (
	"io/ioutil"
	"strings"

	"github.com/onflow/flow-go-sdk"
)

const (
	topshotFile           = "../contracts/TopShot.cdc"
	marketFile            = "../contracts/MarketTopShot.cdc"
	shardedCollectionFile = "../contracts/TopShotShardedCollection.cdc"
	adminReceiverFile     = "../contracts/TopShotAdminReceiver.cdc"
)

// GenerateTopShotContract returns a copy
// of the topshot contract with the import addresses updated
func GenerateTopShotContract(nftAddr flow.Address) []byte {

	topShotCode := ReadFile(topshotFile)
	codeWithNFTAddr := strings.ReplaceAll(string(topShotCode), "02", nftAddr.String())

	return []byte(codeWithNFTAddr)
}

// GenerateTopShotShardedCollectionContract returns a copy
// of the TopShotShardedCollectionContract with the import addresses updated
func GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr flow.Address) []byte {

	shardedCode := ReadFile(shardedCollectionFile)
	codeWithNFTAddr := strings.ReplaceAll(string(shardedCode), "02", nftAddr.String())
	codeWithTopshotAddr := strings.ReplaceAll(string(codeWithNFTAddr), "03", topshotAddr.String())

	return []byte(codeWithTopshotAddr)
}

// GenerateTopshotAdminReceiverContract returns a copy
// of the TopshotAdminReceiver contract with the import addresses updated
func GenerateTopshotAdminReceiverContract(topshotAddr, shardedAddr flow.Address) []byte {

	adminReceiverCode := ReadFile(adminReceiverFile)
	codeWithTopshotAddr := strings.ReplaceAll(string(adminReceiverCode), "03", topshotAddr.String())
	codeWithShardedAddr := strings.ReplaceAll(string(codeWithTopshotAddr), "04", shardedAddr.String())

	return []byte(codeWithShardedAddr)
}

// GenerateTopShotMarketContract returns a copy
// of the TopShotMarketContract with the import addresses updated
func GenerateTopShotMarketContract(fungibletokenAddr, flowtokenAddr, nftAddr, topshotAddr flow.Address) []byte {

	marketCode := ReadFile(marketFile)
	codeWithFunTAddr := strings.ReplaceAll(string(marketCode), "0x04", fungibletokenAddr.String())
	codeWithFlowTAddr := strings.ReplaceAll(string(codeWithFunTAddr), "0x05", flowtokenAddr.String())
	codeWithNFTAddr := strings.ReplaceAll(string(codeWithFlowTAddr), "0x02", nftAddr.String())
	codeWithTopshotAddr := strings.ReplaceAll(string(codeWithNFTAddr), "0x03", topshotAddr.String())

	return []byte(codeWithTopshotAddr)
}

// ReadFile reads a file from the file system
func ReadFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

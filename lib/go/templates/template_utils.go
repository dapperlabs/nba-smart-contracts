package templates

import (
	"fmt"
	"strings"

	_ "github.com/kevinburke/go-bindata"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -prefix ../../../transactions/... -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../transactions/...

const (
	placeholderFungibleTokenAddress   = "0xFUNGIBLETOKENADDRESS"
	placeholderFlowTokenAddress       = "0xFLOWTOKENADDRESS"
	placeholderNFTAddress             = "0xNFTADDRESS"
	placeholderTopShotAddress         = "0xTOPSHOTADDRESS"
	placeholderTopShotMarketAddress   = "0xMARKETADDRESS"
	placeholderTopShotMarketV3Address = "0xMARKETV3ADDRESS"
	placeholderShardedAddress         = "0xSHARDEDADDRESS"
	placeholderAdminReceiverAddress   = "0xADMINRECEIVERADDRESS"
	placeholderDUCAddress             = "0xDUCADDRESS"
	placeholderForwardingAddress      = "0xFORWARDINGADDRESS"
	placeholderMetadataViewsAddress   = "0xMETADATAVIEWSADDRESS"
	placeholderTopShotLockingAddress  = "0xTOPSHOTLOCKINGADDRESS"
	placeholderFastBreakAddress       = "0xFASTBREAKADDRESS"
	placeholderFTSwitchboardAddress   = "0xFUNGIBLETOKENSWITCHBOARDADDRESS"
)

type Environment struct {
	Network                           string
	FungibleTokenAddress              string
	FlowTokenAddress                  string
	NFTAddress                        string
	TopShotAddress                    string
	TopShotMarketAddress              string
	TopShotMarketV3Address            string
	ShardedAddress                    string
	FastBreakAddress                  string
	AdminReceiverAddress              string
	DUCAddress                        string
	ForwardingAddress                 string
	MetadataViewsAddress              string
	TopShotLockingAddress             string
	FungibleTokenMetadataViewsAddress string
	FTSwitchboardAddress              string
	ViewResolverAddress               string
	CrossVMMetadataViewsAddress       string
	EVMAddress                        string
}

func uint32ToCadenceArr(nums []uint32) []byte {
	var s string
	for _, n := range nums {
		s += fmt.Sprintf("%d as UInt32, ", n)
	}
	// slice the last 2 characters off as that's the comma and the whitespace
	return []byte("[" + s[:len(s)-2] + "]")
}

func withHexPrefix(address string) string {
	if address == "" {
		return ""
	}

	if address[0:2] == "0x" {
		return address
	}

	return fmt.Sprintf("0x%s", address)
}

func replaceAddresses(code string, env Environment) string {

	code = strings.ReplaceAll(
		code,
		placeholderFungibleTokenAddress,
		withHexPrefix(env.FungibleTokenAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderFlowTokenAddress,
		withHexPrefix(env.FlowTokenAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderNFTAddress,
		withHexPrefix(env.NFTAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderTopShotAddress,
		withHexPrefix(env.TopShotAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderTopShotMarketAddress,
		withHexPrefix(env.TopShotMarketAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderTopShotMarketV3Address,
		withHexPrefix(env.TopShotMarketV3Address),
	)

	code = strings.ReplaceAll(
		code,
		placeholderShardedAddress,
		withHexPrefix(env.ShardedAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderAdminReceiverAddress,
		withHexPrefix(env.AdminReceiverAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderDUCAddress,
		withHexPrefix(env.DUCAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderForwardingAddress,
		withHexPrefix(env.ForwardingAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderMetadataViewsAddress,
		withHexPrefix(env.MetadataViewsAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderTopShotLockingAddress,
		withHexPrefix(env.TopShotLockingAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderFastBreakAddress,
		withHexPrefix(env.FastBreakAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderFTSwitchboardAddress,
		withHexPrefix(env.FTSwitchboardAddress),
	)

	return code
}

package templates

import (
	"fmt"
	"strings"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -prefix ../../../transactions/... -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../transactions/...

const (
	defaultTopShotAddress = "TOPSHOTADDRESS"
	defaultNFTAddress     = "NFTADDRESS"
	defaultMarketAddress  = "MARKETADDRESS"
	defaultShardedAddress = "SHARDEDADDRESS"
)

func uint32ToCadenceArr(nums []uint32) []byte {
	var s string
	for _, n := range nums {
		s += fmt.Sprintf("%d as UInt32, ", n)
	}
	// slice the last 2 characters off as that's the comma and the whitespace
	return []byte("[" + s[:len(s)-2] + "]")
}

func replaceAddresses(code string, topShotAddr, nftAddr, marketAddr, shardedAddr string) string {

	replacer := strings.NewReplacer("0x"+defaultTopShotAddress, "0x"+topShotAddr,
		"0x"+defaultNFTAddress, "0x"+nftAddr,
		"0x"+defaultMarketAddress, "0x"+nftAddr,
		"0x"+defaultShardedAddress, "0x"+shardedAddr)

	code = replacer.Replace(code)

	return code
}

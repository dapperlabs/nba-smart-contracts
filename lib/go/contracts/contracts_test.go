package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
)

var addrA = "0A"
var addrB = "0B"
var addrC = "0C"
var addrD = "0D"
var addrE = "0E"
var addrF = "0F"
var addrG = "0G"
var network = "mainnet"
var flowEvmContractAddr = "0x1234565789012345657890123456578901234565"

func TestTopShotContract(t *testing.T) {
	contract := contracts.GenerateTopShotContract(addrA, addrA, addrA, addrA, addrA, addrA, addrA, addrA, network, flowEvmContractAddr)
	assert.NotNil(t, contract)
}

func TestTopShotShardedCollectionContract(t *testing.T) {
	contract := contracts.GenerateTopShotShardedCollectionContract(addrA, addrB, addrC)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

func TestTopShotAdminReceiverContract(t *testing.T) {
	contract := contracts.GenerateTopshotAdminReceiverContract(addrA, addrB)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

func TestTopShotMarketContract(t *testing.T) {
	contract := contracts.GenerateTopShotMarketContract(addrA, addrB, addrC, addrD)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

func TestTopShotMarketV3Contract(t *testing.T) {
	contract := contracts.GenerateTopShotMarketV3Contract(addrA, addrB, addrC, addrD, addrE, addrF, addrG)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

func TestFastBreakContract(t *testing.T) {
	contract := contracts.GenerateFastBreakContract(addrA, addrB, addrC, addrD)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
	assert.Contains(t, string(contract), addrB)
}

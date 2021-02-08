package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
)

var addrA = "0A"
var addrB = "0B"
var addrC = "0C"

func TestTopShotContract(t *testing.T) {
	contract := contracts.GenerateTopShotContract(addrA)
	assert.NotNil(t, contract)
}

func TestTopShotShardedCollectionContract(t *testing.T) {
	contract := contracts.GenerateTopShotShardedCollectionContract(addrA, addrB)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

func TestTopShotAdminReceiverContract(t *testing.T) {
	contract := contracts.GenerateTopshotAdminReceiverContract(addrA, addrB)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

func TestTopShotMarketContract(t *testing.T) {
	contract := contracts.GenerateTopShotMarketContract(addrA, addrB, addrC)
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA)
}

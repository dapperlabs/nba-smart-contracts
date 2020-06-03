package contracts_test

import (
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/nba-smart-contracts/contracts"
)

var addrA = flow.HexToAddress("0A")
var addrB = flow.HexToAddress("0B")

func TestTopShotContract(t *testing.T) {
	contract := contracts.GenerateTopShotContract(addrA.Hex())
	assert.NotNil(t, contract)
}

func TestTopShotShardedCollectionContract(t *testing.T) {
	contract := contracts.GenerateTopShotShardedCollectionContract(addrA.Hex(), addrB.Hex())
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA.Hex())
}

func TestTopShotAdminReceiverContract(t *testing.T) {
	contract := contracts.GenerateTopshotAdminReceiverContract(addrA.Hex(), addrB.Hex())
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA.Hex())
}

func TestTopShotMarketContract(t *testing.T) {
	contract := contracts.GenerateTopShotMarketContract(addrA.Hex(), addrB.Hex())
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addrA.Hex())
}

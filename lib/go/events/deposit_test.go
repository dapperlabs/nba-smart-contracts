package events

import (
	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_Deposit(t *testing.T) {
	id := uint64(1234)
	address := flow.HexToAddress("0x12345678")
	to := [8]byte(address)

	depositEventType := cadence.NewEventType(
		utils.TestLocation,
		"TopShot.Deposit",
		[]cadence.Field{
			{
				Identifier: "id",
				Type:       cadence.UInt64Type,
			},
			{
				Identifier: "to",
				Type:       &cadence.OptionalType{},
			},
		},
		nil,
	)

	depositEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt64(id),
		cadence.NewOptional(cadence.NewAddress(to)),
	}).WithType(depositEventType)

	payload, err := jsoncdc.Encode(depositEvent)
	require.NoError(t, err, "failed to encode deposit cadence event")

	decodedDepositEventType, err := DecodeDepositEvent(payload)
	require.NoError(t, err, "failed to decode deposit cadence event")

	assert.Equal(t, id, decodedDepositEventType.Id())
	assert.Equal(t, address.String(), decodedDepositEventType.To())
}

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

func TestCadenceEvents_Withdraw(t *testing.T) {
	id := uint64(1234)
	address := flow.HexToAddress("0x12345678")
	from := [8]byte(address)

	withdrawEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.Withdraw",
		Fields: []cadence.Field{
			{
				Identifier: "id",
				Type:       cadence.UInt64Type{},
			},
			{
				Identifier: "from",
				Type:       cadence.OptionalType{},
			},
		},
		Initializer: []cadence.Parameter{},
	}

	withdrawEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt64(id),
		cadence.NewOptional(cadence.NewAddress(from)),
	}).WithType(&withdrawEventType)

	payload, err := jsoncdc.Encode(withdrawEvent)
	require.NoError(t, err, "failed to encode withdraw cadence event")

	decodedWithdrawEventType, err := DecodeWithdrawEvent(payload)
	require.NoError(t, err, "failed to decode withdraw cadence event")

	assert.Equal(t, id, decodedWithdrawEventType.Id())
	assert.Equal(t, address.String(), decodedWithdrawEventType.From())
}

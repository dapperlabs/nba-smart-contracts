package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)


func TestCadenceEvents_SetCreated(t *testing.T) {
	setID := uint32(1234)
	series := uint32(1234)

	setCreatedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.SetCreated",
		Fields: []cadence.Field{
			{
				Identifier: "setId",
				Type:       cadence.UInt32Type{},
			},
			{
				Identifier: "series",
				Type:       cadence.UInt32Type{},
			},
		},
		Initializer: []cadence.Parameter{},
	}

	setCreatedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt32(setID),
		cadence.NewUInt32(series),
	}).WithType(&setCreatedEventType)

	payload, err := jsoncdc.Encode(setCreatedEvent)
	require.NoError(t, err, "failed to encode set created cadence event")

	decodedSetCreatedEventType, err := DecodeSetCreatedEvent(payload)
	require.NoError(t, err, "failed to decode set created cadence event")

	assert.Equal(t, setID, decodedSetCreatedEventType.SetID())
	assert.Equal(t, series, decodedSetCreatedEventType.Series())
}

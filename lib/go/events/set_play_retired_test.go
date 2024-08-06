package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_SetPlayRetired(t *testing.T) {
	setID := uint32(1234)
	playID := uint32(1234)
	numMoments := uint32(1234)

	setPlayRetiredEventType := cadence.NewEventType(
		utils.TestLocation,
		"TopShot.PlayRetiredFromSet",
		[]cadence.Field{
			{
				Identifier: "setID",
				Type:       cadence.UInt32Type,
			},
			{
				Identifier: "playID",
				Type:       cadence.UInt32Type,
			},
			{
				Identifier: "numMoments",
				Type:       cadence.UInt32Type,
			},
		},
		nil,
	)

	setPlayRetiredEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt32(setID),
		cadence.NewUInt32(playID),
		cadence.NewUInt32(numMoments),
	}).WithType(setPlayRetiredEventType)

	payload, err := jsoncdc.Encode(setPlayRetiredEvent)
	require.NoError(t, err, "failed to encode set play retired cadence event")

	decodedSetPlayRetiredEvent, err := DecodeSetPlayRetiredEvent(payload)
	require.NoError(t, err, "failed to decode set play retired cadence event")

	assert.Equal(t, setID, decodedSetPlayRetiredEvent.SetID())
	assert.Equal(t, playID, decodedSetPlayRetiredEvent.SetID())
	assert.Equal(t, numMoments, decodedSetPlayRetiredEvent.SetID())
}

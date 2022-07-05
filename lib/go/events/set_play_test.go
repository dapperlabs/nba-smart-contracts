package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)


func TestCadenceEvents_PlayAddedToSet(t *testing.T) {
	setID := uint32(1234)
	playID := uint32(1234)

	playAddedToSetEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.PlayAddedToSet",
		Fields: []cadence.Field{
			{
				Identifier: "setId",
				Type:       cadence.UInt32Type{},
			},
			{
				Identifier: "playId",
				Type:       cadence.UInt32Type{},
			},
		},
		Initializer: []cadence.Parameter{},
	}

	playAddedToSetEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt32(setID),
		cadence.NewUInt32(playID),
	}).WithType(&playAddedToSetEventType)

	payload, err := jsoncdc.Encode(playAddedToSetEvent)
	require.NoError(t, err, "failed to encode play added to set cadence event")

	decodedPlayAddedToSetEventType, err := DecodePlayAddedToSetEvent(payload)
	require.NoError(t, err, "failed to decode play added to set cadence event")

	assert.Equal(t, setID, decodedPlayAddedToSetEventType.SetID())
	assert.Equal(t, playID, decodedPlayAddedToSetEventType.PlayID())
}

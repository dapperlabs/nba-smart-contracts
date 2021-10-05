package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_PlayCreated(t *testing.T) {
	var (
		id = uint32(1234)
		playerKey = "playerID"
		playerValue = "player ID"
		teamKey = "teamAtMoment"
		teamValue = "current team"
	)

	playCreatedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.PlayCreated",
		Fields: []cadence.Field{
			{
				Identifier: "id",
				Type:       cadence.UInt32Type{},
			},
			{
				Identifier: "metadata",
				Type:       cadence.DictionaryType{},
			},
		},
		Initializer: []cadence.Parameter{},
	}

	playCreatedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt32(id),
		cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: cadence.NewString(playerKey), Value: cadence.NewString(playerValue)},
			{Key: cadence.NewString(teamKey), Value: cadence.NewString(teamValue)},
		}),
	}).WithType(&playCreatedEventType)

	payload, err := jsoncdc.Encode(playCreatedEvent)
	require.NoError(t, err, "failed to encode play created cadence event")

	decodedPlayCreatedEventType, err := DecodePlayCreatedEvent(payload)
	require.NoError(t, err, "failed to decode play created cadence event")

	assert.Equal(t, id, decodedPlayCreatedEventType.Id())
	assert.Equal(t, map[interface{}]interface{}{
		playerKey: playerValue,
		teamKey: teamValue,
	}, decodedPlayCreatedEventType.MetaData())
}

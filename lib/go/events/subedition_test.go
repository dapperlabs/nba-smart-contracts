package events

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_SubeditionCreated(t *testing.T) {
	var (
		id        = uint32(1234)
		name      = "Subedition #1"
		setKey    = "setID"
		setValue  = "1234"
		playKey   = "playID"
		playValue = "1234"
	)

	playCreatedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.SubeditionCreated",
		Fields: []cadence.Field{
			{
				Identifier: "id",
				Type:       cadence.UInt32Type{},
			},
			{
				Identifier: "name",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "metadata",
				Type:       cadence.DictionaryType{},
			},
		},
		Initializer: []cadence.Parameter{},
	}

	subeditionCreatedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt32(id),
		NewCadenceString(name),
		cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: NewCadenceString(setKey), Value: NewCadenceString(setValue)},
			{Key: NewCadenceString(playKey), Value: NewCadenceString(playValue)},
		}),
	}).WithType(&playCreatedEventType)

	payload, err := jsoncdc.Encode(subeditionCreatedEvent)
	require.NoError(t, err, "failed to encode play created cadence event")

	decodedSubeditionCreatedEventType, err := DecodeSubeditionCreatedEvent(payload)
	require.NoError(t, err, "failed to decode play created cadence event")

	assert.Equal(t, id, decodedSubeditionCreatedEventType.Id())
	assert.Equal(t, name, decodedSubeditionCreatedEventType.Name())
	assert.Equal(t, map[interface{}]interface{}{
		setKey:  setValue,
		playKey: playValue,
	}, decodedSubeditionCreatedEventType.MetaData())
}

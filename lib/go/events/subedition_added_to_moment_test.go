package events

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_SubeditionAddedToMoment(t *testing.T) {
	var (
		subeditionID = uint32(1234)
		momentID     = uint64(1234)
	)

	playCreatedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.SubeditionAddedToMoment",
		Fields: []cadence.Field{
			{
				Identifier: "momentID",
				Type:       cadence.UInt64Type{},
			},
			{
				Identifier: "subeditionID",
				Type:       cadence.UInt32Type{},
			},
		},
		Initializer: []cadence.Parameter{},
	}

	subeditionAddedToMomentEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt64(momentID),
		cadence.NewUInt32(subeditionID),
	}).WithType(&playCreatedEventType)

	payload, err := jsoncdc.Encode(subeditionAddedToMomentEvent)
	require.NoError(t, err, "failed to encode subedition added to moment cadence event")

	decodedPlayCreatedEventType, err := DecodeSubeditionAddedToMomentEvent(payload)
	require.NoError(t, err, "failed to decode subedition added to moment cadence event")

	assert.Equal(t, momentID, decodedPlayCreatedEventType.MomentID())
	assert.Equal(t, subeditionID, decodedPlayCreatedEventType.SubeditionID())
}

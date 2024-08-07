package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_MomentDestroyed(t *testing.T) {
	id := uint64(1234)

	momentDestroyedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.MomentDestroyed",
		Fields: []cadence.Field{
			{
				Identifier: "setID",
				Type:       cadence.UInt64Type{},
			},
		},
	}

	momentDestroyedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt64(id),
	}).WithType(&momentDestroyedEventType)

	payload, err := jsoncdc.Encode(momentDestroyedEvent)
	require.NoError(t, err, "failed to encode moment destroyed cadence event")

	decodeSetLockPurchasedEvent, err := DecodeMomentDestroyedEvent(payload)
	require.NoError(t, err, "failed to decode moment destroyed cadence event")

	assert.Equal(t, id, decodeSetLockPurchasedEvent.Id())
}

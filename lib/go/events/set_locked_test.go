package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_SetLocked(t *testing.T) {
	setID := uint32(1234)

	setLockedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "TopShot.SetLocked",
		Fields: []cadence.Field{
			{
				Identifier: "setID",
				Type:       cadence.UInt32Type{},
			},
		},
	}

	setLockedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt32(setID),
	}).WithType(&setLockedEventType)

	payload, err := jsoncdc.Encode(setLockedEvent)
	require.NoError(t, err, "failed to encode set locked cadence event")

	decodeSetLockPurchasedEvent, err := DecodeSetLockedEvent(payload)
	require.NoError(t, err, "failed to decode set locked cadence event")

	assert.Equal(t, setID, decodeSetLockPurchasedEvent.SetID())
}

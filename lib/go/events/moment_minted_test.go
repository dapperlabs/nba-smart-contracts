package events

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_MomentMinted(t *testing.T) {
	momentID := uint64(1234)
	playID := uint32(1234)
	setID := uint32(1234)
	serialNumber := uint32(1234)
	subeditionID := uint32(1234)

	momentMintedEventType := cadence.NewEventType(
		utils.TestLocation,
		"TopShot.MomentMinted",
		[]cadence.Field{
			{
				Identifier: "momentID",
				Type:       cadence.UInt64Type,
			},
			{
				Identifier: "playID",
				Type:       cadence.UInt32Type,
			},
			{
				Identifier: "setID",
				Type:       cadence.UInt32Type,
			},
			{
				Identifier: "serialNumber",
				Type:       cadence.UInt32Type,
			},
			{
				Identifier: "subeditionID",
				Type:       cadence.UInt32Type,
			},
		},
		nil,
	)

	momentMintedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt64(momentID),
		cadence.NewUInt32(playID),
		cadence.NewUInt32(setID),
		cadence.NewUInt32(serialNumber),
		cadence.NewUInt32(subeditionID),
	}).WithType(momentMintedEventType)

	payload, err := jsoncdc.Encode(momentMintedEvent)
	require.NoError(t, err, "failed to encode moment minted cadence event")

	decodedMomentMintedEventType, err := DecodeMomentMintedEvent(payload)
	require.NoError(t, err, "failed to decode moment minted cadence event")

	assert.Equal(t, momentID, decodedMomentMintedEventType.MomentId())
	assert.Equal(t, playID, decodedMomentMintedEventType.PlayId())
	assert.Equal(t, setID, decodedMomentMintedEventType.SetId())
	assert.Equal(t, serialNumber, decodedMomentMintedEventType.SerialNumber())
	assert.Equal(t, subeditionID, decodedMomentMintedEventType.SubeditionId())

}

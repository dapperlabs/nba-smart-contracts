package events

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestCadenceEvents_Reveal(t *testing.T) {
	var (
		packID = uint64(10)
		salt   = "salt"

		momentID1 = uint64(1)
		momentID2 = uint64(2)
		momentID3 = uint64(3)
		momentIDs = fmt.Sprintf(`%d,%d,%d`, momentID1, momentID2, momentID3)
	)

	revealedEventType := cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "PackNFT.Revealed",
		Fields: []cadence.Field{
			{
				Identifier: "id",
				Type:       cadence.UInt32Type{},
			},
			{
				Identifier: "salt",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "nfts",
				Type:       cadence.StringType{},
			},
		},
	}

	revealedEvent := cadence.NewEvent([]cadence.Value{
		cadence.NewUInt64(packID),
		NewCadenceString(salt),
		NewCadenceString(momentIDs),
	}).WithType(&revealedEventType)

	revealedPayload, err := jsoncdc.Encode(revealedEvent)
	require.NoError(t, err, "failed to encode revealed cadence event")

	decodedRevealedEventType, err := DecodeRevealedEvent(revealedPayload)
	require.NoError(t, err, "failed to decode revealed cadence event")

	assert.Equal(t, packID, decodedRevealedEventType.Id())
	assert.Equal(t, salt, decodedRevealedEventType.Salt())
	assert.Equal(t, momentIDs, decodedRevealedEventType.NFTs())

}

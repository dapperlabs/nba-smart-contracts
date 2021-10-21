package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
)

/// Used to verify set metadata in tests
type SetMetadata struct {
	setID  uint32
	name   string
	series uint32
	plays  []uint32
	//retired {UInt32: Bool}
	locked bool
	//numberMintedPerPlay {UInt32: UInt32}
}

/// Verifies that the epoch metadata is equal to the provided expected values
func verifyQuerySetMetadata(
	t *testing.T,
	b *emulator.Blockchain,
	env templates.Environment,
	expectedMetadata SetMetadata) {

	result := executeScriptAndCheck(t, b, templates.GenerateGetSetMetadataScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(expectedMetadata.setID))})
	metadataFields := result.(cadence.Struct).Fields

	setID := metadataFields[0]
	assertEqual(t, cadence.NewUInt32(expectedMetadata.setID), setID)

	name := metadataFields[1]
	assertEqual(t, CadenceString(expectedMetadata.name), name)

	series := metadataFields[2]
	assertEqual(t, cadence.NewUInt32(expectedMetadata.series), series)

	if len(expectedMetadata.plays) != 0 {
		plays := metadataFields[3].(cadence.Array).Values

		for i, play := range plays {
			expectedPlayID := cadence.NewUInt32(expectedMetadata.plays[i])
			assertEqual(t, expectedPlayID, play)
		}
	}

	locked := metadataFields[5]
	assertEqual(t, cadence.NewBool(expectedMetadata.locked), locked)

}

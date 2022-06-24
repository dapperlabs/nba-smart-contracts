package test

import (
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
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

func updateContract(b *emulator.Blockchain, address sdk.Address, signer crypto.Signer, name string, contractCode []byte) error {
	tx := sdktemplates.UpdateAccountContract(
		address,
		sdktemplates.Contract{
			Name:   name,
			Source: string(contractCode),
		},
	)

	tx.SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	err := tx.SignPayload(address, 0, signer)
	if err != nil {
		return err
	}

	serviceSigner, err := b.ServiceKey().Signer()
	if err != nil {
		return err
	}

	err = tx.SignEnvelope(b.ServiceKey().Address, b.ServiceKey().Index, serviceSigner)
	if err != nil {
		return err
	}

	err = b.AddTransaction(*tx)
	if err != nil {
		return err
	}

	result, err := b.ExecuteNextTransaction()
	if err != nil {
		return err
	}
	if !result.Succeeded() {
		return result.Error
	}

	_, err = b.CommitBlock()
	if err != nil {
		return err
	}

	return nil
}

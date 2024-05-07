package test

import (
	"testing"

	"github.com/onflow/flow-emulator/convert"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/emulator"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
)

// / Used to verify set metadata in tests
type SetMetadata struct {
	setID  uint32
	name   string
	series uint32
	plays  []uint32
	//retired {UInt32: Bool}
	locked bool
	//numberMintedPerPlay {UInt32: UInt32}
}

// / Verifies that the epoch metadata is equal to the provided expected values
func verifyQuerySetMetadata(
	t *testing.T,
	b *emulator.Blockchain,
	env templates.Environment,
	expectedMetadata SetMetadata) {

	result := executeScriptAndCheck(t, b, templates.GenerateGetSetMetadataScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(expectedMetadata.setID))})
	metadataStruct := result.(cadence.Struct)

	//setID := metadataFields[0]
	setID := cadence.SearchFieldByName(metadataStruct, "setID")
	assertEqual(t, cadence.NewUInt32(expectedMetadata.setID), setID)

	//name := metadataFields[1]
	name := cadence.SearchFieldByName(metadataStruct, "name")
	assertEqual(t, CadenceString(expectedMetadata.name), name)

	//series := metadataFields[2]
	series := cadence.SearchFieldByName(metadataStruct, "series")
	assertEqual(t, cadence.NewUInt32(expectedMetadata.series), series)

	//if len(expectedMetadata.plays) != 0 {
	//	plays := metadataFields[3].(cadence.Array).Values
	//
	//	for i, play := range plays {
	//		expectedPlayID := cadence.NewUInt32(expectedMetadata.plays[i])
	//		assertEqual(t, expectedPlayID, play)
	//	}
	//}

	//locked := metadataFields[5]
	locked := cadence.SearchFieldByName(metadataStruct, "locked")
	assertEqual(t, cadence.NewBool(expectedMetadata.locked), locked)

}

type MetadataViewStruct struct {
	CollectionBannerImage string `cadence:"collectionBannerImage"`
	CollectionDescription string `cadence:"collectionDescription"`
	CollectionName        string `cadence:"collectionName"`
	CollectionSquareImage string `cadence:"collectionSquareImage"`
	Description           string `cadence:"description"`
	ExternalURL           string `cadence:"externalURL"`
	Name                  string `cadence:"name"`
	Owner                 string
	PublicPath            string `cadence:"publicPath"`
	RoyaltyReceiversCount uint32 `cadence:"royaltyReceiversCount"`
	StoragePath           string `cadence:"storagePath"`
	Thumbnail             string `cadence:"thumbnail"`
	TraitsCount           uint32 `cadence:"traitsCount"`
	Type                  string `cadence:"type"`
	VideoURL              string `cadence:"videoURL"`
}

func getTopShotMetadata(t *testing.T, b *emulator.Blockchain, env templates.Environment, topshotAddr flow.Address, nftID uint64) MetadataViewStruct {
	metadataResult := executeScriptAndCheck(
		t,
		b,
		templates.GenerateGetNFTMetadataScript(env),
		[][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(nftID))},
	)
	metadataCDCStruct := metadataResult.(cadence.Struct)

	mvs := &MetadataViewStruct{}
	mvs.Owner = cadence.SearchFieldByName(metadataCDCStruct, "owner").(cadence.Address).String()
	err := cadence.DecodeFields(metadataCDCStruct, mvs)
	require.NoError(t, err)

	return *mvs
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

	flowTx := convert.SDKTransactionToFlow(*tx)
	err = b.AddTransaction(*flowTx)
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

// Transfer and start a v1 or v3 sale
func transferAndStartSale(
	t *testing.T,
	b *emulator.Blockchain,
	env templates.Environment,
	marketVersion int,
	momentIndex int,
	totalMoments int,
	price string,
	contractAddr sdk.Address,
	contractSigner crypto.Signer,
	userAddress sdk.Address,
	signer crypto.Signer,
	serviceKeySigner crypto.Signer) {

	startSaleScript := func() []byte {
		switch marketVersion {
		case 1:
			return templates.GenerateStartSaleScript(env)
		default:
			return templates.GenerateStartSaleV3Script(env)
		}
	}

	// transfer two moments to user's account and start sale for both moments
	for i := momentIndex; i < (momentIndex + totalMoments); i++ {
		// transfer moments to user
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), contractAddr)

		_ = tx.AddArgument(cadence.NewAddress(userAddress))
		_ = tx.AddArgument(cadence.NewUInt64(uint64(i)))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, contractAddr}, []crypto.Signer{serviceKeySigner, contractSigner},
			false,
		)
		// Start sale for user's moments
		tx = createTxWithTemplateAndAuthorizer(b, startSaleScript(), userAddress)

		_ = tx.AddArgument(cadence.NewUInt64((uint64(i))))
		_ = tx.AddArgument(CadenceUFix64(price))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, userAddress}, []crypto.Signer{serviceKeySigner, signer},
			false,
		)
	}
}

func VariableArray(cadenceType cadence.Type, values ...cadence.Value) cadence.Array {
	return cadence.NewArray(values).WithType(cadence.NewVariableSizedArrayType(cadenceType))
}

func UInt32Array(values ...int) cadence.Array {
	mapped := make([]cadence.Value, len(values))
	for i, v := range values {
		mapped[i] = cadence.NewUInt32(uint32(v))
	}
	return VariableArray(cadence.UInt32Type, mapped...)
}

func UInt64Array(values ...int) cadence.Array {
	mapped := make([]cadence.Value, len(values))
	for i, v := range values {
		mapped[i] = cadence.NewUInt64(uint64(v))
	}
	return VariableArray(cadence.UInt64Type, mapped...)
}

func CadenceStringDictionary(pairs []cadence.KeyValuePair) cadence.Dictionary {
	return cadence.NewDictionary(pairs).
		WithType(cadence.NewDictionaryType(cadence.StringType, cadence.StringType))
}

func CadenceIntArrayContains(t assert.TestingT, result cadence.Value, vals ...int) {
	intArray := result.(cadence.Array).WithType(cadence.NewVariableSizedArrayType(cadence.IntType))

	interfaceArray := intArray.Values
	assert.Equal(t, len(vals), len(interfaceArray))
	for _, intValue := range interfaceArray {
		switch intValue.Type() {
		case cadence.UInt64Type:
			assert.Contains(t, vals, int(intValue.(cadence.UInt64)))
		case cadence.IntType:
			assert.Contains(t, vals, intValue.(cadence.Int).Big().Int64())
		}
	}
}

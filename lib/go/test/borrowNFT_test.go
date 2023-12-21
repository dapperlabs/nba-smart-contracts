package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"

	"github.com/onflow/flow-go-sdk/crypto"

	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestBorrowNFTSafe(t *testing.T) {
	tb := NewTopShotTestBlockchain(t)
	tb.genericBootstrapping(t) // TODO: Should be able to find a codeGen/declarative way to set things up
	b := tb.Blockchain
	env := tb.env

	t.Run("Should return non nil if the moment id in the collection", func(t *testing.T) {
		for _, momentID := range []uint64{1, 2} {
			r, err := b.ExecuteScript(templates.GenerateBorrowNFTSafeScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(tb.userAddress)), jsoncdc.MustEncode(cadence.UInt64(momentID))})
			assert.NoError(t, err)
			assert.NoError(t, r.Error)
			expectedValue := cadence.NewBool(true)
			assert.Equal(t, expectedValue, r.Value)
		}
	})

	t.Run("Should return nil/empty optional if the moment does not exist in the collection", func(t *testing.T) {
		r, err := b.ExecuteScript(templates.GenerateBorrowNFTSafeScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(tb.userAddress)), jsoncdc.MustEncode(cadence.UInt64(3))})
		assert.NoError(t, err)
		assert.NoError(t, r.Error)
		expectedValue := cadence.NewBool(false)
		assert.Equal(t, expectedValue, r.Value)
	})
}

// genericBootstrapping should get us the blockchain in a state where we can run interesting tests against it
// will need to likely expose more of the generated ids for this to be generally useful
func (tb *topshotTestBlockchain) genericBootstrapping(t *testing.T) {
	b := tb.Blockchain
	serviceKeySigner := tb.serviceKeySigner
	topshotAddr := tb.topshotAdminAddr
	accountKeys := tb.accountKeys
	topshotSigner := tb.topshotAdminSigner
	env := tb.env

	// Create a new user account
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)
	tb.userAddress = joshAddress
	// Create moment collection
	tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")
	hayward := CadenceString("Hayward")
	antetokounmpo := CadenceString("Antetokounmpo")

	// Create plays
	lebronPlayID := uint32(1)
	haywardPlayID := uint32(2)
	antetokounmpoPlayID := uint32(3)

	for _, metadata := range [][]cadence.KeyValuePair{
		{{Key: firstName, Value: lebron}},
		{{Key: firstName, Value: hayward}},
		{{Key: firstName, Value: antetokounmpo}},
	} {
		tb.CreatePlay(t, metadata)
	}

	// Create Set
	genesisSetID := uint32(1)
	tb.CreateSet(t, "Genesis")

	// Add plays to Set
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlaysToSetScript(env), topshotAddr)

	_ = tx.AddArgument(cadence.NewUInt32(genesisSetID))

	plays := []cadence.Value{cadence.NewUInt32(lebronPlayID), cadence.NewUInt32(haywardPlayID), cadence.NewUInt32(antetokounmpoPlayID)}
	_ = tx.AddArgument(cadence.NewArray(plays))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
		false,
	)

	// Mint two moments to joshAddress
	tb.MintMoment(t, genesisSetID, lebronPlayID, joshAddress)
	tb.MintMoment(t, genesisSetID, haywardPlayID, joshAddress)

	//check that moments with ids 1 and 2 exist in josh's collection
	result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
	assert.Equal(t, cadence.NewBool(true), result)
	result = executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
	assert.Equal(t, cadence.NewBool(true), result)
}

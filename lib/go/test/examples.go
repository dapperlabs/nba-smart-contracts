package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/cmd"

	emulator "github.com/dapperlabs/flow-emulator"
)

type BlockchainAPI interface {
	emulator.BlockchainAPI
	CreateAccount(publicKeys []*flow.AccountKey, code []byte) (flow.Address, error)
}

// ReadFile reads a file from the file system
func ReadFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// DownloadFile will download a url a byte slice
func DownloadFile(url string) ([]byte, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// NewEmulator returns a emulator object for testing
func NewEmulator() BlockchainAPI {
	b, err := emulator.NewBlockchain()
	if err != nil {
		panic(err)
	}
	return b
}

// createSignAndSubmit creates a new transaction and submits it
func createSignAndSubmit(
	t *testing.T,
	b emulator.BlockchainAPI,
	template []byte,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) {

	latestBlock, err := b.GetLatestBlock()
	require.NoError(t, err)

	tx := flow.NewTransaction().
		SetScript(template).
		SetGasLimit(9999).
		SetReferenceBlockID(latestBlock.ID).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().ID, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)
	for _, addr := range signerAddresses[1:] {
		tx = tx.AddAuthorizer(addr)
	}

	SignAndSubmit(
		t, b, tx,
		signerAddresses,
		signers,
		shouldRevert,
	)
}

// SignAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// Then submits the transaction to the emulator. If the private keys don't match up with the addresses,
// the transaction will not succeed.
// shouldRevert parameter indicates whether the transaction should fail or not
// This function asserts the correct result and commits the block if it passed
func SignAndSubmit(
	t *testing.T,
	b emulator.BlockchainAPI,
	tx *flow.Transaction,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) {
	// sign transaction playload with each signer other than the first
	for i, signer := range signers[1:] {
		err := tx.SignPayload(signerAddresses[i+1], 0, signer)
		assert.NoError(t, err)
	}
	// sing transaction envelope with the first signer
	err := tx.SignEnvelope(signerAddresses[0], 0, signers[0])
	assert.NoError(t, err)

	Submit(t, b, tx, shouldRevert)
}

// Submit submits a transaction and checks
// if it fails or not
func Submit(
	t *testing.T,
	b emulator.BlockchainAPI,
	tx *flow.Transaction,
	shouldRevert bool,
) {
	// submit the signed transaction
	err := b.AddTransaction(*tx)
	require.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	require.NoError(t, err)

	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
			cmd.PrettyPrintError(result.Error, "", map[string]string{"": ""})
		}
	}

	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

// ExecuteScriptAndCheck executes a script and checks to make sure
// that it succeeded
func ExecuteScriptAndCheck(t *testing.T, b emulator.BlockchainAPI, script []byte, shouldRevert bool) {
	result, err := b.ExecuteScript(script, nil)
	if err != nil {
		t.Log(string(script))
	}
	require.NoError(t, err)
	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
			cmd.PrettyPrintError(result.Error, "", map[string]string{"": ""})
		}
	}
}

// CadenceUFix64 returns a UFix64 value
func CadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)

	if err != nil {
		panic(err)
	}

	return newValue
}

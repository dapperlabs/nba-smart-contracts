package topshottests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/onflow/flow-ft/fttest"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-nft/nfttests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/cmd"
	"github.com/onflow/flow-go-sdk"

	emulator "github.com/dapperlabs/flow-emulator"
)

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
func NewEmulator() *emulator.Blockchain {
	b, err := emulator.NewBlockchain()
	if err != nil {
		panic(err)
	}
	return b
}

// createSignAndSubmit creates a new transaction and submits it
func createSignAndSubmit(
	t *testing.T,
	b *emulator.Blockchain,
	template []byte,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(template).
		SetGasLimit(99999).
		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
		SetPayer(b.RootKey().Address).
		AddAuthorizer(signerAddresses[1])

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
	b *emulator.Blockchain,
	tx *flow.Transaction,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) {
	// sign transaction with each signer
	for i := len(signerAddresses) - 1; i >= 0; i-- {
		signerAddress := signerAddresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.SignEnvelope(signerAddress, 0, signer)
			assert.NoError(t, err)
		} else {
			err := tx.SignPayload(signerAddress, 0, signer)
			assert.NoError(t, err)
		}
	}

	Submit(t, b, tx, shouldRevert)
}

// Submit submits a transaction and checks
// if it fails or not
func Submit(
	t *testing.T,
	b *emulator.Blockchain,
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
func ExecuteScriptAndCheck(t *testing.T, b *emulator.Blockchain, script []byte) {
	result, err := b.ExecuteScript(script)
	require.NoError(t, err)
	if !assert.True(t, result.Succeeded()) {
		t.Log(result.Error.Error())
	}
}

// setupUsersTokens sets up two accounts with an empty Vault
// and a NFT collection
func setupUsersTokens(
	t *testing.T,
	b *emulator.Blockchain,
	ftAddr flow.Address,
	flowAddr flow.Address,
	nftAddr flow.Address,
	topshotAddr flow.Address,
	signerAddresses []flow.Address,
	signerKeys []*flow.AccountKey,
	signers []crypto.Signer,
) {
	// add array of signers to transaction
	for i := 0; i < len(signerAddresses); i++ {
		tx := flow.NewTransaction().
			SetScript(fttest.GenerateCreateTokenScript(ftAddr, flowAddr)).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(signerAddresses[i])

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, signerAddresses[i]},
			[]crypto.Signer{b.RootKey().Signer(), signers[i]},
			false,
		)

		tx = flow.NewTransaction().
			SetScript(nfttests.GenerateCreateCollectionScript(nftAddr, "TopShot", topshotAddr, "MomentCollection")).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(signerAddresses[i])

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, signerAddresses[i]},
			[]crypto.Signer{b.RootKey().Signer(), signers[i]},
			false,
		)
	}
}

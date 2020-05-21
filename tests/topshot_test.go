package topshottests

import (
	"testing"

	"github.com/onflow/nba-smart-contracts/data"
	"github.com/onflow/nba-smart-contracts/templates"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"
)

const (
	NonFungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	NonFungibleTokenInterfaceFile    = "NonFungibleToken.cdc"
	TopShotContractFile              = "../contracts/TopShot.cdc"
)

func TestNFTDeployment(t *testing.T) {
	b := NewEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	_, err := b.CreateAccount(nil, nftCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	topshotCode := ReadFile(TopShotContractFile)
	_, err = b.CreateAccount(nil, topshotCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

func TestMintNFTs(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	_, err := b.CreateAccount(nil, nftCode)
	assert.NoError(t, err)

	// First, deploy the contract
	topshotCode := ReadFile(TopShotContractFile)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, err := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)
	assert.NoError(t, err)

	// ExecuteScriptAndCheck(t, b, templates.GenerateInspectFieldScript(nftAddr, tokenAddr, "currentSeries", 0))
	// ExecuteScriptAndCheck(t, b, templates.GenerateInspectFieldScript(nftAddr, tokenAddr, "nextPlayID", 1))
	// ExecuteScriptAndCheck(t, b, templates.GenerateInspectFieldScript(nftAddr, tokenAddr, "nextSetID", 1))
	// ExecuteScriptAndCheck(t, b, templates.GenerateInspectFieldScript(nftAddr, tokenAddr, "totalSupply", 0))

	// create a new Collection
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Lebron"}

		template, err := templates.GenerateMintPlayScript(topshotAddr, metadata)
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to create multiple new Plays", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Oladipo"}

		template, err := templates.GenerateMintPlayScript(topshotAddr, metadata)
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)

		metadata = data.PlayMetadata{FullName: "Hayward"}

		template, err = templates.GenerateMintPlayScript(topshotAddr, metadata)
		assert.NoError(t, err)

		tx = flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to create a new Set", func(t *testing.T) {

		template, err := templates.GenerateMintSetScript(topshotAddr, "Genesis")
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to add a play to a Set", func(t *testing.T) {

		template, err := templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1)
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to add multiple plays to a Set", func(t *testing.T) {

		template, err := templates.GenerateAddPlaysToSetScript(topshotAddr, 1, []uint32{2, 3})
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to mint a moment", func(t *testing.T) {

		template, err := templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 1)
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to mint a batch of moments", func(t *testing.T) {

		template, err := templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 1, 3, 1)
		assert.NoError(t, err)

		tx := flow.NewTransaction().
			SetScript(template).
			SetGasLimit(20).
			SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
			SetPayer(b.RootKey().Address).
			AddAuthorizer(topshotAddr)

		SignAndSubmit(
			t, b, tx,
			[]flow.Address{b.RootKey().Address, topshotAddr},
			[]crypto.Signer{b.RootKey().Signer(), topshotSigner},
			false,
		)
	})
}

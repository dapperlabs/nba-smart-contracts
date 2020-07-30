package test

import (
	"strings"
	"testing"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/data"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"
)

const (
	NonFungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	NonFungibleTokenInterfaceFile    = "NonFungibleToken.cdc"
)

// This test is for testing the deployment the topshot smart contracts
func TestNFTDeployment(t *testing.T) {
	b := NewEmulator()

	// Should be able to deploy the NFT contract
	// as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, err := b.CreateAccount(nil, nftCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the topshot contract
	// as a new account with no keys.
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAddr, err := b.CreateAccount(nil, topshotCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// deploy the sharded collection contract
	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String())
	shardedAddr, err := b.CreateAccount(nil, shardedCollectionCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the admin receiver contract
	// as a new account with no keys.
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(topshotAddr.String(), shardedAddr.String())
	_, err = b.CreateAccount(nil, adminReceiverCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

// This test tests the pure functionality of the smart contract
func TestMintNFTs(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, nftCode)

	// Deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)

	// Check that that main contract fields were initialized correctly
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 0), false)

	// Deploy the sharded collection contract
	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String())
	shardedAddr, _ := b.CreateAccount(nil, shardedCollectionCode)
	_, _ = b.CommitBlock()

	// Create a new user account
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// Admin sends a transaction to create a play
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.GenerateEmptyPlay("Lebron")),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// Admin sends transactions to create multiple plays
	t.Run("Should be able to create multiple new Plays", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.GenerateEmptyPlay("Oladipo")),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.GenerateEmptyPlay("Hayward")),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.GenerateEmptyPlay("Durant")),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// Check that the return all plays script doesn't fail
		// and that we can return metadata about the plays
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnAllPlaysScript(topshotAddr), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 1, "FullName", "Lebron"), false)

		// These should fail becuase an argument is wrong for each of them
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 1, "Favorite Food", "Lebron"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 1, "FullName", "George"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 10, "FullName", "Lebron"), true)
	})

	// Admin creates a new Set with the name Genesis
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Genesis"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// Check that the set name, ID, and series were initialized correctly.
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 1, "Genesis"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Genesis", 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 1, 0), false)

		// These should fail becuase an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 5, "Genesis"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Gold", 1), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 4, 0), true)
	})

	// Admin sends a transaction that adds play 1 to the set
	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// Admin sends a transaction that adds plays 2 and 3 to the set
	t.Run("Should be able to add multiple plays to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlaysToSetScript(topshotAddr, 1, []uint32{2, 3}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		// Make sure the plays were added correctly and the edition isn't retired or locked
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 1, []int{1, 2, 3}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 1, "false"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 1, "false"), false)

		// These should fail becuase an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 3, []int{1, 2, 3}), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 2, 1, "false"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 6, "false"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 2, "false"), true)
	})

	// Admin sends a transaction that creates a new sharded collection for the admin
	t.Run("Should be able to create new sharded moment collection and store it", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupShardedCollectionScript(topshotAddr, shardedAddr, 32),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// Admin mints a moment that stores it in the admin's collection
	t.Run("Should be able to mint a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// make sure the moment was minted correctly and is stored in the collection with the correct data
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{1}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 1, 1), false)
	})

	// Admin sends a transaction that locks the set
	t.Run("Should be able to lock a set which stops plays from being added", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateLockSetScript(topshotAddr, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// This should fail because the set is locked
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 4),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)

		// Script should return that the set is locked
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 1, "true"), false)
	})

	// Admin sends a transaction that mints a batch of moments
	t.Run("Should be able to mint a batch of moments", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 1, 3, 5),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// Ensure that the correct number of moments have been minted for each edition
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 1, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 3, 5), false)

		// Ensure that the admin's collection and data is correct
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{1, 2, 3, 4, 5, 6}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 1, 1), false)

		// These should fail because an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 10, 1), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 2, 1, 1), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 6, 5), true)
	})

	// Admin sends a transaction to retire a play
	t.Run("Should be able to retire a Play which stops minting", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateRetirePlayScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// Minting from this play should fail becuase it is retired
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)

		// Make sure this edition is retired
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 1, "true"), false)
	})

	// Admin sends a transaction that retires all the plays in a set
	t.Run("Should be able to retire all Plays which stops minting", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateRetireAllPlaysScript(topshotAddr, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// minting should fail
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 3),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// create a new Collection for a user address
	t.Run("Should be able to create a moment Collection", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	// Admin sends a transaction to transfer a moment to a user
	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		// make sure the user received it
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 1), false)
	})

	// Admin sends a transaction to transfer a batch of moments to a user
	t.Run("Should be able to batch transfer moments from a sharded collection", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateBatchTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, []uint64{2, 3, 4}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		// make sure the user received them
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, joshAddress, 2, 1), false)
	})

	// Admin sends a transaction to update the current series
	t.Run("Should be able to change the current series", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateChangeSeriesScript(topshotAddr),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// TODO
	// Test cases:
	// - Mismatched set + play array lengths
	// - 0-length set + plays
	t.Run("Should fail with mismatched Set and Play slice lengths", func(t *testing.T) {
		_, err := templates.GenerateSetPlaysOwnedByAddressScript(topshotAddr, joshAddress, []uint32{1, 2}, []uint32{1})
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "mismatched lengths"))
	})

	// - 2/2
	t.Run("Should be able check how many of a user's moments contribute towards a set Challenge", func(t *testing.T) {
		challengeScript, err := templates.GenerateSetPlaysOwnedByAddressScript(topshotAddr, joshAddress, []uint32{1, 2}, []uint32{1, 2})
		assert.NoError(t, err)

		result, err := b.ExecuteScript(challengeScript, nil)
		assert.NoError(t, err)
		boolResult, ok := result.Value.ToGoValue().(bool)
		assert.True(t, ok)
		assert.True(t, boolResult)
	})

	// - 2/3
	t.Run("Should be able check how many of a user's moments contribute towards a set Challenge", func(t *testing.T) {
		challengeScript, err := templates.GenerateSetPlaysOwnedByAddressScript(topshotAddr, joshAddress, []uint32{1, 2, 3}, []uint32{1, 2, 3})
		assert.NoError(t, err)

		result, err := b.ExecuteScript(challengeScript, nil)
		assert.NoError(t, err)
		boolResult, ok := result.Value.ToGoValue().(bool)
		assert.True(t, ok)
		assert.False(t, boolResult)
	})

	// Make sure the contract fields are correct
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 5), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 6), false)

}

// This test is for ensuring that admin receiver smart contract works correctly
func TestTransferAdmin(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, nftCode)

	// First, deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String())
	shardedAddr, _ := b.CreateAccount(nil, shardedCollectionCode)
	_, _ = b.CommitBlock()

	// Should be able to deploy the admin receiver contract
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(topshotAddr.String(), shardedAddr.String())
	adminAccountKey, adminSigner := accountKeys.NewWithSigner()
	adminAddr, _ := b.CreateAccount([]*flow.AccountKey{adminAccountKey}, adminReceiverCode)
	b.CommitBlock()

	// create a new Collection
	t.Run("Should be able to transfer an admin Capability to the receiver account", func(t *testing.T) {

		createSignAndSubmit(
			t, b,
			templates.GenerateTransferAdminScript(topshotAddr, adminAddr),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// cannot create a new play with the old admin
	t.Run("Shouldn't be able to create a new Play with the old admin account", func(t *testing.T) {
		metadata := data.GenerateEmptyPlay("Lebron")

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// can create a new play with the new admin
	t.Run("Should be able to create a new Play with the new Admin account", func(t *testing.T) {
		metadata := data.GenerateEmptyPlay("Lebron")

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, adminAddr}, []crypto.Signer{b.ServiceKey().Signer(), adminSigner},
			false,
		)
	})
}

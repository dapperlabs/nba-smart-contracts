package test

import (
	"testing"

	"github.com/dapperlabs/nba-smart-contracts/contracts"
	"github.com/dapperlabs/nba-smart-contracts/templates"
	"github.com/dapperlabs/nba-smart-contracts/templates/data"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"
)

const (
	NonFungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/src/contracts/"
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
			templates.GenerateMintPlayScript(topshotAddr, data.PlayMetadata{FullName: "Lebron"}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// Admin sends transactions to create multiple plays
	t.Run("Should be able to create multiple new Plays", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.PlayMetadata{FullName: "Oladipo"}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.PlayMetadata{FullName: "Hayward"}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.PlayMetadata{FullName: "Durant"}),
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

	// Make sure the contract fields are correct
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 5), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 6), false)

}

// This test is similar to the last one,
// but it upgrades the topshot contract after
// and checks if the normal actions are still possible
func TestUpgradeTopshot(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, nftCode)

	// First, deploy the original version of the topshot contract
	topshotCode := contracts.GenerateTopShotV1Contract(nftAddr.String())
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)

	// check the contract fields initialization
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 0), false)

	// deploy the original version of the sharded collection contract
	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionV1Contract(nftAddr.String(), topshotAddr.String())
	shardedCollectionAccountKey, shardedCollectionSigner := accountKeys.NewWithSigner()
	shardedAddr, err := b.CreateAccount([]*flow.AccountKey{shardedCollectionAccountKey}, shardedCollectionCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// create a new user account
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// create a new play
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Lebron"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new set
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Genesis"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// add the play to the set
	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new sharded collection
	t.Run("Should be able to create new sharded moment collection and store it", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupShardedCollectionScript(topshotAddr, shardedAddr, 32),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// mint a moment
	t.Run("Should be able to mint a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{1}), false)
	})

	// lock a set
	t.Run("Should be able to lock a set which stops plays from being added", func(t *testing.T) {

		createSignAndSubmit(
			t, b,
			templates.GenerateLockSetScript(topshotAddr, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 4),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// retire a play
	t.Run("Should be able to retire a Play which stops minting", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateRetirePlayScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// create a new Collection for the user
	t.Run("Should be able to create a moment Collection", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	// transfer a moment from the admin to the user
	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// make sure the moment was transferred correctly
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 1), false)

	// make sure the contract fields are what they should be
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 1), false)

	// Update the topshot account with the upgraded topshot contract code
	// without overwriting any of its state
	t.Run("Should be able to upgrade the topshot code and sharded Code without resetting its fields", func(t *testing.T) {
		// get the contract code for the upgraded contract
		topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
		// submit the transaction which upgrades topshot
		createSignAndSubmit(
			t, b,
			templates.GenerateUnsafeNotInitializingSetCodeScript(topshotCode),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// get the upgraded sharded collection contract code
		shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String())
		// submit the transaction that upgrades the contract
		createSignAndSubmit(
			t, b,
			templates.GenerateUnsafeNotInitializingSetCodeScript(shardedCollectionCode),
			[]flow.Address{b.ServiceKey().Address, shardedAddr}, []crypto.Signer{b.ServiceKey().Signer(), shardedCollectionSigner},
			false,
		)

		// Check to see that the collections still have the correct moments in them
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 1), false)

		// Run scripts from the old smart contract to see that they still work correctly and return
		// the same values as before
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 1), false)

		// Run the new scripts from the updated topshot contract to see that they return the correct values
		// based the state changes that happened before the upgrade
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 1, "Genesis"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Genesis", 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 1, 0), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 1, []int{1}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 1, "true"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 1, "true"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 1, 1), false)

		// New script from the updated sharded collection contract to make sure the
		// shardedcollection still works
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, joshAddress, 1, 1), false)

		// These should fail becuase an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 5, "Genesis"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Gold", 1), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 4, 0), true)

	})

	// Ensure the the locking that happened before the upgrade still is enforced
	t.Run("Should not be able to add a play or mint a moment from the last set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// Cannot create an empty play with no data
	t.Run("Should not be able to create an empty Play", func(t *testing.T) {
		metadata := data.PlayMetadata{}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// Create a new play
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Jordan"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		// Ensure the metadata is correct
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 2, "FullName", "Jordan"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataByFieldScript(topshotAddr, 2, "FullName", "Jordan"), false)
	})

	// Admin tries to create a set with no name, which should fail
	t.Run("Should not be able to create a new Set with no Name", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, ""),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// Admin creates a new set
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Gold"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// Try to add a nonexistant play to a set which should fail
	t.Run("Should not be able to add a play that doesn't exist to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 2, 5),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// Add the correct play to the set
	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// Cannot add a play twice to a set
	t.Run("Should not be able to add a play to a Set if it has already been added", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// admin tries to mint an invalid play/set combo
	t.Run("Shouldn't be able to mint a moment for a play that doesn't exist in a set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 2, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)

		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 2), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 2, 2), true)
	})

	// admin mints a moment
	t.Run("Should be able to mint moments", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// Make sure the moment was minted correctly
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{2}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 2, 2), false)

		// mint a batch of moments
		createSignAndSubmit(
			t, b,
			templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 2, 2, 5),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// ensure the batch was minted correctly
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 5), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{2, 3, 4, 5, 6, 7}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 7, 2), false)
	})

	// admin locks a set
	t.Run("Should be able to lock a set which stops plays from being added", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateLockSetScript(topshotAddr, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// cannot add a play after a set has been locked
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 2, 4),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// admin retires a play
	t.Run("Should be able to retire a Play which stops minting", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateRetirePlayScript(topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// reture an invalid play which should fail
		createSignAndSubmit(
			t, b,
			templates.GenerateRetirePlayScript(topshotAddr, 2, 9),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)

		// cannot mint a moment from a retired play
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// transfer a moment from the admin to a user
	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// make sure the moment was transferred correctly
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, joshAddress, []uint64{1, 2}), false)
	})

	// Check the data from all the plays and sets to ensure the state stil remains correct
	// after upgrading and running transactions
	ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 2, "Gold"), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Gold", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 2, 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 2, []int{2}), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 2, 2, "true"), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 2, "true"), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 2, 2, 6), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, joshAddress, 2, 2), false)

	// These should fail because an argument is wrong
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 10, 1), true)

	// Create a new account
	ericAccountKey, ericSigner := accountKeys.NewWithSigner()
	ericAddress, _ := b.CreateAccount([]*flow.AccountKey{ericAccountKey}, nil)

	// Transfer a moment from one user to another
	t.Run("Should be able to transfer a moment again", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, ericAddress}, []crypto.Signer{b.ServiceKey().Signer(), ericSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentScript(nftAddr, topshotAddr, ericAddress, 2),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		// ensure the transfer happened correctly
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, joshAddress, []uint64{1}), false)

		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, ericAddress, 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, ericAddress, []uint64{2}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, ericAddress, 2, 2), false)
	})
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
		metadata := data.PlayMetadata{FullName: "Lebron"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// can create a new play with the new admin
	t.Run("Should be able to create a new Play with the new Admin account", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Lebron"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, adminAddr}, []crypto.Signer{b.ServiceKey().Signer(), adminSigner},
			false,
		)
	})
}

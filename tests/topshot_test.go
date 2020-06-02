package tests

import (
	"testing"

	"github.com/dapperlabs/nba-smart-contracts/contracts"
	"github.com/dapperlabs/nba-smart-contracts/data"
	"github.com/dapperlabs/nba-smart-contracts/templates"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"
)

const (
	NonFungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/src/contracts/"
	NonFungibleTokenInterfaceFile    = "NonFungibleToken.cdc"
)

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
	topshotCode := contracts.GenerateTopShotContract(nftAddr)
	topshotAddr, err := b.CreateAccount(nil, topshotCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr)
	shardedAddr, err := b.CreateAccount(nil, shardedCollectionCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the admin receiver contract
	// as a new account with no keys.
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(topshotAddr, shardedAddr)
	_, err = b.CreateAccount(nil, adminReceiverCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// topshotCode := contracts.GenerateTopShotContract(nftAddr)
	// _, err = b.UnsafeAccountCodeUpdate(nil, topshotCode)
	// if !assert.NoError(t, err) {
	// 	t.Log(err.Error())
	// }
	// _, err = b.CommitBlock()
	// assert.NoError(t, err)
}

func TestMintNFTs(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, nftCode)

	// First, deploy the contract
	topshotCode := contracts.GenerateTopShotContract(nftAddr)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)

	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 0), false)

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr)
	shardedAddr, _ := b.CreateAccount(nil, shardedCollectionCode)
	_, _ = b.CommitBlock()

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// create a new Collection
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.PlayMetadata{FullName: "Lebron"}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
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

		ExecuteScriptAndCheck(t, b, templates.GenerateReturnAllPlaysScript(topshotAddr), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 1, "FullName", "Lebron"), false)

		// These should fail becuase an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 1, "Favorite Food", "Lebron"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 1, "FullName", "George"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlayMetadataScript(topshotAddr, 10, "FullName", "Lebron"), true)
	})

	// create a new Collection
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Genesis"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 1, "Genesis"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Genesis", 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 1, 0), false)

		// These should fail becuase an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 5, "Genesis"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Gold", 1), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 4, 0), true)
	})

	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to add multiple plays to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlaysToSetScript(topshotAddr, 1, []uint32{2, 3}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 1, []int{1, 2, 3}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 1, "false"), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 1, "false"), false)

		// These should fail becuase an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 3, []int{1, 2, 3}), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 2, 1, "false"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 6, "false"), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 2, "false"), true)
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

		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 1, "true"), false)
	})

	// create a new Collection
	t.Run("Should be able to mint a batch of moments", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 1, 3, 5),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 1, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 3, 5), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{1, 2, 3, 4, 5, 6}), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, topshotAddr, 1, 1), false)

		// These should fail because an argument is wrong
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 2, 1, 1), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 6, 5), true)
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

		ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 1, "true"), false)
	})

	// retire all plays
	t.Run("Should be able to retire all Plays which stops minting", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateRetireAllPlaysScript(topshotAddr, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 1, 3),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// create a new Collection
	t.Run("Should be able to create a moment Collection", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 1), false)
	})

	t.Run("Should be able to batch transfer moments from a sharded collection", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateBatchTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, []uint64{2, 3, 4}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, joshAddress, 2, 1), false)
	})

	t.Run("Should be able to change the current series", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateChangeSeriesScript(topshotAddr),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

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

	// First, deploy the contract
	topshotCode := contracts.GenerateTopShotV1Contract(nftAddr)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)

	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 0), false)

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionV1Contract(nftAddr, topshotAddr)
	shardedAddr, err := b.CreateAccount(nil, shardedCollectionCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// create a new Collection
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Lebron"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Genesis"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

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

	ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 1), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{1}), false)

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

	// create a new Collection
	t.Run("Should be able to create a moment Collection", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 1), false)

	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 2), false)
	ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 1), false)

	// Update the topshot account with the upgraded topshot contract code
	// without overwriting any of its state
	t.Run("Should be able to upgrade the topshot code and sharded Code without resetting its fields", func(t *testing.T) {
		// topshotCode := contracts.GenerateTopShotontract(nftAddr)
		// _, err = b.GenerateUnsafeNotInitializingSetCodeScript(nil, topshotCode)
		// if !assert.NoError(t, err) {
		// 	t.Log(err.Error())
		// }
		// _, err = b.CommitBlock()
		// assert.NoError(t, err)

		// shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr)
		// _, err = b.GenerateUnsafeNotInitializingSetCodeScript(nil, shardedCollectionCode)
		// if !assert.NoError(t, err) {
		// 	t.Log(err.Error())
		// }
		// _, err = b.CommitBlock()
		// assert.NoError(t, err)

		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 1), false)

		// Old scripts
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "currentSeries", "UInt32", 0), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextPlayID", "UInt32", 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "nextSetID", "UInt32", 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectTopshotFieldScript(nftAddr, topshotAddr, "totalSupply", "UInt64", 1), false)

		// New scripts from the updated topshot contract
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 1, "Genesis"), false)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Genesis", 1), false)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 1, 0), false)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 1, []int{1}), false)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 1, 1, "true"), false)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 1, "true"), false)
		// ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 1, 1, 1), false)

		// New script from the updated sharded collection contract
		// ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionDataScript(nftAddr, topshotAddr, joshAddr, 1, 1), false)

		// // These should fail becuase an argument is wrong
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 5, "Genesis"), true)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Gold", 1), true)
		// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 4, 0), true)

	})

	// lock a set
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

	// create a new Collection
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Jordan"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// create a new Collection
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Gold"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// mint a moment
	t.Run("Should be able to mint a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, topshotAddr, 2), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionIDsScript(nftAddr, topshotAddr, topshotAddr, []uint64{2}), false)
	})

	// lock a set
	t.Run("Should be able to lock a set which stops plays from being added", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateLockSetScript(topshotAddr, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 2, 4),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// retire a play
	t.Run("Should be able to retire a Play which stops minting", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateRetirePlayScript(topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintMomentScript(topshotAddr, topshotAddr, 2, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentfromShardedCollectionScript(nftAddr, topshotAddr, shardedAddr, joshAddress, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetNameScript(topshotAddr, 2, "Gold"), false)
	// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetIDsByNameScript(topshotAddr, "Gold", 2), false)
	// ExecuteScriptAndCheck(t, b, templates.GenerateReturnSetSeriesScript(topshotAddr, 2, 0), false)
	// ExecuteScriptAndCheck(t, b, templates.GenerateReturnPlaysInSetScript(topshotAddr, 2, []int{2}), false)
	// ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsEditionRetiredScript(topshotAddr, 2, 2, "true"), false)
	// ExecuteScriptAndCheck(t, b, templates.GenerateReturnIsSetLockedScript(topshotAddr, 2, "true"), false)
	// ExecuteScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(topshotAddr, 2, 2, 1), false)

}

func TestTransferAdmin(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, nftCode)

	// First, deploy the contract
	topshotCode := contracts.GenerateTopShotContract(nftAddr)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr, topshotAddr)
	shardedAddr, _ := b.CreateAccount(nil, shardedCollectionCode)
	_, _ = b.CommitBlock()

	// Should be able to deploy a contract as a new account with no keys.
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(topshotAddr, shardedAddr)
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

	// create a new Collection
	t.Run("Shouldn't be able to create a new Play with the old admin account", func(t *testing.T) {
		metadata := data.PlayMetadata{FullName: "Lebron"}

		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, metadata),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			true,
		)
	})

	// create a new Collection
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

package test

import (
	"context"
	"strings"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/common"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/adapters"
	"github.com/onflow/flow-emulator/emulator"
	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/rs/zerolog"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"

	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	NonFungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/0ceed9bbdaba549328a73c4fca9f30a9e56c49bc/contracts/"
	NonFungibleTokenInterfaceFile    = "NonFungibleToken.cdc"

	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/0ceed9bbdaba549328a73c4fca9f30a9e56c49bc/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"

	FungibleTokenMetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-ft/9a9073c713583061cf614efc4aed0a1572ca0f09/contracts/"
	FungibleTokenMetadataViewsInterfaceFile    = "FungibleTokenMetadataViews.cdc"

	FungibleTokenSwitchboardContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-ft/9a9073c713583061cf614efc4aed0a1572ca0f09/contracts/"
	FungibleTokenSwitchboardInterfaceFile    = "FungibleTokenSwitchboard.cdc"

	ViewResolverContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/0ceed9bbdaba549328a73c4fca9f30a9e56c49bc/contracts/"
	ViewResolverInterfaceFile    = "ViewResolver.cdc"

	emulatorFTAddress           = "ee82856bf20e2aa6"
	emulatorFlowTokenAddress    = "0ae53cb6e3f42a79"
	MetadataFTReplaceAddress    = `"FungibleToken"`
	MetadataNFTReplaceAddress   = `"NonFungibleToken"`
	MetadataViewsReplaceAddress = `"MetadataViews"`
	ViewResolverReplaceAddress  = `"ViewResolver"`
	Network                     = `"mainnet"`
)

// topshotTestContext will expose sugar for common actions needed to bootstrap testing
type topshotTestBlockchain struct {
	*emulator.Blockchain
	env                            templates.Environment
	topshotAdminAddr               flow.Address
	evmAddr                        flow.Address
	crossVMMetadataViewsAddr       flow.Address
	metadataViewsAddr              flow.Address
	fungibleTokenMetadataViewsAddr flow.Address
	viewResolverAddr               flow.Address
	topshotLockingAddr             flow.Address
	serviceKeySigner               crypto.Signer
	topshotAdminSigner             crypto.Signer
	lockingSigner                  crypto.Signer
	accountKeys                    *test.AccountKeys

	userAddress flow.Address
}

func NewTopShotTestBlockchain(t *testing.T) topshotTestBlockchain {
	b := newBlockchain()

	serviceKeySigner, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	accountKeys := test.AccountKeyGenerator()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	viewResolverAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.ViewResolverAddress = viewResolverAddr.String()

	nftAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.NFTAddress = nftAddr.String()

	metadataViewsAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.MetadataViewsAddress = metadataViewsAddr.String()

	fungibleTokenMetadataViewsAddr := flow.HexToAddress("ee82856bf20e2aa6")
	env.FungibleTokenMetadataViewsAddress = fungibleTokenMetadataViewsAddr.String()

	evmAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.EVMAddress = evmAddr.String()

	// Deploy CrossVMMetadataViews contract
	crossVMMetadataViewsCode := contracts.GenerateCrossVMMetadataViewsContract(evmAddr.String(), viewResolverAddr.String())
	crossVMMetadataViewsAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "CrossVMMetadataViews",
			Source: string(crossVMMetadataViewsCode),
		},
	})
	require.NoError(t, err)
	env.CrossVMMetadataViewsAddress = crossVMMetadataViewsAddr.String()

	// Deploy TopShot Locking contract
	lockingKey, lockingSigner := test.AccountKeyGenerator().NewWithSigner()
	topshotLockingCode := contracts.GenerateTopShotLockingContract(nftAddr.String())
	topShotLockingAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{lockingKey}, []sdktemplates.Contract{
		{
			Name:   "TopShotLocking",
			Source: string(topshotLockingCode),
		},
	})
	assert.Nil(t, err)

	env.TopShotLockingAddress = topShotLockingAddr.String()

	royaltyAddr := flow.HexToAddress("ee82856bf20e2aa6")
	env.FTSwitchboardAddress = royaltyAddr.String()
	_, err = b.CommitBlock()
	require.NoError(t, err)

	// Deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(
		emulatorFTAddress,
		nftAddr.String(),
		metadataViewsAddr.String(),
		viewResolverAddr.String(),
		crossVMMetadataViewsAddr.String(),
		evmAddr.String(),
		topShotLockingAddr.String(),
		royaltyAddr.String(),
		Network,
	)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})
	require.NoError(t, err)

	env.TopShotAddress = topshotAddr.String()

	// Update the locking contract with topshot address
	topShotLockingCodeWithRuntimeAddr := contracts.GenerateTopShotLockingContractWithTopShotRuntimeAddr(nftAddr.String(), topshotAddr.String())
	err = updateContract(b, topShotLockingAddr, lockingSigner, "TopShotLocking", topShotLockingCodeWithRuntimeAddr)
	require.NoError(t, err)

	// Grant the TopShot account a TopShotLocking admin
	// In testnet/mainnet the TopShotLocking contract is in the same account as TopShot contract
	tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateTopShotLockingAdminGrantAdminScript(env), topShotLockingAddr)
	tx.AddAuthorizer(topshotAddr)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, topShotLockingAddr, topshotAddr},
		[]crypto.Signer{serviceKeySigner, lockingSigner, topshotSigner},
		false,
	)

	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupSwitchboardScript(env), royaltyAddr)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, royaltyAddr},
		[]crypto.Signer{serviceKeySigner, serviceKeySigner},
		false,
	)

	// Check that that main contract fields were initialized correctly
	result := executeScriptAndCheck(t, b, templates.GenerateGetSeriesScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(0), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetNextPlayIDScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(1), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetNextSetIDScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(1), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetSupplyScript(env), nil)
	assert.Equal(t, cadence.NewUInt64(0), result)

	// Deploy the sharded collection contract
	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String(), viewResolverAddr.String())
	shardedAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TopShotShardedCollection",
			Source: string(shardedCollectionCode),
		},
	})
	require.NoError(t, err)

	_, err = b.CommitBlock()

	require.NoError(t, err)

	env.ShardedAddress = shardedAddr.String()

	return topshotTestBlockchain{
		Blockchain:                     b,
		env:                            env,
		topshotAdminAddr:               topshotAddr,
		evmAddr:                        evmAddr,
		crossVMMetadataViewsAddr:       crossVMMetadataViewsAddr,
		metadataViewsAddr:              metadataViewsAddr,
		topshotLockingAddr:             topShotLockingAddr,
		fungibleTokenMetadataViewsAddr: fungibleTokenMetadataViewsAddr,
		viewResolverAddr:               viewResolverAddr,
		serviceKeySigner:               serviceKeySigner,
		topshotAdminSigner:             topshotSigner,
		lockingSigner:                  lockingSigner,
		accountKeys:                    accountKeys}
}

// This test is for testing the deployment the topshot smart contracts
func TestNFTDeployment(t *testing.T) {
	// all core contracts should be deployed on the test blockchain
	b := NewTopShotTestBlockchain(t)
	env := b.env
	// Should be able to deploy the admin receiver contract
	// as a new account with no keys.
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(env.TopShotAddress, env.ShardedAddress)
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	_, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TopshotAdminReceiver",
			Source: string(adminReceiverCode),
		},
	})
	require.NoError(t, err)

	_, err = b.CommitBlock()
	require.NoError(t, err)
}

// This test tests the pure functionality of the smart contract
func TestMintNFTs(t *testing.T) {
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	env := tb.env
	accountKeys := tb.accountKeys
	topshotAddr := tb.topshotAdminAddr
	serviceKeySigner := tb.serviceKeySigner
	topshotSigner := tb.topshotAdminSigner

	// Create a new user account
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	joshAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{joshAccountKey}, nil)

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")
	oladipo := CadenceString("Oladipo")
	hayward := CadenceString("Hayward")
	durant := CadenceString("Durant")

	playType := CadenceString("PlayType")
	dunk := CadenceString("Dunk")

	// Admin sends a transaction to create a play
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		tb.CreatePlay(t, []cadence.KeyValuePair{{Key: firstName, Value: lebron}, {Key: playType, Value: dunk}})
	})

	t.Run("Should be able to update an existing Play", func(t *testing.T) {
		result := executeScriptAndCheck(t, b, templates.GenerateGetPlayMetadataFieldScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.String("FullName"))})
		assert.Equal(t, CadenceString("Lebron"), result)

		tb.UpdateTagline(t, []cadence.KeyValuePair{{Key: cadence.UInt32(1), Value: CadenceString("lorem ipsum")}})
		result = executeScriptAndCheck(t, b, templates.GenerateGetPlayMetadataFieldScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.String("Tagline"))})
		assert.Equal(t, CadenceString("lorem ipsum"), result)
	})

	// Admin sends transactions to create multiple plays
	t.Run("Should be able to create multiple new Plays", func(t *testing.T) {
		metadatas := [][]cadence.KeyValuePair{
			{{Key: firstName, Value: oladipo}},
			{{Key: firstName, Value: hayward}},
			{{Key: firstName, Value: durant}},
		}

		for _, metadata := range metadatas {
			tb.CreatePlay(t, metadata)
		}
		// Check that the return all plays script doesn't fail
		// and that we can return metadata about the plays
		executeScriptAndCheck(t, b, templates.GenerateGetAllPlaysScript(env), nil)

		result := executeScriptAndCheck(t, b, templates.GenerateGetPlayMetadataFieldScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.String("FullName"))})
		assert.Equal(t, CadenceString("Lebron"), result)
	})

	// Admin creates a new Set with the name Genesis
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		tb.CreateSet(t, "Genesis")

		// Check that the set name, ID, and series were initialized correctly.
		result := executeScriptAndCheck(t, b, templates.GenerateGetSetNameScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, CadenceString("Genesis"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSetIDsByNameScript(env), [][]byte{jsoncdc.MustEncode(cadence.String("Genesis"))})
		idsArray := UInt32Array(1)
		assert.Equal(t, idsArray, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSetSeriesScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewUInt32(0), result)
	})

	t.Run("Should not be able to create set and play data structs that increment the id counter", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSetandPlayDataScript(env), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Check that the play ID and set ID were not incremented
		result := executeScriptAndCheck(t, b, templates.GenerateGetNextPlayIDScript(env), nil)
		assert.Equal(t, cadence.NewUInt32(5), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNextSetIDScript(env), nil)
		assert.Equal(t, cadence.NewUInt32(2), result)

	})

	// Admin sends a transaction that adds play 1 to the set
	t.Run("Should be able to add a play to a Set", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	// Admin sends a transaction that adds plays 2 and 3 to the set
	t.Run("Should be able to add multiple plays to a Set", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlaysToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))

		plays := []cadence.Value{cadence.NewUInt32(2), cadence.NewUInt32(3)}
		_ = tx.AddArgument(cadence.NewArray(plays))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Make sure the plays were added correctly and the edition isn't retired or locked
		result := executeScriptAndCheck(t, b, templates.GenerateGetPlaysInSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})
		playsArray := UInt32Array(1, 2, 3)
		assert.Equal(t, playsArray, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetIsEditionRetiredScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewBool(false), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetIsSetLockedScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewBool(false), result)

	})

	// Admin sends a transaction that creates a new sharded collection for the admin
	t.Run("Should be able to create new sharded moment collection and store it", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupShardedCollectionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt64(32))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	// Admin mints a moment that stores it in the admin's collection
	t.Run("Should be able to mint a moment", func(t *testing.T) {
		tb.MintMoment(t, 1, 1, topshotAddr)

		// make sure the moment was minted correctly and is stored in the collection with the correct data
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCollectionIDsScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr))})
		idsArray := UInt64Array(1)
		assert.Equal(t, idsArray, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

	})

	t.Run("Should be able to get moments metadata", func(t *testing.T) {
		// Tests to ensure that all core metadataviews are resolvable
		expectedMetadataName := "Lebron Dunk"
		expectedMetadataDescription := "lorem ipsum"
		expectedMetadataThumbnail := "https://assets.nbatopshot.com/media/1?width=256"
		expectedMetadataExternalURL := "https://nbatopshot.com/moment/1"
		expectedStoragePath := "/storage/MomentCollection"
		expectedPublicPath := "/public/MomentCollection"
		expectedCollectionName := "NBA-Top-Shot"
		expectedCollectionDescription := "NBA Top Shot is your chance to own, sell, and trade official digital collectibles of the NBA and WNBA's greatest plays and players"
		expectedCollectionSquareImage := "https://nbatopshot.com/static/favicon/favicon.svg"
		expectedCollectionBannerImage := "https://nbatopshot.com/static/img/top-shot-logo-horizontal-white.svg"
		expectedRoyaltyReceiversCount := 1
		expectedTraitsCount := 9
		expectedVideoURL := "https://assets.nbatopshot.com/media/1/video"

		mvs := getTopShotMetadata(t, b, env, topshotAddr, 1)

		assert.Equal(t, expectedMetadataName, mvs.Name)
		assert.Equal(t, expectedMetadataDescription, mvs.Description)
		assert.Equal(t, expectedMetadataThumbnail, mvs.Thumbnail)
		assert.Equal(t, expectedMetadataExternalURL, mvs.ExternalURL)
		assert.Equal(t, expectedStoragePath, mvs.StoragePath)
		assert.Equal(t, expectedPublicPath, mvs.PublicPath)
		assert.Equal(t, expectedCollectionName, mvs.CollectionName)
		assert.Equal(t, expectedCollectionDescription, mvs.CollectionDescription)
		assert.Equal(t, expectedCollectionSquareImage, mvs.CollectionSquareImage)
		assert.Equal(t, expectedCollectionBannerImage, mvs.CollectionBannerImage)
		assert.Equal(t, uint32(expectedRoyaltyReceiversCount), mvs.RoyaltyReceiversCount)
		assert.Equal(t, uint32(expectedTraitsCount), mvs.TraitsCount)
		assert.Equal(t, expectedVideoURL, mvs.VideoURL)

		// Tests that top-shot specific metadata is discoverable on-chain
		expectedPlayID := 1
		expectedSetID := 1
		expectedSerialNumber := 1

		resultTopShot := executeScriptAndCheck(t, b, templates.GenerateGetTopShotMetadataScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		metadataViewTopShot := resultTopShot.(cadence.Struct)
		assert.Equal(t, cadence.UInt32(expectedSerialNumber), cadence.SearchFieldByName(metadataViewTopShot, "serialNumber").(cadence.UInt32))
		assert.Equal(t, cadence.UInt32(expectedPlayID), cadence.SearchFieldByName(metadataViewTopShot, "playID").(cadence.UInt32))
		assert.Equal(t, cadence.UInt32(expectedSetID), cadence.SearchFieldByName(metadataViewTopShot, "setID").(cadence.UInt32))
	})

	// Admin sends a transaction that locks the set
	t.Run("Should be able to lock a set which stops plays from being added", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateLockSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// This should fail because the set is locked
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(4))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			true,
		)

		// Script should return that the set is locked
		result := executeScriptAndCheck(t, b, templates.GenerateGetIsSetLockedScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewBool(true), result)
	})

	// Admin sends a transaction that mints a batch of moments
	t.Run("Should be able to mint a batch of moments", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Ensure that the correct number of moments have been minted for each edition
		result := executeScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(3))})
		assert.Equal(t, cadence.NewUInt32(5), result)

		result = executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCollectionIDsScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr))})
		CadenceIntArrayContains(t, result, 1, 2, 3, 4, 5, 6)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

	})

	t.Run("Should be able to mint a batch of moments and fulfill a pack", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateFulfillPackScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		ids := []cadence.Value{cadence.NewUInt64(6), cadence.NewUInt64(7), cadence.NewUInt64(8)}
		_ = tx.AddArgument(cadence.NewArray(ids))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

	})

	// Admin sends a transaction to retire a play
	t.Run("Should be able to retire a Play which stops minting", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRetirePlayScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Minting from this play should fail becuase it is retired
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			true,
		)

		// Make sure this edition is retired
		result := executeScriptAndCheck(t, b, templates.GenerateGetIsEditionRetiredScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewBool(true), result)
	})

	// Admin sends a transaction that retires all the plays in a set
	t.Run("Should be able to retire all Plays which stops minting", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRetireAllPlaysScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// minting should fail
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			true,
		)

		verifyQuerySetMetadata(t, b, env,
			SetMetadata{
				setID:  1,
				name:   "Genesis",
				series: 0,
				plays:  []uint32{1, 2, 3},
				//retired {UInt32: Bool}
				locked: true,
				//numberMintedPerPlay {UInt32: UInt32}})
			})
	})

	// create a new Collection for a user address
	t.Run("Should be able to create a moment Collection", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)
	})

	// Admin sends a transaction to transfer a moment to a user
	t.Run("Should be able to transfer a moment", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentfromShardedCollectionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
		// make sure the user received it
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)
	})

	// Admin sends a transaction to transfer a batch of moments to a user
	t.Run("Should be able to batch transfer moments from a sharded collection", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchTransferMomentfromShardedCollectionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))

		ids := []cadence.Value{cadence.NewUInt64(2), cadence.NewUInt64(3), cadence.NewUInt64(4)}
		_ = tx.AddArgument(cadence.NewArray(ids))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
		// make sure the user received them
		result := executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assert.Equal(t, cadence.NewUInt32(1), result)
	})

	// Admin sends a transaction to update the current series
	t.Run("Should be able to change the current series", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeSeriesScript(env), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	// Make sure the contract fields are correct
	result := executeScriptAndCheck(t, b, templates.GenerateGetSeriesScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(1), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetNextPlayIDScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(5), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetNextSetIDScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(2), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetSupplyScript(env), nil)
	assert.Equal(t, cadence.NewUInt64(11), result)

}

func (b *topshotTestBlockchain) CreatePlay(t *testing.T, metadata []cadence.KeyValuePair) {
	tx := createTxWithTemplateAndAuthorizer(b.Blockchain, templates.GenerateMintPlayScript(b.env), b.topshotAdminAddr)
	play := cadence.NewDictionary(metadata)
	_ = tx.AddArgument(play)

	signAndSubmit(
		t, b.Blockchain, tx,
		[]flow.Address{b.ServiceKey().Address, b.topshotAdminAddr}, []crypto.Signer{b.serviceKeySigner, b.topshotAdminSigner},
		false,
	)
}

func (b *topshotTestBlockchain) UpdateTagline(t *testing.T, plays []cadence.KeyValuePair) {
	tx := createTxWithTemplateAndAuthorizer(b.Blockchain, templates.GenerateUpdateTaglineScript(b.env), b.topshotAdminAddr)
	tags := cadence.NewDictionary(plays)
	_ = tx.AddArgument(tags)

	signAndSubmit(
		t, b.Blockchain, tx,
		[]flow.Address{b.ServiceKey().Address, b.topshotAdminAddr}, []crypto.Signer{b.serviceKeySigner, b.topshotAdminSigner},
		false,
	)
}

func (b *topshotTestBlockchain) CreateSet(t *testing.T, setName string) {
	tx := createTxWithTemplateAndAuthorizer(b.Blockchain, templates.GenerateMintSetScript(b.env), b.topshotAdminAddr)

	_ = tx.AddArgument(CadenceString(setName))

	signAndSubmit(
		t, b.Blockchain, tx,
		[]flow.Address{b.ServiceKey().Address, b.topshotAdminAddr}, []crypto.Signer{b.serviceKeySigner, b.topshotAdminSigner},
		false,
	)
}

// MintMoment will mint a moment from the supplied set and play id into the recipients account
func (b *topshotTestBlockchain) MintMoment(t *testing.T, setID uint32, playID uint32, recipient flow.Address) {
	tx := createTxWithTemplateAndAuthorizer(b.Blockchain, templates.GenerateMintMomentScript(b.env), b.topshotAdminAddr)

	_ = tx.AddArgument(cadence.NewUInt32(setID))
	_ = tx.AddArgument(cadence.NewUInt32(playID))
	_ = tx.AddArgument(cadence.NewAddress(recipient))

	signAndSubmit(
		t, b.Blockchain, tx,
		[]flow.Address{b.ServiceKey().Address, b.topshotAdminAddr}, []crypto.Signer{b.serviceKeySigner, b.topshotAdminSigner},
		false,
	)
}

// This test is for ensuring that admin receiver smart contract works correctly
func TestTransferAdmin(t *testing.T) {
	b := newBlockchain()

	serviceKeySigner, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	accountKeys := test.AccountKeyGenerator()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	// Should be able to deploy a contract as a new account with no keys.
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	viewResolverAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	nftAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.NFTAddress = nftAddr.String()

	metadataViewsAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.MetadataViewsAddress = metadataViewsAddr.String()

	fungibleTokenMetadataViewsAddr := flow.HexToAddress("ee82856bf20e2aa6")

	env.FungibleTokenMetadataViewsAddress = fungibleTokenMetadataViewsAddr.String()

	evmAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.EVMAddress = evmAddr.String()

	// Deploy CrossVMMetadataViews contract
	crossVMMetadataViewsKey, _ := test.AccountKeyGenerator().NewWithSigner()
	crossVMMetadataViewsCode := contracts.GenerateCrossVMMetadataViewsContract(evmAddr.String(), viewResolverAddr.String())
	crossVMMetadataViewsAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{crossVMMetadataViewsKey}, []sdktemplates.Contract{
		{
			Name:   "CrossVMMetadataViews",
			Source: string(crossVMMetadataViewsCode),
		},
	})
	assert.Nil(t, err)
	env.CrossVMMetadataViewsAddress = crossVMMetadataViewsAddr.String()

	// Deploy TopShot Locking contract
	topshotLockingCode := contracts.GenerateTopShotLockingContract(nftAddr.String())
	topShotLockingAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TopShotLocking",
			Source: string(topshotLockingCode),
		},
	})

	royaltyAddr := flow.HexToAddress("ee82856bf20e2aa6")
	assert.Nil(t, err)
	env.FTSwitchboardAddress = royaltyAddr.String()

	tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupSwitchboardScript(env), royaltyAddr)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, royaltyAddr},
		[]crypto.Signer{serviceKeySigner, serviceKeySigner},
		false,
	)

	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.TopShotLockingAddress = topShotLockingAddr.String()

	// First, deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(
		emulatorFTAddress,
		nftAddr.String(),
		metadataViewsAddr.String(),
		viewResolverAddr.String(),
		crossVMMetadataViewsAddr.String(),
		evmAddr.String(),
		topShotLockingAddr.String(),
		royaltyAddr.String(),
		Network,
	)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})

	env.TopShotAddress = topshotAddr.String()

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String(), viewResolverAddr.String())
	shardedAddr, _ := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TopShotShardedCollection",
			Source: string(shardedCollectionCode),
		},
	})
	_, _ = b.CommitBlock()

	env.ShardedAddress = shardedAddr.String()

	// Should be able to deploy the admin receiver contract
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(topshotAddr.String(), shardedAddr.String())
	adminAccountKey, adminSigner := accountKeys.NewWithSigner()
	adminAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{adminAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopshotAdminReceiver",
			Source: string(adminReceiverCode),
		},
	})

	assert.NoError(t, err)
	b.CommitBlock()

	env.AdminReceiverAddress = adminAddr.String()

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")

	// create a new Collection
	t.Run("Should be able to transfer an admin Capability to the receiver account", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferAdminScript(env), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	tb := topshotTestBlockchain{
		Blockchain:         b,
		env:                env,
		topshotAdminAddr:   adminAddr,
		serviceKeySigner:   serviceKeySigner,
		topshotAdminSigner: adminSigner,
		accountKeys:        accountKeys,
	}
	// can create a new play with the new admin
	t.Run("Should be able to create a new Play with the new Admin account", func(t *testing.T) {
		tb.CreatePlay(t, []cadence.KeyValuePair{{Key: firstName, Value: lebron}})
	})
}

func TestSetPlaysOwnedByAddressScript(t *testing.T) {
	// Setup
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	env := tb.env
	serviceKeySigner := tb.serviceKeySigner
	topshotAddr := tb.topshotAdminAddr
	topshotSigner := tb.topshotAdminSigner

	// Create a new user account
	joshAccountKey, joshSigner := tb.accountKeys.NewWithSigner()
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	joshAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{joshAccountKey}, nil)

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
	metadatas := [][]cadence.KeyValuePair{
		{{Key: firstName, Value: lebron}},
		{{Key: firstName, Value: hayward}},
		{{Key: firstName, Value: antetokounmpo}},
	}
	for _, metadata := range metadatas {
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

	// Mint one moment to topshotAddress
	tb.MintMoment(t, genesisSetID, lebronPlayID, topshotAddr)

	t.Run("Should return true if the address owns moments corresponding to each SetPlay", func(t *testing.T) {
		script := templates.GenerateSetPlaysOwnedByAddressScript(env)

		setIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(genesisSetID), cadence.NewUInt32(genesisSetID)})
		playIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(lebronPlayID), cadence.NewUInt32(haywardPlayID)})

		result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(setIDs), jsoncdc.MustEncode(playIDs)})

		assert.Equal(t, cadence.NewBool(true), result)
	})

	t.Run("Should return false if the address does not own moments corresponding to each SetPlay", func(t *testing.T) {
		script := templates.GenerateSetPlaysOwnedByAddressScript(env)

		setIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(genesisSetID), cadence.NewUInt32(genesisSetID), cadence.NewUInt32(genesisSetID)})
		playIDs := cadence.NewArray([]cadence.Value{cadence.NewUInt32(lebronPlayID), cadence.NewUInt32(haywardPlayID), cadence.NewUInt32(antetokounmpoPlayID)})

		result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(setIDs), jsoncdc.MustEncode(playIDs)})
		assert.Equal(t, cadence.NewBool(false), result)
	})

	// t.Run("Should fail with mismatched Set and Play slice lengths", func(t *testing.T) {
	// 	_, err := templates.GenerateSetPlaysOwnedByAddressScript(topshotAddr, joshAddress, []uint32{1, 2}, []uint32{1})
	// 	assert.Error(t, err)
	// 	assert.True(t, strings.Contains(err.Error(), "mismatched lengths"))
	// })

	// t.Run("Should fail with empty SetPlays", func(t *testing.T) {
	// 	_, err := templates.GenerateSetPlaysOwnedByAddressScript(topshotAddr, joshAddress, []uint32{}, []uint32{})
	// 	assert.Error(t, err)
	// 	assert.True(t, strings.Contains(err.Error(), "no SetPlays"))
	// })
}

func TestDestroyMoments(t *testing.T) {
	// Setup
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	serviceKeySigner := tb.serviceKeySigner
	topshotAddr := tb.topshotAdminAddr
	accountKeys := tb.accountKeys
	topshotSigner := tb.topshotAdminSigner
	env := tb.env

	// Should be able to deploy the token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, env.MetadataViewsAddress, env.FungibleTokenMetadataViewsAddress, "DapperUtilityCoin", "dapperUtilityCoin", "1000.0")
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	tokenAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{tokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "DapperUtilityCoin",
			Source: string(tokenCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.DUCAddress = tokenAddr.String()

	// Setup with the first market contract
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), env.DUCAddress)
	marketAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "Market",
			Source: string(marketCode),
		},
	})
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	env.TopShotMarketAddress = marketAddr.String()

	// Should be able to deploy the third market contract
	marketV3Code := contracts.GenerateTopShotMarketV3Contract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), marketAddr.String(), env.DUCAddress, env.TopShotLockingAddress, env.MetadataViewsAddress)
	marketV3Addr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TopShotMarketV3",
			Source: string(marketV3Code),
		},
	})
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	env.TopShotMarketV3Address = marketV3Addr.String()

	// Should be able to deploy the token forwarding contract
	forwardingCode := fungibleToken.CustomTokenForwarding(defaultfungibleTokenAddr, "DapperUtilityCoin", "dapperUtilityCoin")
	forwardingAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TokenForwarding",
			Source: string(forwardingCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.ForwardingAddress = forwardingAddr.String()

	// Create a new user account
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{joshAccountKey}, nil)

	tokenEnv := fungibleTokenTemplates.Environment{
		FungibleTokenAddress:              defaultfungibleTokenAddr,
		MetadataViewsAddress:              tb.metadataViewsAddr.Hex(),
		ExampleTokenAddress:               tokenAddr.Hex(),
		FungibleTokenMetadataViewsAddress: tb.fungibleTokenMetadataViewsAddr.Hex(),
		ViewResolverAddress:               tb.viewResolverAddr.Hex(),
		TokenForwardingAddress:            forwardingAddr.Hex(),
	}

	createTokenScript := fungibleTokenTemplates.GenerateCreateTokenScript(tokenEnv)
	modifiedCreateTokenScript := strings.Replace(string(createTokenScript), "ExampleToken", "DapperUtilityCoin", -1)

	tx := createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateTokenScript), joshAddress)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	mintTokenScript := fungibleTokenTemplates.GenerateMintTokensScript(tokenEnv)
	modifiedMintTokenScript := strings.Replace(string(mintTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
	tx = createTxWithTemplateAndAuthorizer(b, []byte(modifiedMintTokenScript), tokenAddr)
	_ = tx.AddArgument(cadence.NewAddress(joshAddress))
	_ = tx.AddArgument(CadenceUFix64("80.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
		false,
	)

	// Create moment collection
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)
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

	ducPublicPath := cadence.Path{Domain: common.PathDomainPublic, Identifier: "dapperUtilityCoinReceiver"}

	// Create a marketv1 sale collection for Josh
	// setting himself as the beneficiary with a 15% cut
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), joshAddress)

	_ = tx.AddArgument(ducPublicPath)
	_ = tx.AddArgument(cadence.NewAddress(joshAddress))
	_ = tx.AddArgument(CadenceUFix64("0.15"))
	_ = tx.AddArgument(cadence.NewUInt64(2))
	_ = tx.AddArgument(CadenceUFix64("50.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	// Create a sale collection for josh's account, setting josh as the beneficiary
	// and with a 15% cut
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleV3Script(env), joshAddress)

	_ = tx.AddArgument(ducPublicPath)
	_ = tx.AddArgument(cadence.NewAddress(joshAddress))
	_ = tx.AddArgument(CadenceUFix64("0.15"))
	_ = tx.AddArgument(cadence.NewUInt64(1))
	_ = tx.AddArgument(CadenceUFix64("50.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	momentIDs := []uint64{1, 2}
	// check the price, sale length, and the sale's data
	result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})
	assertEqual(t, cadence.NewInt(2), result)
	for _, momentID := range momentIDs {
		result = executeScriptAndCheck(t, b, templates.GenerateGetSalePriceV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(momentID))})
		assertEqual(t, CadenceUFix64("50.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleSetIDV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(momentID))})
		assertEqual(t, cadence.NewUInt32(1), result)
	}

	t.Run("Should destroy the 2 moments created in Josh account", func(t *testing.T) {
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateDestroyMomentsScript(env), joshAddress)
		_ = tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewUInt64(1), cadence.NewUInt64(2)}))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		// verify that the moments no longer exist in josh's collection
		for _, momentID := range momentIDs {
			r, err := b.ExecuteScript(templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(momentID))})
			assert.NoError(t, err)
			assertEqual(t, cadence.NewBool(false), r.Value)
		}
		// verify no moments in sale collection
		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})
		assertEqual(t, cadence.NewInt(0), result)
	})
}

func TestDestroyMomentsV2(t *testing.T) {
	// Setup
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	serviceKeySigner := tb.serviceKeySigner
	topshotAddr := tb.topshotAdminAddr
	accountKeys := tb.accountKeys
	topshotSigner := tb.topshotAdminSigner
	env := tb.env

	// Should be able to deploy the token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, env.MetadataViewsAddress, env.FungibleTokenMetadataViewsAddress, "DapperUtilityCoin", "dapperUtilityCoin", "1000.0")
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	tokenAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{tokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "DapperUtilityCoin",
			Source: string(tokenCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.DUCAddress = tokenAddr.String()

	// Setup with the first market contract
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), env.DUCAddress)
	marketAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "Market",
			Source: string(marketCode),
		},
	})
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	env.TopShotMarketAddress = marketAddr.String()

	// Should be able to deploy the third market contract
	marketV3Code := contracts.GenerateTopShotMarketV3Contract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), marketAddr.String(), env.DUCAddress, env.TopShotLockingAddress, env.MetadataViewsAddress)
	marketV3Addr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TopShotMarketV3",
			Source: string(marketV3Code),
		},
	})
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	env.TopShotMarketV3Address = marketV3Addr.String()

	// Should be able to deploy the token forwarding contract
	forwardingCode := fungibleToken.CustomTokenForwarding(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage)
	forwardingAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "TokenForwarding",
			Source: string(forwardingCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.ForwardingAddress = forwardingAddr.String()

	// Create a new user account
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{joshAccountKey}, nil)

	tokenEnv := fungibleTokenTemplates.Environment{
		FungibleTokenAddress:              defaultfungibleTokenAddr,
		MetadataViewsAddress:              tb.metadataViewsAddr.Hex(),
		ExampleTokenAddress:               tokenAddr.Hex(),
		FungibleTokenMetadataViewsAddress: tb.fungibleTokenMetadataViewsAddr.Hex(),
		ViewResolverAddress:               tb.viewResolverAddr.Hex(),
		TokenForwardingAddress:            forwardingAddr.Hex(),
	}

	createTokenScript := fungibleTokenTemplates.GenerateCreateTokenScript(tokenEnv)
	modifiedCreateTokenScript := strings.Replace(string(createTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
	tx := createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateTokenScript), joshAddress)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	mintTokenScript := fungibleTokenTemplates.GenerateMintTokensScript(tokenEnv)
	modifiedMintTokenScript := strings.Replace(string(mintTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
	tx = createTxWithTemplateAndAuthorizer(b, []byte(modifiedMintTokenScript), tokenAddr)
	_ = tx.AddArgument(cadence.NewAddress(joshAddress))
	_ = tx.AddArgument(CadenceUFix64("80.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
		false,
	)

	// Create moment collection
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)
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

	// Mint three moments to joshAddress
	tb.MintMoment(t, genesisSetID, lebronPlayID, joshAddress)
	tb.MintMoment(t, genesisSetID, haywardPlayID, joshAddress)
	tb.MintMoment(t, genesisSetID, haywardPlayID, joshAddress)

	ducPublicPath := cadence.Path{Domain: common.PathDomainPublic, Identifier: "dapperUtilityCoinReceiver"}

	// Create a marketv1 sale collection for Josh
	// setting himself as the beneficiary with a 15% cut
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), joshAddress)

	_ = tx.AddArgument(ducPublicPath)
	_ = tx.AddArgument(cadence.NewAddress(joshAddress))
	_ = tx.AddArgument(CadenceUFix64("0.15"))
	_ = tx.AddArgument(cadence.NewUInt64(2))
	_ = tx.AddArgument(CadenceUFix64("50.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	// Create a sale collection for josh's account, setting josh as the beneficiary
	// and with a 15% cut
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleV3Script(env), joshAddress)

	_ = tx.AddArgument(ducPublicPath)
	_ = tx.AddArgument(cadence.NewAddress(joshAddress))
	_ = tx.AddArgument(CadenceUFix64("0.15"))
	_ = tx.AddArgument(cadence.NewUInt64(1))
	_ = tx.AddArgument(CadenceUFix64("50.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	// lock the third moment
	tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTopShotLockingLockMomentScript(env), joshAddress)

	_ = tx.AddArgument(cadence.NewUInt64(3))
	duration, _ := cadence.NewUFix64("31536000.0")
	_ = tx.AddArgument(duration)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	// Verify moment is locked
	result := executeScriptAndCheck(t, b, templates.GenerateGetMomentIsLockedScript(env), [][]byte{
		jsoncdc.MustEncode(cadence.Address(joshAddress)),
		jsoncdc.MustEncode(cadence.UInt64(3)),
	})
	assertEqual(t, cadence.NewBool(true), result)

	t.Run("Should destroy the 3 moments created in Josh account, including locked", func(t *testing.T) {
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateDestroyMomentsV2Script(env), joshAddress)
		_ = tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewUInt64(1), cadence.NewUInt64(2), cadence.NewUInt64(3)}))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		// verify that the moments no longer exist in josh's collection
		for _, momentID := range []uint64{1, 2, 3} {
			r, err := b.ExecuteScript(templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(momentID))})
			assert.NoError(t, err)
			assertEqual(t, cadence.NewBool(false), r.Value)
		}
	})
}

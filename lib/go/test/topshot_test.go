package test

import (
	"strings"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/lib/go/templates"

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
	NonFungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	NonFungibleTokenInterfaceFile    = "NonFungibleToken.cdc"

	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"

	emulatorFTAddress         = "ee82856bf20e2aa6"
	emulatorFlowTokenAddress  = "0ae53cb6e3f42a79"
	MetadataFTReplaceAddress  = `"./utility/FungibleToken.cdc"`
	MetadataNFTReplaceAddress = `"./NonFungibleToken.cdc"`
	Network                   = `"mainnet"`
)

// topshotTestContext will expose sugar for common actions needed to bootstrap testing
type topshotTestBlockchain struct {
	*emulator.Blockchain
	env                templates.Environment
	topshotAdminAddr   flow.Address
	serviceKeySigner   crypto.Signer
	topshotAdminSigner crypto.Signer
	accountKeys        *test.AccountKeys

	userAddress flow.Address
}

func NewTopShotTestBlockchain(t *testing.T) topshotTestBlockchain {
	b := newBlockchain()

	serviceKeySigner, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	accountKeys := test.AccountKeyGenerator()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})

	env.NFTAddress = nftAddr.String()

	// Should be able to deploy a contract as a new account with no keys.
	metadataViewsCode, _ := DownloadFile(MetadataViewsContractsBaseURL + MetadataViewsInterfaceFile)
	parsedMetadataContract := strings.Replace(string(metadataViewsCode), MetadataFTReplaceAddress, "0x"+emulatorFTAddress, 1)
	parsedMetadataContract = strings.Replace(parsedMetadataContract, MetadataNFTReplaceAddress, "0x"+nftAddr.String(), 1)
	metadataViewsAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "MetadataViews",
			Source: parsedMetadataContract,
		},
	})
	// errString := err.Error()
	// log.Print(errString)
	assert.Nil(t, err)

	env.MetadataViewsAddress = metadataViewsAddr.String()

	// Deploy TopShot Locking contract
	topshotLockingCode := contracts.GenerateTopShotLockingContract(nftAddr.String())
	topShotLockingAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "TopShotLocking",
			Source: string(topshotLockingCode),
		},
	})

	royaltyAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.TopShotLockingAddress = topShotLockingAddr.String()

	// Deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(emulatorFTAddress, nftAddr.String(), metadataViewsAddr.String(), topShotLockingAddr.String(), royaltyAddr.String(), Network)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})

	env.TopShotAddress = topshotAddr.String()

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
	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String())
	shardedAddr, _ := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "TopShotShardedCollection",
			Source: string(shardedCollectionCode),
		},
	})
	_, _ = b.CommitBlock()

	env.ShardedAddress = shardedAddr.String()

	return topshotTestBlockchain{
		Blockchain:         b,
		env:                env,
		topshotAdminAddr:   topshotAddr,
		serviceKeySigner:   serviceKeySigner,
		topshotAdminSigner: topshotSigner,
		accountKeys:        accountKeys}
}

// This test is for testing the deployment the topshot smart contracts
func TestNFTDeployment(t *testing.T) {
	// all core contracts should be deployed on the test blockchain
	b := NewTopShotTestBlockchain(t)
	env := b.env
	// Should be able to deploy the admin receiver contract
	// as a new account with no keys.
	adminReceiverCode := contracts.GenerateTopshotAdminReceiverContract(env.TopShotAddress, env.ShardedAddress)
	_, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "TopshotAdminReceiver",
			Source: string(adminReceiverCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)
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
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

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
		expectedPrivatePath := "/private/MomentCollection"
		expectedCollectionName := "NBA-Top-Shot"
		expectedCollectionDescription := "NBA Top Shot is your chance to own, sell, and trade official digital collectibles of the NBA and WNBA's greatest plays and players"
		expectedCollectionSquareImage := "https://nbatopshot.com/static/img/og/og.png"
		expectedCollectionBannerImage := "https://nbatopshot.com/static/img/top-shot-logo-horizontal-white.svg"
		expectedRoyaltyReceiversCount := 1
		expectedTraitsCount := 6
		expectedVideoURL := "https://assets.nbatopshot.com/media/1/video"

		resultNFT := executeScriptAndCheck(t, b, templates.GenerateGetNFTMetadataScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		metadataViewNFT := resultNFT.(cadence.Struct)
		assert.Equal(t, cadence.String(expectedMetadataName), metadataViewNFT.Fields[0])
		assert.Equal(t, cadence.String(expectedMetadataDescription), metadataViewNFT.Fields[1])
		assert.Equal(t, cadence.String(expectedMetadataThumbnail), metadataViewNFT.Fields[2])
		assert.Equal(t, cadence.String(expectedMetadataExternalURL), metadataViewNFT.Fields[5])
		assert.Equal(t, cadence.String(expectedStoragePath), metadataViewNFT.Fields[6])
		assert.Equal(t, cadence.String(expectedPublicPath), metadataViewNFT.Fields[7])
		assert.Equal(t, cadence.String(expectedPrivatePath), metadataViewNFT.Fields[8])
		assert.Equal(t, cadence.String(expectedCollectionName), metadataViewNFT.Fields[9])
		assert.Equal(t, cadence.String(expectedCollectionDescription), metadataViewNFT.Fields[10])
		assert.Equal(t, cadence.String(expectedCollectionSquareImage), metadataViewNFT.Fields[11])
		assert.Equal(t, cadence.String(expectedCollectionBannerImage), metadataViewNFT.Fields[12])
		assert.Equal(t, cadence.UInt32(expectedRoyaltyReceiversCount), metadataViewNFT.Fields[13])
		assert.Equal(t, cadence.UInt32(expectedTraitsCount), metadataViewNFT.Fields[14])
		assert.Equal(t, cadence.String(expectedVideoURL), metadataViewNFT.Fields[15])

		// Tests that top-shot specific metadata is discoverable on-chain
		expectedPlayID := 1
		expectedSetID := 1
		expectedSerialNumber := 1

		resultTopShot := executeScriptAndCheck(t, b, templates.GenerateGetTopShotMetadataScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		metadataViewTopShot := resultTopShot.(cadence.Struct)
		assert.Equal(t, cadence.UInt32(expectedSerialNumber), metadataViewTopShot.Fields[26])
		assert.Equal(t, cadence.UInt32(expectedPlayID), metadataViewTopShot.Fields[27])
		assert.Equal(t, cadence.UInt32(expectedSetID), metadataViewTopShot.Fields[28])
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
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, _ := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})

	env.NFTAddress = nftAddr.String()

	// Should be able to deploy a contract as a new account with no keys.
	metadataViewsCode, _ := DownloadFile(MetadataViewsContractsBaseURL + MetadataViewsInterfaceFile)
	parsedMetadataContract := strings.Replace(string(metadataViewsCode), MetadataFTReplaceAddress, "0x"+emulatorFTAddress, 1)
	parsedMetadataContract = strings.Replace(parsedMetadataContract, MetadataNFTReplaceAddress, "0x"+nftAddr.String(), 1)
	metadataViewsAddr, _ := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "MetadataViews",
			Source: parsedMetadataContract,
		},
	})

	env.MetadataViewsAddress = metadataViewsAddr.String()

	// Deploy TopShot Locking contract
	topshotLockingCode := contracts.GenerateTopShotLockingContract(nftAddr.String())
	topShotLockingAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "TopShotLocking",
			Source: string(topshotLockingCode),
		},
	})

	royaltyAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.TopShotLockingAddress = topShotLockingAddr.String()

	// First, deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(emulatorFTAddress, nftAddr.String(), metadataViewsAddr.String(), topShotLockingAddr.String(), royaltyAddr.String(), Network)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})

	env.TopShotAddress = topshotAddr.String()

	shardedCollectionCode := contracts.GenerateTopShotShardedCollectionContract(nftAddr.String(), topshotAddr.String())
	shardedAddr, _ := b.CreateAccount(nil, []sdktemplates.Contract{
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
	adminAddr, _ := b.CreateAccount([]*flow.AccountKey{adminAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopshotAdminReceiver",
			Source: string(adminReceiverCode),
		},
	})
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
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

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
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage, "1000.0")
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
	tokenAddr, err := b.CreateAccount([]*flow.AccountKey{tokenAccountKey}, []sdktemplates.Contract{
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
	marketAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
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
	marketV3Code := contracts.GenerateTopShotMarketV3Contract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), marketAddr.String(), env.DUCAddress, env.TopShotLockingAddress)
	marketV3Addr, err := b.CreateAccount(nil, []sdktemplates.Contract{
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
	forwardingAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
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
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	tx := createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName), joshAddress)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
		false,
	)

	tx = createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateMintTokensScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, joshAddress, defaultTokenName, 80), tokenAddr)

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

	ducPublicPath := cadence.Path{Domain: "public", Identifier: "dapperUtilityCoinReceiver"}

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
			assert.Contains(t, r.Error.Error(), "NFT does not exist in the collection")
		}
		// verify no moments in sale collection
		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})
		assertEqual(t, cadence.NewInt(0), result)
	})
}

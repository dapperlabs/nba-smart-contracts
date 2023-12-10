package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

// This test tests the pure functionality of the smart contract
func TestSubeditions(t *testing.T) {
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

	playIDString := CadenceString("PlayID")
	setIDString := CadenceString("SetID")
	value1 := CadenceString("1")
	value3 := CadenceString("3")
	subedition111Name := CadenceString("Subedition PlayID:1 SetID:1 SubeditionID: 1")
	subedition112Name := CadenceString("Subedition PlayID:1 SetID:1 SubeditionID: 2")
	subedition133Name := CadenceString("Subedition PlayID:3 SetID:1 SubeditionID: 3")
	subedition134Name := CadenceString("Subedition PlayID:3 SetID:1 SubeditionID: 4")
	var result cadence.Value
	// Admin sends a transaction to create a play
	t.Run("Should be able to create a new Play", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata := []cadence.KeyValuePair{{Key: firstName, Value: lebron}, {Key: playType, Value: dunk}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	// Admin sends transactions to create multiple plays
	t.Run("Should be able to create multiple new Plays", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata := []cadence.KeyValuePair{{Key: firstName, Value: oladipo}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata = []cadence.KeyValuePair{{Key: firstName, Value: hayward}}
		play = cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata = []cadence.KeyValuePair{{Key: firstName, Value: durant}}
		play = cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Check that the return all plays script doesn't fail
		// and that we can return metadata about the plays
		executeScriptAndCheck(t, b, templates.GenerateGetAllPlaysScript(env), nil)

		result = executeScriptAndCheck(t, b, templates.GenerateGetPlayMetadataFieldScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.String("FullName"))})
		assert.Equal(t, CadenceString("Lebron"), result)
	})

	// Admin creates a new Set with the name Genesis
	t.Run("Should be able to create a new Set", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(env), topshotAddr)

		_ = tx.AddArgument(CadenceString("Genesis"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Check that the set name, ID, and series were initialized correctly.
		result := executeScriptAndCheck(t, b, templates.GenerateGetSetNameScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, CadenceString("Genesis"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSetIDsByNameScript(env), [][]byte{jsoncdc.MustEncode(cadence.String("Genesis"))})
		assert.Equal(t, UInt32Array(1), result)

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
		result = executeScriptAndCheck(t, b, templates.GenerateGetNextPlayIDScript(env), nil)
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

	t.Run("Should be able to create new subedition admin resource", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateNewSubeditionAdminResourceScript(env), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	t.Run("Should be able to create multiple new Subeditions", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSubeditionScript(env), topshotAddr)

		name := subedition111Name
		metadata := []cadence.KeyValuePair{{Key: setIDString, Value: value1}, {Key: playIDString, Value: value1}}
		subeditionMetadata := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(name)
		_ = tx.AddArgument(subeditionMetadata)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSubeditionScript(env), topshotAddr)

		name = subedition112Name
		metadata = []cadence.KeyValuePair{{Key: setIDString, Value: value1}, {Key: playIDString, Value: value1}}
		subeditionMetadata = cadence.NewDictionary(metadata)
		_ = tx.AddArgument(name)
		_ = tx.AddArgument(subeditionMetadata)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSubeditionScript(env), topshotAddr)

		name = subedition133Name
		metadata = []cadence.KeyValuePair{{Key: setIDString, Value: value1}, {Key: playIDString, Value: value3}}
		subeditionMetadata = cadence.NewDictionary(metadata)
		_ = tx.AddArgument(name)
		_ = tx.AddArgument(subeditionMetadata)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSubeditionScript(env), topshotAddr)

		name = subedition134Name
		metadata = []cadence.KeyValuePair{{Key: setIDString, Value: value1}, {Key: playIDString, Value: value3}}
		subeditionMetadata = cadence.NewDictionary(metadata)
		_ = tx.AddArgument(name)
		_ = tx.AddArgument(subeditionMetadata)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNextSubeditionIDScript(env), nil)
		assert.Equal(t, cadence.NewUInt32(5), result)

		// Check that the return all subeditions script works properly
		// and that we can return metadata about the plays
		result = executeScriptAndCheck(t, b, templates.GenerateGetAllSubeditionScript(env), nil)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSubeditionByIDScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1))})

		metadataFields := result.(cadence.Struct).Fields

		metadata = []cadence.KeyValuePair{{Key: playIDString, Value: value1}, {Key: setIDString, Value: value1}}
		subeditionMetadata = CadenceStringDictionary(metadata)
		assert.Equal(t, cadence.NewUInt32(1), metadataFields[0])
		assert.Equal(t, subedition111Name, metadataFields[1])
		assert.Equal(t, subeditionMetadata, metadataFields[2])
	})

	t.Run("Should be able to link nft to subedition", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetNFTsubedtitionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt64(100))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNFTSubeditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt64(100))})
		assert.Equal(t, cadence.NewUInt32(1), result)
	})

	// Admin mints a moment that stores it in the admin's collection
	t.Run("Should be able to mint a moment with subedition #1", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentWithSubeditionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// make sure the moment was minted correctly and is stored in the collection with the correct data
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCollectionIDsScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr))})
		CadenceIntArrayContains(t, result, 1)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

	})

	t.Run("Should be able to mint a moment with subedition #2", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentWithSubeditionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(2))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// make sure the moment was minted correctly and is stored in the collection with the correct data
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCollectionIDsScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr))})
		CadenceIntArrayContains(t, result, 1, 2)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

	})

	t.Run("Should be able to get moments metadata", func(t *testing.T) {
		// Tests to ensure that all core metadataviews are resolvable
		expectedMetadataName := "Lebron Dunk"
		expectedMetadataDescription := "A series 0 Genesis moment with serial number 1"
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
	t.Run("Should be able to mint a batch of moments with subedition #1", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentWithSubeditionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Ensure that the correct number of moments have been minted for each edition
		result := executeScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewUInt32(2), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(3))})
		assert.Equal(t, cadence.NewUInt32(5), result)

		result = executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCollectionIDsScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr))})
		CadenceIntArrayContains(t, result, 3, 5, 4, 2, 6, 7, 1)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

	})

	t.Run("Should be able to mint a batch of moments with subedition #2", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentWithSubeditionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(cadence.NewUInt32(4))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Ensure that the correct number of moments have been minted for each edition
		result := executeScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))})
		assert.Equal(t, cadence.NewUInt32(2), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNumMomentsInEditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(3))})
		assert.Equal(t, cadence.NewUInt32(10), result)

		result = executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCollectionIDsScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr))})
		CadenceIntArrayContains(t, result, 3, 8, 9, 10, 12, 4, 2, 6, 11, 7, 1, 5)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

	})

	t.Run("Should be able to get moment's subedition", func(t *testing.T) {
		//check separately minted moments
		result = executeScriptAndCheck(t, b, templates.GenerateGetNFTSubeditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewUInt32(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNFTSubeditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt64(2))})
		assert.Equal(t, cadence.NewUInt32(2), result)

		//check batch minted moments
		result = executeScriptAndCheck(t, b, templates.GenerateGetNFTSubeditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt64(3))})
		assert.Equal(t, cadence.NewUInt32(3), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNFTSubeditionScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt64(8))})
		assert.Equal(t, cadence.NewUInt32(4), result)
	})

	t.Run("Should be able to check moment's serial number", func(t *testing.T) {
		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(3))})
		assert.Equal(t, cadence.NewUInt32(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(4))})
		assert.Equal(t, cadence.NewUInt32(2), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(5))})
		assert.Equal(t, cadence.NewUInt32(3), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(6))})
		assert.Equal(t, cadence.NewUInt32(4), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(7))})
		assert.Equal(t, cadence.NewUInt32(5), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(8))})
		assert.Equal(t, cadence.NewUInt32(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(9))})
		assert.Equal(t, cadence.NewUInt32(2), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(10))})
		assert.Equal(t, cadence.NewUInt32(3), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(11))})
		assert.Equal(t, cadence.NewUInt32(4), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSerialNumScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(topshotAddr)), jsoncdc.MustEncode(cadence.UInt64(12))})
		assert.Equal(t, cadence.NewUInt32(5), result)
	})

	t.Run("Should be able to mint a batch of moments with subedition and fulfill a pack", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentWithSubeditionScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(3))
		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(cadence.NewUInt32(1))
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
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentWithSubeditionScript(env), topshotAddr)

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
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentWithSubeditionScript(env), topshotAddr)

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
		result = executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
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
		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentSetScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
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
	result = executeScriptAndCheck(t, b, templates.GenerateGetSeriesScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(1), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetNextPlayIDScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(5), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetNextSetIDScript(env), nil)
	assert.Equal(t, cadence.NewUInt32(2), result)

	result = executeScriptAndCheck(t, b, templates.GenerateGetSupplyScript(env), nil)
	assert.Equal(t, cadence.NewUInt64(17), result)

}

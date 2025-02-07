package test

import (
	"context"
	"testing"
	"time"

	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/adapters"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests all the main functionality of the TopShot Locking contract
func TestFastBreak(t *testing.T) {
	//b := newBlockchain()
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	serviceKeySigner, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	accountKeys := test.AccountKeyGenerator()

	env := tb.env
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	viewResolverAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.ViewResolverAddress = viewResolverAddr.String()

	nftAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.NFTAddress = nftAddr.String()

	metadataViewsAddr := flow.HexToAddress("f8d6e0586b0a20c7")
	env.MetadataViewsAddress = metadataViewsAddr.String()

	// Deploy TopShot Locking contract
	lockingKey, lockingSigner := test.AccountKeyGenerator().NewWithSigner()
	topshotLockingCode := contracts.GenerateTopShotLockingContract(nftAddr.String())
	topShotLockingAddr, topShotLockingAddrErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{lockingKey}, []sdktemplates.Contract{
		{
			Name:   "TopShotLocking",
			Source: string(topshotLockingCode),
		},
	})
	assert.Nil(t, topShotLockingAddrErr)
	env.TopShotLockingAddress = topShotLockingAddr.String()

	topShotRoyaltyAddr := flow.HexToAddress("ee82856bf20e2aa6")
	env.FTSwitchboardAddress = topShotRoyaltyAddr.String()

	// Deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(defaultfungibleTokenAddr, nftAddr.String(), metadataViewsAddr.String(), viewResolverAddr.String(), topShotLockingAddr.String(), topShotRoyaltyAddr.String(), Network)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, topshotAddrErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})
	assert.Nil(t, topshotAddrErr)
	env.TopShotAddress = topshotAddr.String()

	// Update the locking contract with topshot address
	topShotLockingCodeWithRuntimeAddr := contracts.GenerateTopShotLockingContractWithTopShotRuntimeAddr(nftAddr.String(), topshotAddr.String())
	updateErr := updateContract(b, topShotLockingAddr, lockingSigner, "TopShotLocking", topShotLockingCodeWithRuntimeAddr)
	assert.Nil(t, updateErr)

	// Should be able to deploy the token contract
	tokenCode := fungibleToken.CustomToken(
		defaultfungibleTokenAddr,
		env.MetadataViewsAddress,
		env.FungibleTokenMetadataViewsAddress,
		"DapperUtilityCoin",
		"dapperUtilityCoin",
		"1000.0",
	)
	tokenAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "DapperUtilityCoin",
			Source: string(tokenCode),
		},
	})
	assert.NoError(t, err)
	env.DUCAddress = tokenAddr.String()

	// Setup with the first market contract
	marketAccountKey, marketSigner := accountKeys.NewWithSigner()
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), env.DUCAddress)
	marketAddr, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{marketAccountKey}, []sdktemplates.Contract{
		{
			Name:   "Market",
			Source: string(marketCode),
		},
	})
	assert.NoError(t, err)
	env.TopShotMarketAddress = marketAddr.String()

	// Should be able to deploy the third market contract
	marketV3Code := contracts.GenerateTopShotMarketV3Contract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), marketAddr.String(), env.DUCAddress, env.TopShotLockingAddress, env.MetadataViewsAddress)

	tx1 := sdktemplates.AddAccountContract(
		marketAddr,
		sdktemplates.Contract{
			Name:   "TopShotMarketV3",
			Source: string(marketV3Code),
		},
	)
	tx1.
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, marketAddr},
		[]crypto.Signer{signer, marketSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)
	env.TopShotMarketV3Address = marketAddr.String()

	// Deploy Fast Break
	fastBreakKey, fastBreakSigner := test.AccountKeyGenerator().NewWithSigner()
	fastBreakCode := contracts.GenerateFastBreakContract(nftAddr.String(), topshotAddr.String(), metadataViewsAddr.String(), env.TopShotMarketV3Address)
	fastBreakAddr, fastBreakAddrErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{fastBreakKey}, []sdktemplates.Contract{
		{
			Name:   "FastBreakV1",
			Source: string(fastBreakCode),
		},
	})
	require.NoError(t, fastBreakAddrErr)
	env.FastBreakAddress = fastBreakAddr.String()

	// create a new user account
	aliceAccountKey, aliceSigner := accountKeys.NewWithSigner()
	aliceAddress, aliceAddressErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{aliceAccountKey}, nil)
	require.NoError(t, aliceAddressErr)

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")
	playType := CadenceString("PlayType")
	dunk := CadenceString("Dunk")

	// Create Play
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)
		metadata := []cadence.KeyValuePair{{Key: firstName, Value: lebron}, {Key: playType, Value: dunk}}
		play := cadence.NewDictionary(metadata)

		arg0Err := tx.AddArgument(play)
		assert.Nil(t, arg0Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	// Create Set
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(env), topshotAddr)

		arg0Err := tx.AddArgument(CadenceString("Genesis"))
		assert.Nil(t, arg0Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	// Add Play to Set
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		arg0Err := tx.AddArgument(cadence.NewUInt32(1))
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewUInt32(1))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	var (
		//fast break run

		tomorrow         = time.Now().Add(24 * time.Hour)
		nextWeek         = time.Now().Add(7 * 24 * time.Hour)
		lastWeek         = time.Now().Add(-7 * 24 * time.Hour)
		fastBreakRunId   = "abc-123"
		fastBreakRunName = "R0"
		runStart         = lastWeek.Unix()
		runEnd           = nextWeek.Unix()
		fatigueModeOn    = true

		//fast break
		fastBreakID                  = "def-456"
		fastBreakName                = "fb0"
		submissionDeadline           = tomorrow.Unix()
		numPlayers            uint64 = 1
		fastBreakStartedState uint8  = 1
		playerId              uint64

		//fast break stat
		statName           = "POINTS"
		statRawType uint8  = 0
		valueNeeded uint64 = 30
	)

	t.Run("oracle should be able to create a fast break run", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateRunScript(env), fastBreakAddr)
		cdcId, cdcIdErr := cadence.NewString(fastBreakRunId)
		assert.Nil(t, cdcIdErr)

		cdcName, cdcNameErr := cadence.NewString(fastBreakRunName)
		assert.Nil(t, cdcNameErr)

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cdcName)
		assert.Nil(t, arg1Err)

		arg2Err := tx.AddArgument(cadence.NewUInt64(uint64(runStart)))
		assert.Nil(t, arg2Err)

		arg3Err := tx.AddArgument(cadence.NewUInt64(uint64(runEnd)))
		assert.Nil(t, arg3Err)

		arg4Err := tx.AddArgument(cadence.NewBool(fatigueModeOn))
		assert.Nil(t, arg4Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, fastBreakSigner},
			false,
		)

	})

	t.Run("oracle should be able to create a fast break game", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateGameScript(env), fastBreakAddr)
		cdcId, cdcIdErr := cadence.NewString(fastBreakID)
		assert.Nil(t, cdcIdErr)

		cdcName, cdcNameErr := cadence.NewString(fastBreakName)
		assert.Nil(t, cdcNameErr)

		cdcFbrId, cdcFbrIdErr := cadence.NewString(fastBreakRunId)
		assert.Nil(t, cdcFbrIdErr)

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cdcName)
		assert.Nil(t, arg1Err)

		arg2Err := tx.AddArgument(cdcFbrId)
		assert.Nil(t, arg2Err)

		arg3Err := tx.AddArgument(cadence.NewUInt64(uint64(submissionDeadline)))
		assert.Nil(t, arg3Err)

		arg4Err := tx.AddArgument(cadence.NewUInt64(numPlayers))
		assert.Nil(t, arg4Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, fastBreakSigner},
			false,
		)

		// Check that that main contract fields were initialized correctly
		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakScript(env), [][]byte{jsoncdc.MustEncode(cdcId)})
		assert.NotNil(t, result)

		resultId := cadence.SearchFieldByName(result.(cadence.Optional).Value.(cadence.Struct), "id")
		assert.Equal(t, cadence.String(fastBreakID), resultId)
	})

	t.Run("oracle should be able to add a stat to a fast break game", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddStatToGameScript(env), fastBreakAddr)
		cdcId, cdcIdErr := cadence.NewString(fastBreakID)
		assert.Nil(t, cdcIdErr)

		cdcName, cdcNameErr := cadence.NewString(statName)
		assert.Nil(t, cdcNameErr)

		cdcType := cadence.NewUInt8(statRawType)

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cdcName)
		assert.Nil(t, arg1Err)

		arg2Err := tx.AddArgument(cdcType)
		assert.Nil(t, arg2Err)

		arg3Err := tx.AddArgument(cadence.NewUInt64(valueNeeded))
		assert.Nil(t, arg3Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, fastBreakSigner},
			false,
		)

		// Check that that main contract fields were initialized correctly
		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakStatsScript(env), [][]byte{jsoncdc.MustEncode(cdcId)})
		interfaceArray := result.(cadence.Array)
		assert.Len(t, interfaceArray.Values, 1)
		assert.NotNil(t, result)
	})

	t.Run("player should be able to create a moment collection", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), aliceAddress)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, aliceAddress}, []crypto.Signer{serviceKeySigner, aliceSigner},
			false,
		)
	})

	t.Run("player should be able to setup game wallet", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateFastBreakCreateAccountScript(env), aliceAddress)
		playerName, playerNameErr := cadence.NewString("houseofhufflepuff")
		assert.Nil(t, playerNameErr)

		arg0Err := tx.AddArgument(playerName)
		assert.Nil(t, arg0Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, aliceAddress}, []crypto.Signer{serviceKeySigner, aliceSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateCurrentPlayerScript(env), nil)
		require.NotNil(t, result)
		playerId = uint64(result.(cadence.UInt64))
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

	// mint moment 1 to alice
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentScript(env), topshotAddr)

		arg0Err := tx.AddArgument(cadence.NewUInt32(1))
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewUInt32(1))
		assert.Nil(t, arg1Err)

		arg2Err := tx.AddArgument(cadence.NewAddress(aliceAddress))
		assert.Nil(t, arg2Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	t.Run("player should not be able play fast break without top shots", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePlayFastBreakScript(env), aliceAddress)
		cdcId, _ := cadence.NewString(fastBreakID)
		ids := []cadence.Value{cadence.NewUInt64(2)}

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewArray(ids))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, aliceAddress}, []crypto.Signer{serviceKeySigner, aliceSigner},
			true,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakTokenCountScript(env), nil)
		assert.Equal(t, cadence.NewUInt64(0), result)
	})

	t.Run("player should be able to play fast break", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePlayFastBreakScript(env), aliceAddress)
		cdcId, _ := cadence.NewString(fastBreakID)
		cdcTopShots := []cadence.Value{cadence.NewUInt64(1)}

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewArray(cdcTopShots))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, aliceAddress}, []crypto.Signer{serviceKeySigner, aliceSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakTokenCountScript(env), nil)
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

	t.Run("player should not be able to resubmit fast break", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePlayFastBreakScript(env), aliceAddress)
		cdcId, _ := cadence.NewString(fastBreakID)
		ids := []cadence.Value{cadence.NewUInt64(1)}

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewArray(ids))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, aliceAddress}, []crypto.Signer{serviceKeySigner, aliceSigner},
			true,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakTokenCountScript(env), nil)
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

	t.Run("oracle should be able to update status of fast break", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateUpdateFastBreakGameScript(env), fastBreakAddr)
		cdcId, _ := cadence.NewString(fastBreakID)
		cdcState := cadence.NewUInt8(fastBreakStartedState)

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cdcState)
		assert.Nil(t, arg1Err)

		// winner
		arg2Err := tx.AddArgument(cadence.NewUInt64(playerId))
		assert.Nil(t, arg2Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, fastBreakSigner},
			false,
		)
	})

	t.Run("oracle should be to score a submission to fast break", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateScoreFastBreakSubmissionScript(env), fastBreakAddr)

		cdcId, _ := cadence.NewString(fastBreakID)

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewAddress(aliceAddress))
		assert.Nil(t, arg1Err)

		arg2Err := tx.AddArgument(cadence.NewUInt64(100))
		assert.Nil(t, arg2Err)

		arg3Err := tx.AddArgument(cadence.NewBool(true))
		assert.Nil(t, arg3Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, fastBreakSigner},
			false,
		)

		result := executeScriptAndCheck(
			t,
			b,
			templates.GenerateGetPlayerScoreScript(env),
			[][]byte{jsoncdc.MustEncode(cdcId), jsoncdc.MustEncode(cadence.NewAddress(aliceAddress))},
		)
		assert.Equal(t, cadence.NewUInt64(100), result)

		fastBreakRunIdCadence, _ := cadence.NewString(fastBreakRunId)

		result = executeScriptAndCheck(
			t,
			b,
			templates.GenerateGetPlayerWinCountForRunScript(env),
			[][]byte{jsoncdc.MustEncode(fastBreakRunIdCadence), jsoncdc.MustEncode(cadence.NewAddress(aliceAddress))},
		)
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

}

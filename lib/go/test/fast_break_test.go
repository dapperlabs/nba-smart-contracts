package test

import (
	"context"
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
	"strings"
	"testing"
	"time"
)

// Tests all the main functionality of the TopShot Locking contract
func TestFastBreak(t *testing.T) {
	b := newBlockchain()

	serviceKeySigner, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	accountKeys := test.AccountKeyGenerator()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	viewResolverCode, _ := DownloadFile(ViewResolverContractsBaseURL + ViewResolverInterfaceFile)
	viewResolverAddr, err := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "ViewResolver",
			Source: string(viewResolverCode),
		},
	})

	assert.Nil(t, err)
	env.ViewResolverAddress = viewResolverAddr.String()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, nftCodeErr := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	parsedNFTContract := strings.Replace(string(nftCode), ViewResolverReplaceAddress, "0x"+viewResolverAddr.String(), 1)
	assert.Nil(t, nftCodeErr)
	nftAddr, nftAddrErr := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: parsedNFTContract,
		},
	})
	assert.Nil(t, nftAddrErr)
	env.NFTAddress = nftAddr.String()

	// Should be able to deploy a contract as a new account with no keys.
	metadataViewsCode, metadataViewsCodeErr := DownloadFile(MetadataViewsContractsBaseURL + MetadataViewsInterfaceFile)
	assert.Nil(t, metadataViewsCodeErr)
	parsedMetadataContract := strings.Replace(string(metadataViewsCode), MetadataFTReplaceAddress, "0x"+emulatorFTAddress, 1)
	parsedMetadataContract = strings.Replace(parsedMetadataContract, MetadataNFTReplaceAddress, "0x"+nftAddr.String(), 1)
	parsedMetadataContract = strings.Replace(parsedMetadataContract, ViewResolverReplaceAddress, "0x"+viewResolverAddr.String(), 1)
	metadataViewsAddr, metadataViewsAddrErr := adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
		{
			Name:   "MetadataViews",
			Source: parsedMetadataContract,
		},
	})
	assert.Nil(t, metadataViewsAddrErr)
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

	fungibleTokenSwitchboardCode, _ := DownloadFile(FungibleTokenSwitchboardContractsBaseURL + FungibleTokenSwitchboardInterfaceFile)
	parsedFungibleSwitchboardContract := strings.Replace(string(fungibleTokenSwitchboardCode), MetadataFTReplaceAddress, "0x"+emulatorFTAddress, 1)
	topShotRoyaltyAccountKey, _ := accountKeys.NewWithSigner()
	topShotRoyaltyAddr, topShotRoyaltyAddrErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{topShotRoyaltyAccountKey}, []sdktemplates.Contract{
		{
			Name:   "FungibleTokenSwitchboard",
			Source: parsedFungibleSwitchboardContract,
		},
	})
	assert.Nil(t, err)
	env.FTSwitchboardAddress = topShotRoyaltyAddr.String()
	assert.Nil(t, topShotRoyaltyAddrErr)

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

	// Deploy Fast Break
	fastBreakKey, fastBreakSigner := test.AccountKeyGenerator().NewWithSigner()
	fastBreakCode := contracts.GenerateFastBreakContract(nftAddr.String(), topshotAddr.String())
	fastBreakAddr, fastBreakAddrErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{fastBreakKey}, []sdktemplates.Contract{
		{
			Name:   "FastBreakV1",
			Source: string(fastBreakCode),
		},
	})
	assert.Nil(t, fastBreakAddrErr)
	env.FastBreakAddress = fastBreakAddr.String()
	assert.Nil(t, err)

	// create a new user account
	jerAccountKey, jerSigner := accountKeys.NewWithSigner()
	jerAddress, jerAddressErr := adapter.CreateAccount(context.Background(), []*flow.AccountKey{jerAccountKey}, nil)
	assert.Nil(t, jerAddressErr)

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
		interfaceArray := result.ToGoValue().([]interface{})
		resultId := interfaceArray[0].(string)
		assert.NotNil(t, result)
		assert.Equal(t, fastBreakID, resultId)
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
		interfaceArray := result.ToGoValue().([]interface{})
		assert.Equal(t, 1, len(interfaceArray))
		assert.NotNil(t, result)
	})

	t.Run("player should be able to create a moment collection", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), jerAddress)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
			false,
		)
	})

	t.Run("player should be able to setup game wallet", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateFastBreakCreateAccountScript(env), jerAddress)
		playerName, playerNameErr := cadence.NewString("houseofhufflepuff")
		assert.Nil(t, playerNameErr)

		arg0Err := tx.AddArgument(playerName)
		assert.Nil(t, arg0Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateCurrentPlayerScript(env), nil)
		playerId = result.ToGoValue().(uint64)
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

	// mint moment 1 to jer
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentScript(env), topshotAddr)

		arg0Err := tx.AddArgument(cadence.NewUInt32(1))
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewUInt32(1))
		assert.Nil(t, arg1Err)

		arg2Err := tx.AddArgument(cadence.NewAddress(jerAddress))
		assert.Nil(t, arg2Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	t.Run("player should not be able play fast break without top shots", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePlayFastBreakScript(env), jerAddress)
		cdcId, _ := cadence.NewString(fastBreakID)
		ids := []cadence.Value{cadence.NewUInt64(2)}

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewArray(ids))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
			true,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakTokenCountScript(env), nil)
		assert.Equal(t, cadence.NewUInt64(0), result)
	})

	t.Run("player should be able to play fast break", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePlayFastBreakScript(env), jerAddress)
		cdcId, _ := cadence.NewString(fastBreakID)
		cdcTopShots := []cadence.Value{cadence.NewUInt64(1)}

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewArray(cdcTopShots))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakTokenCountScript(env), nil)
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

	t.Run("player should not be able to resubmit fast break", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePlayFastBreakScript(env), jerAddress)
		cdcId, _ := cadence.NewString(fastBreakID)
		ids := []cadence.Value{cadence.NewUInt64(1)}

		arg0Err := tx.AddArgument(cdcId)
		assert.Nil(t, arg0Err)

		arg1Err := tx.AddArgument(cadence.NewArray(ids))
		assert.Nil(t, arg1Err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
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

		arg1Err := tx.AddArgument(cadence.NewAddress(jerAddress))
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
			[][]byte{jsoncdc.MustEncode(cdcId), jsoncdc.MustEncode(cadence.NewAddress(jerAddress))},
		)
		assert.Equal(t, cadence.NewUInt64(100), result)
	})

}

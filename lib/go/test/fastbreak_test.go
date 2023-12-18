package test

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
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
	env.MetadataViewsAddress = metadataViewsAddr.String()

	// Deploy TopShot Locking contract
	lockingKey, lockingSigner := test.AccountKeyGenerator().NewWithSigner()
	topshotLockingCode := contracts.GenerateTopShotLockingContract(nftAddr.String())
	topShotLockingAddr, err := b.CreateAccount([]*flow.AccountKey{lockingKey}, []sdktemplates.Contract{
		{
			Name:   "TopShotLocking",
			Source: string(topshotLockingCode),
		},
	})
	env.TopShotLockingAddress = topShotLockingAddr.String()

	topShotRoyaltyAddr, err := b.CreateAccount([]*flow.AccountKey{lockingKey}, []sdktemplates.Contract{})

	// Deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(defaultfungibleTokenAddr, nftAddr.String(), metadataViewsAddr.String(), topShotLockingAddr.String(), topShotRoyaltyAddr.String(), Network)
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, _ := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})
	env.TopShotAddress = topshotAddr.String()

	// Update the locking contract with topshot address
	topShotLockingCodeWithRuntimeAddr := contracts.GenerateTopShotLockingContractWithTopShotRuntimeAddr(nftAddr.String(), topshotAddr.String())
	err = updateContract(b, topShotLockingAddr, lockingSigner, "TopShotLocking", topShotLockingCodeWithRuntimeAddr)
	assert.Nil(t, err)

	// Deploy Fast Break
	fastBreakKey, _ := test.AccountKeyGenerator().NewWithSigner()
	fastBreakCode := contracts.GenerateFastBreakContract(nftAddr.String(), topshotAddr.String())
	fastBreakAddr, err := b.CreateAccount([]*flow.AccountKey{fastBreakKey}, []sdktemplates.Contract{
		{
			Name:   "FastBreak",
			Source: string(fastBreakCode),
		},
	})
	env.FastBreakAddress = fastBreakAddr.String()
	fmt.Println(err)
	assert.Nil(t, err)

	// create a new user account
	jerAccountKey, jerSigner := accountKeys.NewWithSigner()
	jerAddress, _ := b.CreateAccount([]*flow.AccountKey{jerAccountKey}, nil)

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")
	playType := CadenceString("PlayType")
	dunk := CadenceString("Dunk")

	// Create Play
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)
		metadata := []cadence.KeyValuePair{{Key: firstName, Value: lebron}, {Key: playType, Value: dunk}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	// Create Set
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(env), topshotAddr)

		_ = tx.AddArgument(CadenceString("Genesis"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}

	// Add Play to Set
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

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
		fastBreakID               = "def-456"
		fastBreakName             = "fb0"
		isPublic                  = true
		submissionDeadline        = tomorrow.Unix()
		numPlayers         uint64 = 1

		//fast break stat
		statName           = "POINTS"
		statType           = "INDIVIDUAL"
		valueNeeded uint64 = 30
	)

	t.Run("oracle should be able to create a fast break run", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateRunScript(env), fastBreakAddr)
		cdcId, _ := cadence.NewString(fastBreakRunId)
		cdcName, _ := cadence.NewString(fastBreakRunName)

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cdcName)
		_ = tx.AddArgument(cadence.NewUInt64(uint64(runStart)))
		_ = tx.AddArgument(cadence.NewUInt64(uint64(runEnd)))
		_ = tx.AddArgument(cadence.NewBool(fatigueModeOn))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

	})

	t.Run("oracle should be able to create a fast break game", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateGameScript(env), fastBreakAddr)
		cdcId, _ := cadence.NewString(fastBreakID)
		cdcName, _ := cadence.NewString(fastBreakName)
		cdcFbrId, _ := cadence.NewString(fastBreakRunId)

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cdcName)
		_ = tx.AddArgument(cdcFbrId)
		_ = tx.AddArgument(cadence.NewBool(isPublic))
		_ = tx.AddArgument(cadence.NewUInt64(uint64(submissionDeadline)))
		_ = tx.AddArgument(cadence.NewUInt64(numPlayers))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Check that that main contract fields were initialized correctly
		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakScript(env), [][]byte{jsoncdc.MustEncode(cdcId)})
		fmt.Println(result)
		assert.NotNil(t, result)
	})

	t.Run("oracle should be able to add a stat to a fast break game", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddStatToGameScript(env), fastBreakAddr)
		cdcId, _ := cadence.NewString(fastBreakID)
		cdcName, _ := cadence.NewString(statName)
		cdcType, _ := cadence.NewString(statType)

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cdcName)
		_ = tx.AddArgument(cdcType)
		_ = tx.AddArgument(cadence.NewUInt64(valueNeeded))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

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
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
			false,
		)
	})

	// mint moment 1 to jer
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewAddress(jerAddress))

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

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cadence.NewArray(ids))

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
		ids := []cadence.Value{cadence.NewUInt64(1)}

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cadence.NewArray(ids))

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

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cadence.NewArray(ids))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, jerAddress}, []crypto.Signer{serviceKeySigner, jerSigner},
			true,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetFastBreakTokenCountScript(env), nil)
		assert.Equal(t, cadence.NewUInt64(1), result)
	})

}

package test

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/onflow/flow-go-sdk/crypto"
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
	assert.Nil(t, err)

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

	// Mint Moment 1
	{
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	}
	//momentId := uint64(1)

	var (
		fastBreakRunId          = "abc-123"
		fastBreakRunName        = "R0"
		runStart         uint64 = 1702847010
		runEnd           uint64 = 1703227010
		fatigueModeOn           = true
	)
	t.Run("Oracle should be able to create a Fast Break Run", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateRunScript(env), fastBreakAddr)
		cdcId, _ := cadence.NewString(fastBreakRunId)
		cdcName, _ := cadence.NewString(fastBreakRunName)

		_ = tx.AddArgument(cdcId)
		_ = tx.AddArgument(cdcName)
		_ = tx.AddArgument(cadence.NewUInt64(runStart))
		_ = tx.AddArgument(cadence.NewUInt64(runEnd))
		_ = tx.AddArgument(cadence.NewBool(fatigueModeOn))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, fastBreakAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		//// Verify moment is locked
		//result := executeScriptAndCheck(t, b, templates.GenerateGetMomentIsLockedScript(env), [][]byte{
		//	jsoncdc.MustEncode(cadence.Address(topshotAddr)),
		//	jsoncdc.MustEncode(cadence.UInt64(momentId)),
		//})
		//assertEqual(t, cadence.NewBool(true), result)
		//assert.Equal(t, 1, 1)
	})

}

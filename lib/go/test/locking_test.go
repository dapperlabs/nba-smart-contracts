package test

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
	"time"
)

const CadenceUFix64Factor = 100000000

// Tests all the main functionality of the V1 Market
func TestTopShotLocking(t *testing.T) {
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

	// Deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String(), metadataViewsAddr.String(), topShotLockingAddr.String())
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

	// Should be able to deploy the token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage, "1000.0")
	tokenAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "DapperUtilityCoin",
			Source: string(tokenCode),
		},
	})
	env.DUCAddress = tokenAddr.String()

	// Setup with the first market contract
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, nftAddr.String(), topshotAddr.String(), env.DUCAddress)
	marketAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "Market",
			Source: string(marketCode),
		},
	})
	env.TopShotMarketAddress = marketAddr.String()

	// Should be able to deploy the third market contract
	marketV3Code := contracts.GenerateTopShotMarketV3Contract(defaultfungibleTokenAddr, nftAddr.String(), topshotAddr.String(), marketAddr.String(), env.DUCAddress, topShotLockingAddr.String())
	marketV3Addr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "TopShotMarketV3",
			Source: string(marketV3Code),
		},
	})
	env.TopShotMarketV3Address = marketV3Addr.String()

	// Create a new user account
	//joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	//joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

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
	momentId := uint64(1)

	t.Run("Should be able to lock a moment for 1 year", func(t *testing.T) {
		expectedExpiryTime := time.Now().Add(31536000 * time.Second)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateGetTopShotLockingLockMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		duration, err := cadence.NewUFix64("31536000.0")
		if err != nil {
			log.Printf("err %+v", err)
		}
		_ = tx.AddArgument(duration)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Verify moment is locked
		result := executeScriptAndCheck(t, b, templates.GenerateGetMomentIsLockedScript(env), [][]byte{
			jsoncdc.MustEncode(cadence.Address(topshotAddr)),
			jsoncdc.MustEncode(cadence.UInt64(momentId)),
		})
		assertEqual(t, cadence.NewBool(true), result)

		// Verify moment is locked for 1 year
		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentIsLockExpiryScript(env), [][]byte{
			jsoncdc.MustEncode(cadence.Address(topshotAddr)),
			jsoncdc.MustEncode(cadence.UInt64(momentId)),
		})
		resultTime := time.Unix(int64(result.ToGoValue().(uint64)/CadenceUFix64Factor), 0)
		// Flow block time has a 10-second time accuracy, not relevant since locking is in month timescale
		assert.WithinDuration(t, expectedExpiryTime, resultTime, 10*time.Second)
	})

	t.Run("Admin should be able to mark the moment as unlockable then the owner should be able to unlock it", func(t *testing.T) {
		// locking admin marks the moment unlockable
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAdminMarkMomentUnlockableScript(env), topShotLockingAddr)

		_ = tx.AddArgument(cadence.Address(topshotAddr))
		_ = tx.AddArgument(cadence.UInt64(momentId))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topShotLockingAddr}, []crypto.Signer{serviceKeySigner, lockingSigner},
			false,
		)

		// Attempt to unlock
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateGetTopShotLockingUnlockMomentScript(env), topshotAddr)
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	t.Run("Should be able to lock a moment then unlock when the duration has expired", func(t *testing.T) {
		// Lock for 0 seconds
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateGetTopShotLockingLockMomentScript(env), topshotAddr)
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		duration, _ := cadence.NewUFix64("0.0")
		_ = tx.AddArgument(duration)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Verify that the moment is locked
		result := executeScriptAndCheck(t, b, templates.GenerateGetMomentIsLockedScript(env), [][]byte{
			jsoncdc.MustEncode(cadence.Address(topshotAddr)),
			jsoncdc.MustEncode(cadence.UInt64(momentId)),
		})
		assertEqual(t, cadence.NewBool(true), result)

		// Attempt to unlock
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateGetTopShotLockingUnlockMomentScript(env), topshotAddr)
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Verify that the moment is not locked
		result = executeScriptAndCheck(t, b, templates.GenerateGetMomentIsLockedScript(env), [][]byte{
			jsoncdc.MustEncode(cadence.Address(topshotAddr)),
			jsoncdc.MustEncode(cadence.UInt64(momentId)),
		})
		assertEqual(t, cadence.NewBool(false), result)
	})

	t.Run("Should be unable to withdraw or transfer the moment if locked", func(t *testing.T) {
		// Lock the moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateGetTopShotLockingLockMomentScript(env), topshotAddr)
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		duration, _ := cadence.NewUFix64("0.0")
		_ = tx.AddArgument(duration)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// Transfer script must fail
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)
		_ = tx.AddArgument(cadence.NewAddress(tokenAddr)) // recipient is irrelevant, MUST fail before
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			true,
		)
	})

	ducPublicPath := cadence.Path{Domain: "public", Identifier: "dapperUtilityCoinReceiver"}
	t.Run("Should be unable to list a locked moment for sale Market V1", func(t *testing.T) {
		// Note moment is locked from previous test run

		tx := createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), topshotAddr)
		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		_ = tx.AddArgument(CadenceUFix64("50.0"))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			true,
		)
	})

	t.Run("Should be unable to list a locked moment for sale MarketV3", func(t *testing.T) {
		// Note moment is locked from previous test run
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleV3Script(env), topshotAddr)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(momentId))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			true,
		)
	})
}

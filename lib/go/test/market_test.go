package test

import (
	"context"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/common"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/adapters"
	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/rs/zerolog"

	"strings"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultfungibleTokenAddr = "ee82856bf20e2aa6"
	defaultTokenName         = "ExampleToken"
	defaultTokenStorage      = "exampleToken"
)

func TestMarketDeployment(t *testing.T) {
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	topshotAddr := tb.topshotAdminAddr
	accountKeys := tb.accountKeys
	env := tb.env

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	// Should be able to deploy a token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, env.MetadataViewsAddress, env.FungibleTokenMetadataViewsAddress, "DapperUtilityCoin", "dapperUtilityCoin", "1000.0")
	tokenAccountKey, _ := accountKeys.NewWithSigner()
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

	// Should be able to deploy the market contract
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), tokenAddr.String())
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

	marketV3Code := contracts.GenerateTopShotMarketV3Contract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), marketAddr.String(), tokenAddr.String(), env.TopShotLockingAddress, env.MetadataViewsAddress)
	_, err = adapter.CreateAccount(context.Background(), nil, []sdktemplates.Contract{
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
}

// Tests all the main functionality of the V1 Market
func TestMarketV1(t *testing.T) {
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	topshotAddr := tb.topshotAdminAddr
	topshotSigner := tb.topshotAdminSigner
	serviceKeySigner := tb.serviceKeySigner
	accountKeys := tb.accountKeys
	env := tb.env

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	// Should be able to deploy the token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, env.MetadataViewsAddress, env.FungibleTokenMetadataViewsAddress, "DapperUtilityCoin", "dapperUtilityCoin", "1000.0")
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
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

	// Should be able to deploy the market contract
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, env.NFTAddress, topshotAddr.String(), tokenAddr.String())
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

	// create two new accounts
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{bastianAccountKey}, nil)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{joshAccountKey}, nil)

	//ducStoragePath := cadence.Path{Domain: "storage", Identifier: "dapperUtilityCoinVault"}
	ducPublicPath := cadence.Path{Domain: common.PathDomainPublic, Identifier: "dapperUtilityCoinReceiver"}

	tokenEnv := fungibleTokenTemplates.Environment{
		FungibleTokenAddress:              defaultfungibleTokenAddr,
		MetadataViewsAddress:              tb.metadataViewsAddr.Hex(),
		ExampleTokenAddress:               tokenAddr.Hex(),
		FungibleTokenMetadataViewsAddress: tb.fungibleTokenMetadataViewsAddr.Hex(),
		ViewResolverAddress:               tb.viewResolverAddr.Hex(),
		TokenForwardingAddress:            forwardingAddr.Hex(),
	}

	// Setup both accounts to have DUC and a sale collection
	t.Run("Should be able to setup both users' accounts to use the market", func(t *testing.T) {
		// create a Vault for bastian
		createTokenScript := fungibleTokenTemplates.GenerateCreateTokenScript(tokenEnv)
		modifiedCreateTokenScript := strings.Replace(string(createTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
		tx := createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateTokenScript), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		// create a Vault for Josh
		tx = createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateTokenScript), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		// Mint tokens to bastian's vault
		mintTokenScript := fungibleTokenTemplates.GenerateMintTokensScript(tokenEnv)
		modifiedMintTokenScript := strings.Replace(string(mintTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
		tx = createTxWithTemplateAndAuthorizer(b, []byte(modifiedMintTokenScript), tokenAddr)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("80.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

		// Create a sale collection for josh's account, setting bastian as the beneficiary
		// and with a 15% cut
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSaleScript(env), joshAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)
	})

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")

	// Admin sends transactions to create a play, set, and moments
	t.Run("Should be able to setup a play, set, and mint moment", func(t *testing.T) {
		// create a new play
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata := []cadence.KeyValuePair{{Key: firstName, Value: lebron}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// create a new set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(env), topshotAddr)

		_ = tx.AddArgument(CadenceString("Genesis"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// add the play to the set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// mint a batch of moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt64(6))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// setup bastian's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		// setup josh's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		// transfer a moment to josh's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// transfer a moment to bastian's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	t.Run("Can put an NFT up for sale", func(t *testing.T) {
		// start a sale with the moment josh owns, setting its price to 80
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateStartSaleScript(env), joshAddress)

		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("80.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)
		// check the price, sale length, and the sale's data
		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePriceScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assertEqual(t, CadenceUFix64("80.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})
		assertEqual(t, cadence.NewInt(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleSetIDScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assertEqual(t, cadence.NewUInt32(1), result)
	})

	t.Run("Cannot buy an NFT for less than the sale price", func(t *testing.T) {
		// bastian tries to buy the moment for only 9 tokens
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("9.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT for more than the sale price", func(t *testing.T) {
		// bastian tries to buy the moment for too many tokens
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("90.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT that is not for sale", func(t *testing.T) {
		// bastian tries to buy the wrong moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("80.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			true,
		)
	})

	t.Run("Can buy an NFT that is for sale", func(t *testing.T) {
		// bastian sends the correct amount of tokens to buy it
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("80.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
		inspectVaultScript := fungibleTokenTemplates.GenerateInspectVaultScript(tokenEnv)
		modifiedInspectVaultScript := strings.Replace(string(inspectVaultScript), "ExampleToken", "DapperUtilityCoin", -1)
		// make sure that the cut was taken correctly and that josh receied the purchasing tokens
		result := executeScriptAndCheck(t, b, []byte(modifiedInspectVaultScript), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assert.Equal(t, CadenceUFix64("12.0"), result)
		result = executeScriptAndCheck(t, b, []byte(modifiedInspectVaultScript), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})
		assert.Equal(t, CadenceUFix64("68.0"), result)

		// make sure bastian received the purchase's moment
		result = executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.NewBool(true), result)
	})

	t.Run("Can create a sale and put an NFT up for sale in one transaction", func(t *testing.T) {
		// Bastian creates a new sale collection object and puts the moment for sale,
		// setting himself as the beneficiary with a 15% cut
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), bastianAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
		// Make sure that moment id 2 is for sale for 50 tokens and the data is correct
		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePriceScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, CadenceUFix64("50.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, cadence.NewInt(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleSetIDScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, cadence.NewUInt32(1), result)
	})

	t.Run("Cannot change the price of a moment that isn't for sale", func(t *testing.T) {
		// try to change the price of the wrong moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePriceScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(CadenceUFix64("40.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			true,
		)
	})

	t.Run("Can change the price of a sale", func(t *testing.T) {
		// change the price of the moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePriceScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("40.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
		// make sure the price has been changed
		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePriceScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, CadenceUFix64("40.0"), result)
	})

	t.Run("Can change the cut percentage of a sale", func(t *testing.T) {
		// change the cut percentage for the sale collection to 18%
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePercentageScript(env), bastianAddress)

		_ = tx.AddArgument(CadenceUFix64("0.18"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
		// make sure the percentage was changed correctly
		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePercentageScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, CadenceUFix64("0.18"), result)
	})

	t.Run("Cannot withdraw a moment that doesn't exist from a sale", func(t *testing.T) {
		// bastian tries to withdraw the wrong moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateWithdrawFromSaleScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(7))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			true,
		)
		// make sure nothing was withdrawn
		result := executeScriptAndCheck(t, b, templates.GenerateGetSaleLenScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, cadence.NewInt(1), result)
	})

	t.Run("Can withdraw a moment from a sale", func(t *testing.T) {
		// bastian withdraws the correct moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateWithdrawFromSaleScript(env), bastianAddress)
		_ = tx.AddArgument(cadence.NewUInt64(2))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetSaleLenScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, cadence.NewInt(0), result)
	})

	t.Run("Can use the create and start sale to start a sale even if there is already sale in storage", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), bastianAddress)
		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.10"))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("100.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
		// Make sure that moment id 2 is for sale for 50 tokens and the data is correct
		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePriceScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, CadenceUFix64("100.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, cadence.NewInt(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleSetIDScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, cadence.NewUInt32(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSalePercentageScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, CadenceUFix64("0.18"), result)
	})

	t.Run("Can create a forwarder resource to forward tokens to a different account", func(t *testing.T) {
		createForwarderScript := fungibleTokenTemplates.GenerateCreateForwarderScript(tokenEnv)
		modifiedCreateForwarderScript := strings.Replace(string(createForwarderScript), "ExampleToken", "DapperUtilityCoin", -1)
		tx := createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateForwarderScript), bastianAddress)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
	})

	t.Run("Can change the owner capability of a sale", func(t *testing.T) {
		// change the price of the moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeOwnerReceiverScript(env), bastianAddress)
		_ = tx.AddArgument(ducPublicPath)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
		inspectVaultScript := fungibleTokenTemplates.GenerateInspectVaultScript(tokenEnv)
		modifiedInspectVaultScript := strings.Replace(string(inspectVaultScript), "ExampleToken", "DapperUtilityCoin", -1)
		executeScriptAndCheck(t, b, []byte(modifiedInspectVaultScript), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
	})

	t.Run("Can mint tokens and buy a moment with them so the tokens are forwarded", func(t *testing.T) {

		// mint tokens and buy the moment in the same tx

		template := templates.GenerateMintTokensAndBuyScript(env)

		tx := createTxWithTemplateAndAuthorizer(b, template, tokenAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("100.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

		// make sure josh received the purchase's moment
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assert.Equal(t, cadence.NewBool(true), result)

		inspectVaultScript := fungibleTokenTemplates.GenerateInspectVaultScript(tokenEnv)
		modifiedInspectVaultScript := strings.Replace(string(inspectVaultScript), "ExampleToken", "DapperUtilityCoin", -1)
		executeScriptAndCheck(t, b, []byte(modifiedInspectVaultScript), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})
	})
}

func TestMarketV3(t *testing.T) {
	tb := NewTopShotTestBlockchain(t)
	b := tb.Blockchain
	topshotAddr := tb.topshotAdminAddr
	topshotSigner := tb.topshotAdminSigner
	serviceKeySigner := tb.serviceKeySigner
	accountKeys := tb.accountKeys
	env := tb.env

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	// Should be able to deploy the token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, env.MetadataViewsAddress, env.FungibleTokenMetadataViewsAddress, "DapperUtilityCoin", "dapperUtilityCoin", "1000.0")
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
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

	// create two new accounts
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{bastianAccountKey}, nil)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{joshAccountKey}, nil)

	tokenEnv := fungibleTokenTemplates.Environment{
		FungibleTokenAddress:              defaultfungibleTokenAddr,
		MetadataViewsAddress:              tb.metadataViewsAddr.Hex(),
		ExampleTokenAddress:               tokenAddr.Hex(),
		FungibleTokenMetadataViewsAddress: tb.fungibleTokenMetadataViewsAddr.Hex(),
		ViewResolverAddress:               tb.viewResolverAddr.Hex(),
		TokenForwardingAddress:            forwardingAddr.Hex(),
	}

	// Setup both accounts to have DUC and a sale collection
	t.Run("Should be able to setup both users' accounts to use DUC", func(t *testing.T) {

		createTokenScript := fungibleTokenTemplates.GenerateCreateTokenScript(tokenEnv)
		modifiedCreateTokenScript := strings.Replace(string(createTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
		tx := createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateTokenScript), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		// create a Vault for Josh
		tx = createTxWithTemplateAndAuthorizer(b, []byte(modifiedCreateTokenScript), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		mintTokenScript := fungibleTokenTemplates.GenerateMintTokensScript(tokenEnv)
		modifiedMintTokenScript := strings.Replace(string(mintTokenScript), "ExampleToken", "DapperUtilityCoin", -1)
		tx = createTxWithTemplateAndAuthorizer(b, []byte(modifiedMintTokenScript), tokenAddr)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("80.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)
	})

	firstName := CadenceString("FullName")
	lebron := CadenceString("Lebron")

	// Admin sends transactions to create a play, set, and moments
	t.Run("Should be able to setup a play, set, and mint moment", func(t *testing.T) {
		// create a new play
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata := []cadence.KeyValuePair{{Key: firstName, Value: lebron}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// create a new set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(env), topshotAddr)

		_ = tx.AddArgument(CadenceString("Genesis"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// add the play to the set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// mint a batch of moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt64(17))
		_ = tx.AddArgument(cadence.NewAddress(topshotAddr))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// setup bastian's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		// setup josh's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), tokenAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

		// transfer a moment to josh's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// transfer a moment to bastian's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)

		// transfer a moment to bastian's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewUInt64(3))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{serviceKeySigner, topshotSigner},
			false,
		)
	})

	ducPublicPath := cadence.Path{Domain: common.PathDomainPublic, Identifier: "dapperUtilityCoinReceiver"}

	t.Run("Should be able to create a V1 sale collection and V3 sale collection in the same account", func(t *testing.T) {

		// Bastian creates a new sale collection object and puts the moment for sale,
		// setting himself as the beneficiary with a 15% cut
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), bastianAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		// Create a sale collection for josh's account, setting bastian as the beneficiary
		// and with a 15% cut
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleV3Script(env), bastianAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(3))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
	})

	t.Run("Should not be able to put a moment up for sale in v3 that isn't in the main collection", func(t *testing.T) {

		// Should fail because the moment isn't in the user's collection
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleV3Script(env), bastianAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(4))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			true,
		)
	})

	t.Run("Should be able to cancel sales in the v1 and v3 collections", func(t *testing.T) {

		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePriceV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, CadenceUFix64("50.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleLenV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, cadence.NewInt(2), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSaleSetIDV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, cadence.NewUInt32(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSalePercentageV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress))})
		assertEqual(t, CadenceUFix64("0.15"), result)

		// Should not fail if an ID is not for sale
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCancelSaleV3Script(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(4))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCancelSaleV3Script(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(3))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCancelSaleV3Script(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(2))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)
	})

	t.Run("Should start the sales again and purchase", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), bastianAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateStartSaleV3Script(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewUInt64(3))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{serviceKeySigner, bastianSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetSalePriceV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assertEqual(t, CadenceUFix64("50.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetSalePriceV3Script(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(3))})
		assertEqual(t, CadenceUFix64("50.0"), result)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintTokensAndBuyV3Script(env), tokenAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(3))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

		// Should fail because the price is wrong
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMultiContractP2PPurchaseScript(env), tokenAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("40.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			true,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMultiContractP2PPurchaseScript(env), tokenAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

	})

	t.Run("Should fail purchases for tokens that don't exist in the collection", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMultiContractP2PPurchaseScript(env), tokenAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(5))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			true,
		)
	})

	t.Run("V3 transactions should still work for V1", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), joshAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("50.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCancelSaleV3Script(env), joshAddress)

		_ = tx.AddArgument(cadence.NewUInt64(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(env), joshAddress)

		_ = tx.AddArgument(ducPublicPath)
		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(CadenceUFix64("0.15"))
		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("60.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{serviceKeySigner, joshSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMultiContractP2PPurchaseScript(env), tokenAddr)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))
		_ = tx.AddArgument(CadenceUFix64("60.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)
	})

	t.Run("Purchase a group of moments from market v1", func(t *testing.T) {
		// transfer and start sale for 2 moments with a price of 50.00 each
		// starting at moment index 6 for josh via market v1
		transferAndStartSale(t,
			b,
			env,
			1,
			6,
			2,
			"50.00",
			topshotAddr,
			topshotSigner,
			joshAddress,
			joshSigner,
			serviceKeySigner)

		// transfer and start sale for 2 moments with a price of 100.00 each
		// starting at moment index 8 for bastian via market v1
		transferAndStartSale(t,
			b,
			env,
			1,
			8,
			2,
			"100.00",
			topshotAddr,
			topshotSigner,
			bastianAddress,
			bastianSigner,
			serviceKeySigner)

		// Buy group of moments from both josh & bastian
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePurchaseGroupOfMomentsScript(env), tokenAddr)

		joshMomentIDs := []cadence.Value{cadence.NewUInt64(6), cadence.NewUInt64(7)}
		bastianMomentIDs := []cadence.Value{cadence.NewUInt64(8), cadence.NewUInt64(9)}

		group := cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: cadence.NewAddress(joshAddress), Value: cadence.NewArray(joshMomentIDs)},
			{Key: cadence.NewAddress(bastianAddress), Value: cadence.NewArray(bastianMomentIDs)},
		})

		_ = tx.AddArgument(group)
		_ = tx.AddArgument(CadenceUFix64("300.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

		// make sure non-fungible receiver received the purchased moments
		for i := 6; i < 10; i++ {
			result := executeScriptAndCheck(t,
				b,
				templates.GenerateIsIDInCollectionScript(env),
				[][]byte{jsoncdc.MustEncode(cadence.Address(tokenAddr)),
					jsoncdc.MustEncode(cadence.UInt64(i))},
			)
			assert.Equal(t, cadence.NewBool(true), result)
		}
	})

	t.Run("Purchase a group of moments from market V1 or V3", func(t *testing.T) {
		// transfer and start sale for 2 moments with a price of 50.00 each
		// starting at index 10 for josh via market v1
		transferAndStartSale(t,
			b,
			env,
			1,
			10,
			2,
			"50.00",
			topshotAddr,
			topshotSigner,
			joshAddress,
			joshSigner,
			serviceKeySigner)

		// transfer and start sale for 2 moments with a price of 100.00 each
		// starting at index 12 for bastian via market v3
		transferAndStartSale(t,
			b,
			env,
			3,
			12,
			2,
			"100.00",
			topshotAddr,
			topshotSigner,
			bastianAddress,
			bastianSigner,
			serviceKeySigner)

		// Buy group of moments from both josh & bastian
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePurchaseGroupOfMomentsScript(env), tokenAddr)

		joshMomentIDs := []cadence.Value{cadence.NewUInt64(10), cadence.NewUInt64(11)}
		bastianMomentIDs := []cadence.Value{cadence.NewUInt64(12), cadence.NewUInt64(13)}

		group := cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: cadence.NewAddress(joshAddress), Value: cadence.NewArray(joshMomentIDs)},
			{Key: cadence.NewAddress(bastianAddress), Value: cadence.NewArray(bastianMomentIDs)},
		})

		_ = tx.AddArgument(group)
		_ = tx.AddArgument(CadenceUFix64("300.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			false,
		)

		// make sure non-fungible receiver received the purchased moments from both markets
		for i := 10; i < 14; i++ {
			result := executeScriptAndCheck(t,
				b,
				templates.GenerateIsIDInCollectionScript(env),
				[][]byte{jsoncdc.MustEncode(cadence.Address(tokenAddr)),
					jsoncdc.MustEncode(cadence.UInt64(i))},
			)
			assert.Equal(t, cadence.NewBool(true), result)
		}
	})

	t.Run("Should fail: Purchase a group of moments with a purchaseAmount lesser than the sum of all moment prices", func(t *testing.T) {
		// transfer and start sale for 2 moments with a price of 50.00 each
		// starting at index 14 for josh via market v1
		transferAndStartSale(t,
			b,
			env,
			1,
			14,
			2,
			"50.00",
			topshotAddr,
			topshotSigner,
			joshAddress,
			joshSigner,
			serviceKeySigner)

		// transfer and start sale for 2 moments with a price of 100.00 each
		// starting at index 16 for bastian via market v1
		transferAndStartSale(t,
			b,
			env,
			1,
			16,
			2,
			"100.00",
			topshotAddr,
			topshotSigner,
			bastianAddress,
			bastianSigner,
			serviceKeySigner)

		// Buy group of moments from both josh & bastian
		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePurchaseGroupOfMomentsScript(env), tokenAddr)

		joshMomentIDs := []cadence.Value{cadence.NewUInt64(14), cadence.NewUInt64(15)}
		bastianMomentIDs := []cadence.Value{cadence.NewUInt64(16), cadence.NewUInt64(17)}

		group := cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: cadence.NewAddress(joshAddress), Value: cadence.NewArray(joshMomentIDs)},
			{Key: cadence.NewAddress(bastianAddress), Value: cadence.NewArray(bastianMomentIDs)},
		})

		_ = tx.AddArgument(group)
		_ = tx.AddArgument(CadenceUFix64("200.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{serviceKeySigner, tokenSigner},
			true,
		)

		// make sure non-fungible receiver did not receive the purchased moments from the failed TX
		for i := 14; i < 18; i++ {
			result := executeScriptAndCheck(
				t,
				b,
				templates.GenerateIsIDInCollectionScript(env),
				[][]byte{jsoncdc.MustEncode(cadence.Address(tokenAddr)),
					jsoncdc.MustEncode(cadence.UInt64(i))},
			)
			assert.Equal(t, cadence.NewBool(false), result)
		}
	})
}

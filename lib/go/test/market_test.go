package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/lib/go/templates"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultfungibleTokenAddr = "ee82856bf20e2aa6"
	defaultTokenName         = "DapperUtilityCoin"
	defaultTokenStorage      = "dapperUtilityCoin"
)

func TestMarketDeployment(t *testing.T) {
	b := newBlockchain()

	// Should be able to deploy the NFT contract
	// as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the topshot contract
	// as a new account with no keys.
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a token contract
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage, "1000.0")
	_, err = b.CreateAccount(nil, []sdktemplates.Contract{
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
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, nftAddr.String(), topshotAddr.String())
	_, err = b.CreateAccount(nil, []sdktemplates.Contract{
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
}

// Tests all the main functionality of Sale Collections
func TestMarket(t *testing.T) {
	b := newBlockchain()

	accountKeys := test.AccountKeyGenerator()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	// Should be able to deploy the NFT contract
	// as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.NFTAddress = nftAddr.String()

	// Should be able to deploy the topshot contract
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, err := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, []sdktemplates.Contract{
		{
			Name:   "TopShot",
			Source: string(topshotCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	env.TopShotAddress = topshotAddr.String()

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

	// Should be able to deploy the token contract
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, nftAddr.String(), topshotAddr.String())
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

	// create two new accounts
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, err := b.CreateAccount([]*flow.AccountKey{bastianAccountKey}, nil)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	//ducStoragePath := cadence.Path{Domain: "storage", Identifier: "dapperUtilityCoinVault"}
	ducPublicPath := cadence.Path{Domain: "public", Identifier: "dapperUtilityCoinReceiver"}

	// Setup both accounts to have DUC and a sale collection
	t.Run("Should be able to setup both users' accounts to use the market", func(t *testing.T) {

		// create a Vault for bastian
		tx := createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// create a Vault for Josh
		tx = createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		// Mint tokens to bastian's vault
		tx = createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateMintTokensScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, 80), tokenAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
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
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	firstName := cadence.NewString("FullName")
	lebron := cadence.NewString("Lebron")

	// Admin sends transactions to create a play, set, and moments
	t.Run("Should be able to setup a play, set, and mint moment", func(t *testing.T) {
		// create a new play
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(env), topshotAddr)

		metadata := []cadence.KeyValuePair{{Key: firstName, Value: lebron}}
		play := cadence.NewDictionary(metadata)
		_ = tx.AddArgument(play)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// create a new set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewString("Genesis"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// add the play to the set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewUInt32(1))
		_ = tx.AddArgument(cadence.NewUInt32(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
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
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// setup bastian's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// setup josh's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(env), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		// transfer a moment to josh's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(joshAddress))
		_ = tx.AddArgument(cadence.NewUInt64(1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// transfer a moment to bastian's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(env), topshotAddr)

		_ = tx.AddArgument(cadence.NewAddress(bastianAddress))
		_ = tx.AddArgument(cadence.NewUInt64(2))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
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
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// make sure that the cut was taken correctly and that josh receied the purchasing tokens
		executeScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, 12), nil)
		executeScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, joshAddress, defaultTokenName, 68), nil)

		// make sure bastian received the purchase's moment
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(bastianAddress)), jsoncdc.MustEncode(cadence.UInt64(1))})
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
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
		tx := createTxWithTemplateAndAuthorizer(b, fungibleTokenTemplates.GenerateCreateForwarderScript(flow.HexToAddress(defaultfungibleTokenAddr), forwardingAddr, tokenAddr, "DapperUtilityCoin"), bastianAddress)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
	})

	t.Run("Can change the owner capability of a sale", func(t *testing.T) {
		// change the price of the moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeOwnerReceiverScript(env), bastianAddress)
		_ = tx.AddArgument(ducPublicPath)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		executeScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, tokenAddr, defaultTokenName, 1000.0), nil)
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
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		// make sure josh received the purchase's moment
		result := executeScriptAndCheck(t, b, templates.GenerateIsIDInCollectionScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress)), jsoncdc.MustEncode(cadence.UInt64(2))})
		assert.Equal(t, cadence.NewBool(true), result)

		executeScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, tokenAddr, defaultTokenName, 1100.0), nil)
	})
}

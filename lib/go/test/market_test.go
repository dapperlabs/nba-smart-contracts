package test

import (
	"testing"

	fungibleToken "github.com/onflow/flow-ft/lib/go/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/lib/go/templates"

	"github.com/dapperlabs/nba-smart-contracts/lib/go/contracts"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/templates/data"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
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
	b := NewEmulator()

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
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

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

	// Should be able to deploy the token forwarding contract
	forwardingCode := fungibleToken.CustomTokenForwarding(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage)
	forwardingAddr, err := b.CreateAccount(nil, forwardingCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// create two new accounts
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, err := b.CreateAccount([]*flow.AccountKey{bastianAccountKey}, nil)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// Setup both accounts to have DUC and a sale collection
	t.Run("Should be able to setup both users' accounts to use the market", func(t *testing.T) {

		// create a Vault for bastian
		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// create a Vault for Josh
		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		// Mint tokens to bastian's vault
		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateMintTokensScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, 80),
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		// Create a sale collection for josh's account, setting bastian as the beneficiary
		// and with a 15% cut
		createSignAndSubmit(
			t, b,
			templates.GenerateCreateSaleScript(marketAddr, bastianAddress, defaultTokenStorage, .15),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	// Admin sends transactions to create a play, set, and moments
	t.Run("Should be able to setup a play, set, and mint moment", func(t *testing.T) {
		// create a new play
		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.GenerateEmptyPlay("Lebron")),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// create a new set
		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Genesis"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// add the play to the set
		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// mint a batch of moments
		createSignAndSubmit(
			t, b,
			templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 1, 1, 6),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// setup bastian's account to hold topshot moments
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// setup josh's account to hold topshot moments
		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		// transfer a moment to josh's account
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentScript(nftAddr, topshotAddr, joshAddress, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// transfer a moment to bastian's account
		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentScript(nftAddr, topshotAddr, bastianAddress, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	t.Run("Can put an NFT up for sale", func(t *testing.T) {
		// start a sale with the moment josh owns, setting its price to 80
		createSignAndSubmit(
			t, b,
			templates.GenerateStartSaleScript(topshotAddr, marketAddr, 1, 80),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
		// check the price, sale length, and the sale's data
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, joshAddress, 1, 80), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, joshAddress, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, joshAddress, 1, 1), false)
	})

	t.Run("Cannot buy an NFT for less than the sale price", func(t *testing.T) {
		// bastian tries to buy the moment for only 9 tokens
		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 9),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT for more than the sale price", func(t *testing.T) {
		// bastian tries to buy the moment for too many tokens
		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 90),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT that is not for sale", func(t *testing.T) {
		// bastian tries to buy the wrong moment
		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 2, 80),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Can buy an NFT that is for sale", func(t *testing.T) {
		// bastian sends the correct amount of tokens to buy it
		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 80),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// make sure that the cut was taken correctly and that josh receied the purchasing tokens
		ExecuteScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, 12), false)
		ExecuteScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, joshAddress, defaultTokenName, 68), false)

		// make sure bastian received the purchase's moment
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, bastianAddress, 1), false)
	})

	t.Run("Can create a sale and put an NFT up for sale in one transaction", func(t *testing.T) {
		// Bastian creates a new sale collection object and puts the moment for sale,
		// setting himself as the beneficiary with a 15% cut
		createSignAndSubmit(
			t, b,
			templates.GenerateCreateAndStartSaleScript(topshotAddr, marketAddr, bastianAddress, defaultTokenStorage, .15, 2, 50),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// Make sure that moment id 2 is for sale for 50 tokens and the data is correct
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), false)
	})

	t.Run("Cannot change the price of a moment that isn't for sale", func(t *testing.T) {
		// try to change the price of the wrong moment
		createSignAndSubmit(
			t, b,
			templates.GenerateChangePriceScript(topshotAddr, marketAddr, 5, 40),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), false)
	})

	t.Run("Can change the price of a sale", func(t *testing.T) {
		// change the price of the moment
		createSignAndSubmit(
			t, b,
			templates.GenerateChangePriceScript(topshotAddr, marketAddr, 2, 40),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// make sure the price has been changed
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 40), false)
	})

	t.Run("Can change the cut percentage of a sale", func(t *testing.T) {
		// change the cut percentage for the sale collection to 18%
		createSignAndSubmit(
			t, b,
			templates.GenerateChangePercentageScript(topshotAddr, marketAddr, .18),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// make sure the percentage was changed correctly
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSalePercentageScript(marketAddr, bastianAddress, .18), false)
	})

	t.Run("Cannot withdraw a moment that doesn't exist from a sale", func(t *testing.T) {
		// bastian tries to withdraw the wrong moment
		createSignAndSubmit(
			t, b,
			templates.GenerateWithdrawFromSaleScript(topshotAddr, marketAddr, 7),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
		// make sure nothing was withdrawn
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
	})

	t.Run("Can withdraw a moment from a sale", func(t *testing.T) {
		// bastian withdraws the correct moment
		createSignAndSubmit(
			t, b,
			templates.GenerateWithdrawFromSaleScript(topshotAddr, marketAddr, 2),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), true)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 0), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), true)
	})

	t.Run("Can use the create and start sale to start a sale even if there is already sale in storage", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateCreateAndStartSaleScript(topshotAddr, marketAddr, bastianAddress, defaultTokenStorage, .10, 2, 100),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// Make sure that moment id 2 is for sale for 50 tokens and the data is correct
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 100), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), false)

		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSalePercentageScript(marketAddr, bastianAddress, .10), false)
	})

	t.Run("Can create a forwarder resource to forward tokens to a different account", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateCreateForwarderScript(flow.HexToAddress(defaultfungibleTokenAddr), forwardingAddr, tokenAddr, "DapperUtilityCoin"),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
	})

	t.Run("Can change the owner capability of a sale", func(t *testing.T) {
		// change the price of the moment
		createSignAndSubmit(
			t, b,
			templates.GenerateChangeOwnerReceiverScript(flow.HexToAddress(defaultfungibleTokenAddr), topshotAddr, marketAddr, "dapperUtilityCoinReceiver"),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		ExecuteScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, tokenAddr, defaultTokenName, 1000.0), false)
	})

	t.Run("Can mint tokens and buy a moment with them so the tokens are forwarded", func(t *testing.T) {

		// mint tokens and buy the moment in the same tx

		template := templates.GenerateMintTokensAndBuyScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, bastianAddress, joshAddress, defaultTokenName, defaultTokenStorage, 2, 100)

		createSignAndSubmit(
			t, b,
			template,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		// make sure josh received the purchase's moment
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 2), false)

		ExecuteScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, tokenAddr, defaultTokenName, 1100.0), false)
	})
}

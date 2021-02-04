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

	// create two new accounts
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, err := b.CreateAccount([]*flow.AccountKey{bastianAccountKey}, nil)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

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
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateSaleScript(marketAddr, bastianAddress, defaultTokenStorage, .15), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	// Admin sends transactions to create a play, set, and moments
	t.Run("Should be able to setup a play, set, and mint moment", func(t *testing.T) {
		// create a new play
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMintPlayScript(topshotAddr, data.GenerateEmptyPlay("Lebron")), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// create a new set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateMintSetScript(topshotAddr, "Genesis"), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// add the play to the set
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// mint a batch of moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 1, 1, 6), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// setup bastian's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(nftAddr, topshotAddr), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// setup josh's account to hold topshot moments
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateSetupAccountScript(nftAddr, topshotAddr), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		// transfer a moment to josh's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(nftAddr, topshotAddr, joshAddress, 1), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		// transfer a moment to bastian's account
		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateTransferMomentScript(nftAddr, topshotAddr, bastianAddress, 2), topshotAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	t.Run("Can put an NFT up for sale", func(t *testing.T) {
		// start a sale with the moment josh owns, setting its price to 80
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateStartSaleScript(topshotAddr, marketAddr, 1, 80), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
		// check the price, sale length, and the sale's data
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleScript(marketAddr, joshAddress, 1, 80), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleLenScript(marketAddr, joshAddress, 1), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, joshAddress, 1, 1), false)
	})

	t.Run("Cannot buy an NFT for less than the sale price", func(t *testing.T) {
		// bastian tries to buy the moment for only 9 tokens
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 9), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT for more than the sale price", func(t *testing.T) {
		// bastian tries to buy the moment for too many tokens
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 90), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT that is not for sale", func(t *testing.T) {
		// bastian tries to buy the wrong moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 2, 80), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Can buy an NFT that is for sale", func(t *testing.T) {
		// bastian sends the correct amount of tokens to buy it
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateBuySaleScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 80), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		// make sure that the cut was taken correctly and that josh receied the purchasing tokens
		ExecuteScriptAndCheckShouldFail(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, 12), false)
		ExecuteScriptAndCheckShouldFail(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, joshAddress, defaultTokenName, 68), false)

		// make sure bastian received the purchase's moment
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, bastianAddress, 1), false)
	})

	t.Run("Can create a sale and put an NFT up for sale in one transaction", func(t *testing.T) {
		// Bastian creates a new sale collection object and puts the moment for sale,
		// setting himself as the beneficiary with a 15% cut
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(topshotAddr, marketAddr, bastianAddress, defaultTokenStorage, .15, 2, 50), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// Make sure that moment id 2 is for sale for 50 tokens and the data is correct
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), false)
	})

	t.Run("Cannot change the price of a moment that isn't for sale", func(t *testing.T) {
		// try to change the price of the wrong moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePriceScript(topshotAddr, marketAddr, 5, 40), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), false)
	})

	t.Run("Can change the price of a sale", func(t *testing.T) {
		// change the price of the moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePriceScript(topshotAddr, marketAddr, 2, 40), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// make sure the price has been changed
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 40), false)
	})

	t.Run("Can change the cut percentage of a sale", func(t *testing.T) {
		// change the cut percentage for the sale collection to 18%
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePercentageScript(topshotAddr, marketAddr, .18), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// make sure the percentage was changed correctly
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSalePercentageScript(marketAddr, bastianAddress, .18), false)
	})

	t.Run("Cannot withdraw a moment that doesn't exist from a sale", func(t *testing.T) {
		// bastian tries to withdraw the wrong moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateWithdrawFromSaleScript(topshotAddr, marketAddr, 7), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
		// make sure nothing was withdrawn
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
	})

	t.Run("Can withdraw a moment from a sale", func(t *testing.T) {
		// bastian withdraws the correct moment
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateWithdrawFromSaleScript(topshotAddr, marketAddr, 2), bastianAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), true)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 0), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), true)
	})

	t.Run("Can use the create and start sale to start a sale even if there is already sale in storage", func(t *testing.T) {
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateCreateAndStartSaleScript(topshotAddr, marketAddr, bastianAddress, defaultTokenStorage, .10, 2, 100), bastianAddress)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		// Make sure that moment id 2 is for sale for 50 tokens and the data is correct
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 100), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), false)

		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectSalePercentageScript(marketAddr, bastianAddress, .10), false)
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
		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeOwnerReceiverScript(flow.HexToAddress(defaultfungibleTokenAddr), topshotAddr, marketAddr, "dapperUtilityCoinReceiver"), bastianAddress)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		ExecuteScriptAndCheckShouldFail(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, tokenAddr, defaultTokenName, 1000.0), false)
	})

	t.Run("Can mint tokens and buy a moment with them so the tokens are forwarded", func(t *testing.T) {

		// mint tokens and buy the moment in the same tx

		template := templates.GenerateMintTokensAndBuyScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, topshotAddr, marketAddr, bastianAddress, joshAddress, defaultTokenName, defaultTokenStorage, 2, 100)

		tx := createTxWithTemplateAndAuthorizer(b, template, tokenAddr)
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		// make sure josh received the purchase's moment
		ExecuteScriptAndCheckShouldFail(t, b, templates.GenerateInspectCollectionScript(nftAddr, topshotAddr, joshAddress, 2), false)

		ExecuteScriptAndCheckShouldFail(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, tokenAddr, defaultTokenName, 1100.0), false)
	})
}

package test

import (
	"testing"

	fungibleToken "github.com/onflow/flow-ft/contracts"
	fungibleTokenTemplates "github.com/onflow/flow-ft/test"

	"github.com/dapperlabs/nba-smart-contracts/contracts"
	"github.com/dapperlabs/nba-smart-contracts/templates"
	"github.com/dapperlabs/nba-smart-contracts/templates/data"
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
	nftAddr, err := b.CreateAccount(nil, nftCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the topshot contract
	// as a new account with no keys.
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAddr, err := b.CreateAccount(nil, topshotCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage)
	tokenAddr, err := b.CreateAccount(nil, []byte(tokenCode))
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, nftAddr.String(), topshotAddr.String(), tokenAddr.String())
	_, err = b.CreateAccount(nil, marketCode)
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)
}

func TestMarket(t *testing.T) {
	b := NewEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy the NFT contract
	// as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, err := b.CreateAccount(nil, nftCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy the topshot contract
	// as a new account with no keys.
	topshotCode := contracts.GenerateTopShotContract(nftAddr.String())
	topshotAccountKey, topshotSigner := accountKeys.NewWithSigner()
	topshotAddr, err := b.CreateAccount([]*flow.AccountKey{topshotAccountKey}, topshotCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	tokenCode := fungibleToken.CustomToken(defaultfungibleTokenAddr, defaultTokenName, defaultTokenStorage)
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
	tokenAddr, err := b.CreateAccount([]*flow.AccountKey{tokenAccountKey}, []byte(tokenCode))
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	marketCode := contracts.GenerateTopShotMarketContract(defaultfungibleTokenAddr, nftAddr.String(), topshotAddr.String(), tokenAddr.String())
	marketAddr, err := b.CreateAccount(nil, marketCode)
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	// create two new accounts
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, err := b.CreateAccount([]*flow.AccountKey{bastianAccountKey}, nil)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// Admin sends a transaction to create a play
	t.Run("Should be able to setup both users' accounts to use the market", func(t *testing.T) {

		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName, defaultTokenStorage),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateCreateTokenScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, defaultTokenName, defaultTokenStorage),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			fungibleTokenTemplates.GenerateMintTokensScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, defaultTokenStorage, 80),
			[]flow.Address{b.ServiceKey().Address, tokenAddr}, []crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateCreateSaleScript(marketAddr, bastianAddress, defaultTokenStorage, .15),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	// Admin sends a transaction to create a play
	t.Run("Should be able to setup a play, set, and mint moment", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateMintPlayScript(topshotAddr, data.PlayMetadata{FullName: "Lebron"}),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateMintSetScript(topshotAddr, "Genesis"),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateAddPlayToSetScript(topshotAddr, 1, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateBatchMintMomentScript(topshotAddr, topshotAddr, 1, 1, 6),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateSetupAccountScript(nftAddr, topshotAddr),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentScript(nftAddr, topshotAddr, joshAddress, 1),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)

		createSignAndSubmit(
			t, b,
			templates.GenerateTransferMomentScript(nftAddr, topshotAddr, bastianAddress, 2),
			[]flow.Address{b.ServiceKey().Address, topshotAddr}, []crypto.Signer{b.ServiceKey().Signer(), topshotSigner},
			false,
		)
	})

	t.Run("Can put an NFT up for sale", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateStartSaleScript(topshotAddr, marketAddr, 1, 80),
			[]flow.Address{b.ServiceKey().Address, joshAddress}, []crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, joshAddress, 1, 80), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, joshAddress, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, joshAddress, 1, 1), false)
	})

	t.Run("Cannot buy an NFT for less than the sale price", func(t *testing.T) {

		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 9),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Cannot buy an NFT that is not for sale", func(t *testing.T) {

		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 2, 80),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			true,
		)
	})

	t.Run("Can buy an NFT that is for sale", func(t *testing.T) {

		createSignAndSubmit(
			t, b,
			templates.GenerateBuySaleScript(tokenAddr, topshotAddr, marketAddr, joshAddress, defaultTokenName, defaultTokenStorage, 1, 80),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		ExecuteScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, bastianAddress, defaultTokenName, defaultTokenStorage, 12), false)
		ExecuteScriptAndCheck(t, b, fungibleTokenTemplates.GenerateInspectVaultScript(flow.HexToAddress(defaultfungibleTokenAddr), tokenAddr, joshAddress, defaultTokenName, defaultTokenStorage, 68), false)
	})

	t.Run("Can create a sale and put an NFT up for sale in one transaction", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateCreateAndStartSaleScript(topshotAddr, marketAddr, bastianAddress, defaultTokenStorage, .15, 2, 50),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 50), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1), false)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleMomentDataScript(nftAddr, topshotAddr, marketAddr, bastianAddress, 2, 1), false)
	})

	t.Run("Can change the price of a sale", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateChangePriceScript(topshotAddr, marketAddr, 2, 40),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSaleScript(marketAddr, bastianAddress, 2, 40), false)
	})

	t.Run("Can change the cut percentage of a sale", func(t *testing.T) {
		createSignAndSubmit(
			t, b,
			templates.GenerateChangePercentageScript(topshotAddr, marketAddr, .18),
			[]flow.Address{b.ServiceKey().Address, bastianAddress}, []crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)
		ExecuteScriptAndCheck(t, b, templates.GenerateInspectSalePercentageScript(marketAddr, bastianAddress, .18), false)
	})
}

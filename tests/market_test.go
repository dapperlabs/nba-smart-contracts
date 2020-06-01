package tests

import (
	"testing"

	"github.com/dapperlabs/nba-smart-contracts/contracts"
	"github.com/stretchr/testify/assert"
)

const (
	FungibleTokenContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-ft/master/src/contracts/"
	FungibleTokenInterfaceFile    = "FungibleToken.cdc"
	FlowTokenFile                 = "FlowToken.cdc"
	MarketContractFile            = "../contracts/MarketTopShot.cdc"
)

func TestMarketDeployment(t *testing.T) {
	b := NewEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode, _ := DownloadFile(NonFungibleTokenContractsBaseURL + NonFungibleTokenInterfaceFile)
	nftAddr, err := b.CreateAccount(nil, nftCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	topshotCode := contracts.GenerateTopShotContract(nftAddr)
	_, err = b.CreateAccount(nil, topshotCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	ftCode, _ := DownloadFile(FungibleTokenContractsBaseURL + FungibleTokenInterfaceFile)
	_, err = b.CreateAccount(nil, ftCode)
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	// flowTokenCode, _ := DownloadFile(FungibleTokenContractsBaseURL + FlowTokenFile)
	// codeWithFTAddr := strings.ReplaceAll(string(flowTokenCode), "0x02", "0x04")
	// _, err = b.CreateAccount(nil, []byte(codeWithFTAddr))
	// if !assert.NoError(t, err) {
	// 	t.Log(err.Error())
	// }
	// _, err = b.CommitBlock()
	// assert.NoError(t, err)

	// // Should be able to deploy a contract as a new account with no keys.
	// marketCode := ReadFile(MarketContractFile)
	// _, err = b.CreateAccount(nil, marketCode)
	// if !assert.Nil(t, err) {
	// 	t.Log(err.Error())
	// }
	// _, err = b.CommitBlock()
	// require.NoError(t, err)
}

//
// first deploy the FT, NFT, and market code
// tokenCode, _ := DownloadFile(FungibleTokenContractsBaseURL + FlowTokenFile)
// codeWithFTAddr := strings.ReplaceAll(string(tokenCode), "0x02", "0x04")
// tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
// tokenAddr, err := b.CreateAccount([]*flow.AccountKey{tokenAccountKey}, []byte(codeWithFTAddr))
// assert.Nil(t, err)
// _, err = b.CommitBlock()
// require.NoError(t, err)

// marketCode := ReadFile(MarketContractFile)
// marketAddr, err := b.CreateAccount(nil, marketCode)
// assert.Nil(t, err)
// _, err = b.CommitBlock()
// require.NoError(t, err)

// create two new accounts
// bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
// bastianAddress, err := b.CreateAccount([]*flow.AccountKey{bastianAccountKey}, nil)

// joshAccountKey, joshSigner := accountKeys.NewWithSigner()
// joshAddress, err := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

// t.Run("Should be able to create FTs and NFT collections in each accounts storage", func(t *testing.T) {
// 	// create Fungible tokens and NFTs in each accounts storage and store references
// 	setupUsersTokens(
// 		t, b, ftAddr, tokenAddr, nftAddr, topshotAddr,
// 		[]flow.Address{bastianAddress, joshAddress, topshotAddr},
// 		[]*flow.AccountKey{bastianAccountKey, joshAccountKey, topshotAccountKey},
// 		[]crypto.Signer{bastianSigner, joshSigner, topshotSigner},
// 	)

// 	tx := flow.NewTransaction().
// 		SetScript(fttest.GenerateMintTokensScript(ftAddr, tokenAddr, joshAddress, 30)).
// 		SetGasLimit(100).
// 		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
// 		SetPayer(b.RootKey().Address).
// 		AddAuthorizer(tokenAddr)

// 	SignAndSubmit(
// 		t, b, tx,
// 		[]flow.Address{b.RootKey().Address, tokenAddr},
// 		[]crypto.Signer{b.RootKey().Signer(), tokenSigner},
// 		false,
// 	)
// })

// t.Run("Can create sale collection", func(t *testing.T) {
// 	tx := flow.NewTransaction().
// 		SetScript(templates.GenerateCreateSaleScript(ftAddr, topshotAddr, marketAddr, 0.15)).
// 		SetGasLimit(100).
// 		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
// 		SetPayer(b.RootKey().Address).
// 		AddAuthorizer(bastianAddress)

// 	SignAndSubmit(
// 		t, b, tx,
// 		[]flow.Address{b.RootKey().Address, bastianAddress},
// 		[]crypto.Signer{b.RootKey().Signer(), bastianSigner},
// 		false,
// 	)

// 	result, err := b.ExecuteScript(templates.GenerateInspectSaleLenScript(marketAddr, bastianAddress, 0))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}
// })

// t.Run("Can put an NFT up for sale", func(t *testing.T) {
// 	tx := flow.NewTransaction().
// 		SetScript(GenerateStartSaleScript(nftAddr, marketAddr, 0, 10)).
// 		SetGasLimit(100).
// 		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
// 		SetPayer(b.RootKey().Address).
// 		AddAuthorizer(bastianAddress)

// 	SignAndSubmit(
// 		t, b, tx,
// 		[]flow.Address{b.RootKey().Address, bastianAddress},
// 		[]crypto.Signer{b.RootKey().Signer(), bastianSigner},
// 		false,
// 	)

// 	// Assert that the account's collection is correct
// 	result, err := b.ExecuteScript(GenerateInspectSaleScript(marketAddr, bastianAddress, 0, 10))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}

// 	result, err = b.ExecuteScript(GenerateInspectSaleLenScript(marketAddr, bastianAddress, 1))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}
// })

// t.Run("Cannot buy an NFT for less than the sale price", func(t *testing.T) {
// 	tx := flow.NewTransaction().
// 		SetScript(GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, bastianAddress, 0, 9)).
// 		SetGasLimit(100).
// 		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
// 		SetPayer(b.RootKey().Address).
// 		AddAuthorizer(joshAddress)

// 	SignAndSubmit(
// 		t, b, tx,
// 		[]flow.Address{b.RootKey().Address, joshAddress},
// 		[]crypto.Signer{b.RootKey().Signer(), joshSigner},
// 		true,
// 	)
// })

// t.Run("Cannot buy an NFT that is not for sale", func(t *testing.T) {
// 	tx := flow.NewTransaction().
// 		SetScript(GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, bastianAddress, 2, 10)).
// 		SetGasLimit(100).
// 		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
// 		SetPayer(b.RootKey().Address).
// 		AddAuthorizer(joshAddress)

// 	SignAndSubmit(
// 		t, b, tx,
// 		[]flow.Address{b.RootKey().Address, joshAddress},
// 		[]crypto.Signer{b.RootKey().Signer(), joshSigner},
// 		true,
// 	)
// })

// t.Run("Can buy an NFT that is for sale", func(t *testing.T) {
// 	tx := flow.NewTransaction().
// 		SetScript(GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, bastianAddress, 0, 10)).
// 		SetGasLimit(100).
// 		SetProposalKey(b.RootKey().Address, b.RootKey().ID, b.RootKey().SequenceNumber).
// 		SetPayer(b.RootKey().Address).
// 		AddAuthorizer(joshAddress)

// 	SignAndSubmit(
// 		t, b, tx,
// 		[]flow.Address{b.RootKey().Address, joshAddress},
// 		[]crypto.Signer{b.RootKey().Signer(), joshSigner},
// 		false,
// 	)

// 	result, err := b.ExecuteScript(GenerateInspectVaultScript(tokenAddr, bastianAddress, 8.5))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}

// 	result, err = b.ExecuteScript(GenerateInspectVaultScript(tokenAddr, joshAddress, 20))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}

// 	result, err = b.ExecuteScript(GenerateInspectVaultScript(tokenAddr, marketAddr, 1.5))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}

// 	// Assert that the accounts' collections are correct
// 	result, err = b.ExecuteScript(GenerateInspectCollectionLenScript(nftAddr, bastianAddress, 0))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}

// 	result, err = b.ExecuteScript(GenerateInspectCollectionScript(nftAddr, joshAddress, 0))
// 	require.NoError(t, err)
// 	if !assert.True(t, result.Succeeded()) {
// 		t.Log(result.Error.Error())
// 	}
// })

//}

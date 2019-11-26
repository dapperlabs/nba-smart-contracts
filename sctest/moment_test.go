package sctest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/sdk/keys"
)

const (
	MoldContractFile = "./contracts/topshot.cdc"
)

func TestMoldDeployment(t *testing.T) {
	b := newEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	tokenCode := ReadFile(MoldContractFile)
	_, err := b.CreateAccount(nil, tokenCode, GetNonce())
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	b.CommitBlock()
}

func TestCreateMoment(t *testing.T) {
	b := newEmulator()

	// First, deploy the contract
	tokenCode := ReadFile(MoldContractFile)
	contractAddr, err := b.CreateAccount(nil, tokenCode, GetNonce())
	assert.Nil(t, err)

	t.Run("Should be able to create a mold collection", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateMoldCollectionScript(contractAddr),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the account's collection is correct
		_, err = b.ExecuteScript(GenerateInspectMoldCollectionScript(contractAddr, b.RootAccountAddress(), 1, false))
		if !assert.Nil(t, err) {
			t.Log(err.Error())
		}
	})

	t.Run("Should be able to create molds within a collection", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateMintMoldScript(contractAddr, "\"KOBE\"", 2),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the account's collection is correct
		_, err = b.ExecuteScript(GenerateInspectMoldCollectionScript(contractAddr, b.RootAccountAddress(), 1, true))
		if !assert.Nil(t, err) {
			t.Log(err.Error())
		}

		tx = flow.Transaction{
			Script:         GenerateMintMoldScript(contractAddr, "\"Lebron\"", 1),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the account's collection is correct
		_, err = b.ExecuteScript(GenerateInspectMoldCollectionScript(contractAddr, b.RootAccountAddress(), 2, true))
		if !assert.Nil(t, err) {
			t.Log(err.Error())
		}
	})

	t.Run("Should be able to create a momentFactory that reference molds", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateMintMomentFactoryScript(contractAddr),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)
	})

	// create a new account
	bastianPrivateKey := randomKey()
	bastianPublicKey := bastianPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	bastianAddress, err := b.CreateAccount([]flow.AccountPublicKey{bastianPublicKey}, nil, GetNonce())

	t.Run("Should be able to create a momentCollection using the MomentFactory molds", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateMomentCollectionScript(contractAddr, b.RootAccountAddress()),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{bastianAddress},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey(), bastianPrivateKey}, []flow.Address{b.RootAccountAddress(), bastianAddress}, false)
	})

	t.Run("Should be able to mint a moment and deposit it into a user's collection", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateMintMomentScript(contractAddr, 2, "Uncommon", bastianAddress),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the account's collection is correct
		tx = flow.Transaction{
			Script:         GenerateInspectMomentScript(contractAddr, 1, "Lebron"),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{bastianAddress},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey(), bastianPrivateKey}, []flow.Address{b.RootAccountAddress(), bastianAddress}, false)
	})

	t.Run("Shouldn't be able to mint a moment that has reached the limit for its rarity", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateMintMomentScript(contractAddr, 2, "Uncommon", bastianAddress),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, true)

		// Assert that the account's collection is correct
		tx = flow.Transaction{
			Script:         GenerateInspectMomentScript(contractAddr, 1, "Lebron"),
			Nonce:          GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{bastianAddress},
		}

		SignAndSubmit(tx, b, t, []flow.AccountPrivateKey{b.RootKey(), bastianPrivateKey}, []flow.Address{b.RootAccountAddress(), bastianAddress}, false)
	})

}



// The main Fungible Token interface. Other token contracts will implement
// this interface
//
pub contract interface FungibleToken {

    // The total number of tokens in existence
    // it is up to the implementer to ensure that total supply 
    // stays accurate and up to date
    pub var totalSupply: UFix64

    // event that is emitted when the contract is created
    pub event FungibleTokenInitialized(initialSupply: UFix64)

    // event that is emitted when tokens are withdrawn from a Vault
    pub event Withdraw(amount: UFix64, from: Address?)

    // event that is emitted when tokens are deposited to a Vault
    pub event Deposit(amount: UFix64, to: Address?)

    // Provider
    // 
    // Interface that enforces the requirements for withdrawing
    // tokens from the implementing type.
    //
    // We don't enforce requirements on self.balance here because
    // it leaves open the possibility of creating custom providers
    // that don't necessarily need their own balance.
    //
    pub resource interface Provider {

        // withdraw
        //
        // Function that subtracts tokens from the owner's Vault
        // and returns a Vault resource (@Vault) with the removed tokens.
        //
        // The function's access level is public, but this isn't a problem
        // because even the public functions are not fully public at first.
        // anyone in the network can call them, but only if the owner grants
        // them access by publishing a resource that exposes the withdraw
        // function.
        //
        pub fun withdraw(amount: UFix64): @Vault {
            post {
                // `result` refers to the return value
                result.balance == amount:
                    "Withdrawal amount must be the same as the balance of the withdrawn Vault"
            }
        }
    }

    // Receiver 
    //
    // Interface that enforces the requirements for depositing
    // tokens into the implementing type
    //
    // We don't include a condition that checks the balance because
    // we want to give users the ability to make custom Receivers that
    // can do custom things with the tokens, like split them up and
    // send them to different places.
    //
    pub resource interface Receiver {

        // deposit
        //
        // Function that can be called to deposit tokens 
        // into the implementing resource type
        //
        pub fun deposit(from: @Vault) {
            pre {
                from.balance > UFix64(0):
                    "Deposit balance must be positive"
            }
        }
    }

    // Balance 
    
    // Interface that contains the balance field of the Vault
    // and enforces that when new Vault's are created, the balance
    // is initialized correctly.
    //
    pub resource interface Balance {

        // The total balance of the account's tokens
        pub var balance: UFix64

        init(balance: UFix64) {
            post {
                self.balance == balance: 
                    "Balance must be initialized to the initial balance"
            }
        }
    }

    // Vault
    // 
    // The resource that contains the functions to send and receive tokens.
    // 
    // The declaration of a concrete type in a contract interface means that
    // every Fungible Token contract that implements this interface
    // must define a concrete Vault object that
    // conforms to the Provider, Receiver, and Balace interfaces
    // and includes these fields and functions
    //
    pub resource Vault: Provider, Receiver, Balance {
        // The total balance of the accounts tokens
        pub var balance: UFix64

        // must declare init to conform to the Balance interface
        init(balance: UFix64)

        // withdraw subtracts `amount` from the vaults balance and
        // returns a vault object with the subtracted balance
        pub fun withdraw(amount: UFix64): @Vault {
            pre {
                self.balance >= amount: 
                    "Amount withdrawn must be less than or equal than the balance of the Vault"
            }
            post {
                // use the keywork before() to get the value of the 
                // specified field before function execution
                self.balance == before(self.balance) - amount:
                    "New Vault balance must be the difference of the previous balance and the withdrawn Vault"
            }
        }

        // deposit takes a vault object as a parameter and adds
        // its balance to the balance of the stored vault
        pub fun deposit(from: @Vault) {
            post {
                self.balance == before(self.balance) + before(from.balance):
                    "New Vault balance must be the sum of the previous balance and the deposited Vault"
            }
        }
    }

    // createEmptyVault
    // 
    // Any user can call this function to create a new Vault object
    // that has balance == 0
    //
    pub fun createEmptyVault(): @Vault {
        post {
            result.balance == UFix64(0): "The newly created Vault must have zero balance"
        }
    }
}

/**

# FlowToken example contract

Authors:

- Joshua Hannan <joshua.hannan@dapperlabs.com>
- Dieter Shirley <dete@dapperlabs.com>
- Bastian Mueller <bastian@dapperlabs.com>

It is not part of the standard, but just shows how most tokens
should implement the standard, including the Flow network token itself.

The FlowToken contract only needs to be deployed in one account.
The only part of the contract that would be stored in each user's account
is the Vault object, below

The implementation does not need to redefine the interfaces that are
already defined in the Fungible Token interface

*/

pub contract FlowToken: FungibleToken {

    // Total supply of flow tokens in existence
    pub var totalSupply: UFix64

    // Event that is emitted when the contract is created
    pub event FungibleTokenInitialized(initialSupply: UFix64)

    // Event that is emitted when tokens are withdrawn from a Vault
    pub event Withdraw(amount: UFix64, from: Address?)

    // Event that is emitted when tokens are deposited to a Vault
    pub event Deposit(amount: UFix64, to: Address?)

    // Event that is emitted when new tokens are minted
    pub event Mint(amount: UFix64)

    // Event that is emitted when tokens are destroyed
    pub event Burn(amount: UFix64)

    // Event that is emitted when a mew minter resource is created
    pub event MinterCreated(allowedAmount: UFix64)

    // Vault
    //
    // Each user stores an instance of only the Vault in their storage
    // The functions in the Vault and governed by the pre and post conditions
    // in FungibleToken when they are called.
    // The checks happen at runtime whenever a function is called.
    //
    // Resources can only be created in the context of the contract that they
    // are defined in, so there is no way for a malicious user to create Vaults
    // out of thin air. A special Minter resource needs to be defined to mint
    // new tokens.
    //
    pub resource Vault: FungibleToken.Provider, FungibleToken.Receiver, FungibleToken.Balance {

        // holds the balance of a users tokens
        pub var balance: UFix64

        // initialize the balance at resource creation time
        init(balance: UFix64) {
            self.balance = balance
        }

        // withdraw
        //
        // Function that takes an integer amount as an argument
        // and withdraws that amount from the Vault.
        // It creates a new temporary Vault that is used to hold
        // the money that is being transferred. It returns the newly
        // created Vault to the context that called so it can be deposited
        // elsewhere.
        //
        pub fun withdraw(amount: UFix64): @Vault {
            self.balance = self.balance - amount
            emit Withdraw(amount: amount, from: self.owner?.address)
            return <-create Vault(balance: amount)
        }

        // deposit
        //
        // Function that takes a Vault object as an argument and adds
        // its balance to the balance of the owners Vault.
        // It is allowed to destroy the sent Vault because the Vault
        // was a temporary holder of the tokens. The Vault's balance has
        // been consumed and therefore can be destroyed.
        pub fun deposit(from: @Vault) {
            self.balance = self.balance + from.balance
            emit Deposit(amount: from.balance, to: self.owner?.address)
            destroy from
        }
    }

    // createEmptyVault
    //
    // Function that creates a new Vault with a balance of zero
    // and returns it to the calling context. A user must call this function
    // and store the returned Vault in their storage in order to allow their
    // account to be able to receive deposits of this token type.
    //
    pub fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0.0)
    }

    // MintAndBurn
    //
    // Resource object that token admin accounts could hold
    // to mint and burn new tokens.
    //
    pub resource MintAndBurn {

        // the amount of tokens that the minter is allowed to mint
        pub var allowedAmount: UFix64

        // mintTokens
        //
        // Function that mints new tokens, adds them to the total Supply,
        // and returns them to the calling context
        //
        pub fun mintTokens(amount: UFix64): @Vault {
            pre {
                amount > UFix64(0): "Amount minted must be greater than zero"
                amount <= self.allowedAmount: "Amount minted must be less than the allowed amount"
            }
            FlowToken.totalSupply = FlowToken.totalSupply + amount
            self.allowedAmount = self.allowedAmount - amount
            emit Mint(amount: amount)
            return <-create Vault(balance: amount)
        }

        // burnTokens
        //
        // Function that takes a Vault as an argument, subtracts its balance
        // from the total supply, then destroys the Vault,
        // thereby removing the tokens from existence.
        //
        // Returns the amount that was burnt.
        //
        pub fun burnTokens(from: @Vault): UFix64 {
            let amount = from.balance
            FlowToken.totalSupply = FlowToken.totalSupply - from.balance
            destroy from
            emit Burn(amount: amount)
            return amount
        }

        // createNewMinter
        //
        // Function that creates and returns a new minter resource
        //
        pub fun createNewMinter(allowedAmount: UFix64): @MintAndBurn {
            emit MinterCreated(allowedAmount: allowedAmount)
            return <-create MintAndBurn(allowedAmount: allowedAmount)
        }

        init(allowedAmount: UFix64) {
            self.allowedAmount = allowedAmount
        }
    }

    // The initializer for the contract. All fields in the contract must
    // be initialized at deployment. This is just an example of what
    // an implementation could do in the initializer.
    //
    // The numbers are arbitrary.
    //
    init() {
        // Initialize the totalSupply field to the initial balance
        self.totalSupply = 1000.0

        // Create the Vault with the total supply of tokens and save it in storage
        //
        let vault <- create Vault(balance: self.totalSupply)
        self.account.save(<-vault, to: /storage/flowTokenVault)

        // Create a public capability to the stored Vault that only exposes
        // the `deposit` method through the `Receiver` interface
        //
        self.account.link<&{FungibleToken.Receiver}>(
            /public/flowTokenReceiver,
            target: /storage/flowTokenVault
        )

        // Create a public capability to the stored Vault that only exposes
        // the `balance` field through the `Balance` interface
        //
        self.account.link<&{FungibleToken.Balance}>(
            /public/flowTokenBalance,
            target: /storage/flowTokenVault
        )

        // Create a new MintAndBurn resource and store it in account storage
        let mintAndBurn <- create MintAndBurn(allowedAmount: 100.0)
        self.account.save(<-mintAndBurn, to: /storage/flowTokenMintAndBurn)

        // Emit an event that shows that the contract was initialized
        emit FungibleTokenInitialized(initialSupply: self.totalSupply)
    }
}


 

 


// The Fungible Token standard interface that all Fungible Tokens
// would have to conform to

access(all) contract interface FungibleToken {

    // The total number of tokens in existence
    // it is up to the implementor to ensure that total supply 
    // stays accurate and up to date
    access(all) var totalSupply: UInt64

    // event that is emitted when the contract is created
    access(all) event FungibleTokenInitialized(initialSupply: UInt64)

    // event that is emitted when tokens are withdrawn from a Vault
    access(all) event Withdraw(amount: UInt64)

    // event that is emitted when tokens are deposited to a Vault
    access(all) event Deposit(amount: UInt64)

    // Interface that enforces the requirements for withdrawing
    // tokens from the implementing type
    //
    access(all) resource interface Provider {
        pub fun withdraw(amount: UInt64): @Vault {
            post {
                result.balance == amount:
                    "Withdrawal amount must be the same as the balance of the withdrawn Vault"
            }
        }
    }

    // Interface that enforces the requirements for depositing
    // tokens into the implementing type
    //
    // We don't include a condition that checks the balance because
    // we want to give users the ability to make custom Receivers that
    // can do custom things with the tokens
    access(all) resource interface Receiver {
        pub fun deposit(from: @Vault) {
            pre {
                from.balance > UInt64(0):
                    "Deposit balance must be positive"
            }
        }
    }

    // Interface that contains the balance field of the Vault
    //
    access(all) resource interface Balance {
        pub var balance: UInt64
    }

    // The declaration of a concrete type in a contract interface means that
    // every Fungible Token contract that implements this interface
    // must define a concrete Vault object that
    // conforms to the Provider, Receiver, and Balace interfaces
    // and includes these fields and functions
    //
    access(all) resource Vault: Provider, Receiver, Balance {
        // keeps track of the total balance of the accounts tokens
        access(all) var balance: UInt64

        init(balance: UInt64) {
            post {
                self.balance == balance: "Balance must be initialized to the initial balance"
                // cannot get interface fields from within resource
                //self.balance <= self.totalSupply: "Balance must be less than total supply"
            }
        }

        // withdraw subtracts `amount` from the vaults balance and
        // returns a vault object with the subtracted balance
        access(all) fun withdraw(amount: UInt64): @Vault {
            pre {
                self.balance >= amount: "Amount withdrawn must be less than the balance of the Vault!"
            }
        }

        // deposit takes a vault object as a parameter and adds
        // its balance to the balance of the stored vault, then
        // destroys the sent vault because its balance has been consumed
        access(all) fun deposit(from: @Vault) {
            post {
                self.balance == before(self.balance) + before(from.balance):
                    "New Vault balance must be the sum of the previous balance and the deposited Vault"
            }
        }
    }

    // Any user can call this function to create a new Vault object
    // that has balance = 0
    //
    access(all) fun createEmptyVault(): @Vault {
        post {
            result.balance == UInt64(0): "The newly created Vault must have zero balance"
        }
    }
}


// This is an Example Implementation of the Fungible Token Standard
// It is not part of the standard, but just shows how most tokens would implememt the standard
//
access(all) contract FlowToken: FungibleToken {

    access(all) var totalSupply: UInt64

    // event that is emitted when the contract is created
    access(all) event FungibleTokenInitialized(initialSupply: UInt64)

    // event that is emitted when tokens are withdrawn from a Vault
    access(all) event Withdraw(amount: UInt64)

    // event that is emitted when tokens are deposited to a Vault
    access(all) event Deposit(amount: UInt64)

    access(all) resource Vault: FungibleToken.Provider, FungibleToken.Receiver, FungibleToken.Balance {
        
        access(all) var balance: UInt64

        init(balance: UInt64) {
            self.balance = balance
        }

        access(all) fun withdraw(amount: UInt64): @Vault {
            self.balance = self.balance - amount
            emit Withdraw(amount: amount)
            return <-create Vault(balance: amount)
        }
        
        access(all) fun deposit(from: @Vault) {
            self.balance = self.balance + from.balance
            emit Deposit(amount: from.balance)
            destroy from
        }
    }

    access(all) fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0)
    }

    // This function is included here purely for testing purposes and would not be
    // included in an actual implementation
    //
    access(all) fun createVault(initialBalance: UInt64): @Vault {
        return <-create Vault(balance: initialBalance)
    }

    init() {
        self.totalSupply = 1000

        // create the Vault with the initial balance and put it in storage
        let oldVault <- self.account.storage[Vault] <- create Vault(balance: 1000)
        destroy oldVault

        // Create a private reference to the Vault that has all the fields and methods
        self.account.storage[&Vault] = &self.account.storage[Vault] as &Vault

        // Create a public reference to the Vault that only exposes the deposit method
        self.account.published[&FungibleToken.Receiver] = &self.account.storage[Vault] as &FungibleToken.Receiver
        self.account.published[&FungibleToken.Balance] = &self.account.storage[Vault] as &FungibleToken.Balance

        emit FungibleTokenInitialized(initialSupply: self.totalSupply)
    }
}
 
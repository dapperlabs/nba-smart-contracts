

// The main Fungible Token interface. Other token contracts will implement
// this interface
//
pub contract interface FungibleToken {

    // The total number of tokens in existence
    // it is up to the implementer to ensure that total supply 
    // stays accurate and up to date
    pub var totalSupply: UInt64

    // event that is emitted when the contract is created
    pub event FungibleTokenInitialized(initialSupply: UInt64)

    // event that is emitted when tokens are withdrawn from a Vault
    pub event Withdraw(amount: UInt64)

    // event that is emitted when tokens are deposited to a Vault
    pub event Deposit(amount: UInt64)

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
        pub fun withdraw(amount: UInt64): @Vault {
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
                from.balance > UInt64(0):
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
        pub var balance: UInt64

        init(balance: UInt64) {
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
        pub var balance: UInt64

        // must declare init to conform to the Balance interface
        init(balance: UInt64)

        // withdraw subtracts `amount` from the vaults balance and
        // returns a vault object with the subtracted balance
        pub fun withdraw(amount: UInt64): @Vault {
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
            result.balance == UInt64(0): "The newly created Vault must have zero balance"
        }
    }
}


// This is an Example Implementation of the Fungible Token Standard
// It is not part of the standard, but just shows how most tokens would implememt the standard
//
pub contract FlowToken: FungibleToken {

    pub var totalSupply: UInt64

    // event that is emitted when the contract is created
    pub event FungibleTokenInitialized(initialSupply: UInt64)

    // event that is emitted when tokens are withdrawn from a Vault
    pub event Withdraw(amount: UInt64)

    // event that is emitted when tokens are deposited to a Vault
    pub event Deposit(amount: UInt64)

    pub resource Vault: FungibleToken.Provider, FungibleToken.Receiver, FungibleToken.Balance {
        
        pub var balance: UInt64

        init(balance: UInt64) {
            self.balance = balance
        }

        pub fun withdraw(amount: UInt64): @Vault {
            self.balance = self.balance - amount
            emit Withdraw(amount: amount)
            return <-create Vault(balance: amount)
        }
        
        pub fun deposit(from: @Vault) {
            self.balance = self.balance + from.balance
            emit Deposit(amount: from.balance)
            destroy from
        }
    }

    pub fun createEmptyVault(): @Vault {
        return <-create Vault(balance: 0)
    }

    // This function is included here purely for testing purposes and would not be
    // included in an actual implementation
    //
    pub fun createVault(initialBalance: UInt64): @Vault {
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
 
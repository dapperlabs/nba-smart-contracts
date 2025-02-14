import "FungibleToken"
import "FlowToken"

import "EVM"

/// Transfers $FLOW from the signer's account Cadence Flow balance to the recipient's hex-encoded EVM address.
///
transaction(recipientEVMAddressHex: String, amount: UFix64) {

    var sentVault: @FlowToken.Vault
    let recipientEVMAddress: EVM.EVMAddress
    let recipientPreBalance: UFix64

    prepare(signer: auth(BorrowValue, SaveValue) &Account) {
        // Borrow a reference to the signer's FlowToken.Vault and withdraw the amount
        let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(
                from: /storage/flowTokenVault
            ) ?? panic("Could not borrow reference to the owner's Vault!")
        self.sentVault <- vaultRef.withdraw(amount: amount) as! @FlowToken.Vault

        // Get the recipient's EVM address
        self.recipientEVMAddress = EVM.addressFromString(recipientEVMAddressHex)

        // Get the recipient's balance before the transfer to check the amount transferred
        self.recipientPreBalance = self.recipientEVMAddress.balance().inFLOW()
    }

    execute {
        // Deposit the amount to the recipient's EVM address
        self.recipientEVMAddress.deposit(from: <-self.sentVault)
    }

    post {
        self.recipientEVMAddress.balance().inFLOW() == self.recipientPreBalance + amount:
            "Problem transferring value to EVM address"
    }
}

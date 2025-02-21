import Crypto
import "FungibleToken"
import "FlowToken"
import "EVM"

/// Creates a new flow account with single full-weight key, initial FLOW tokens funding, and
/// a COA EVM account in storage.
///
/// @param pubKey: String - public key to be added to the account
/// @param amount: UFix64 - amount of FLOW tokens to be transferred to the new account
///
transaction(pubKey: String, amount: UFix64) {
    let sentVault: @FlowToken.Vault
    let receiverRef: &{FungibleToken.Receiver}

    prepare(signer: auth(BorrowValue) &Account) {
        // Create new account
        let account = Account(payer: signer)

        // Add public key with full weight, SHA2_256, and ECDSA_P256
        account.keys.add(
            publicKey: PublicKey(
                publicKey: pubKey.decodeHex(),
                signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
            ),
            hashAlgorithm: HashAlgorithm.SHA2_256,
            weight: 1000.0
        )

        // Set COA paths
        let storagePath = StoragePath(identifier: "evm")!
        let publicPath = PublicPath(identifier: "evm")!

        // Create and save new COA in new account's storage
        account.storage.save<@EVM.CadenceOwnedAccount>(<- EVM.createCadenceOwnedAccount(), to: storagePath)

        // Issue and publish capability to the COA
        account.capabilities.unpublish(publicPath)
        account.capabilities.publish(account.capabilities.storage.issue<&EVM.CadenceOwnedAccount>(storagePath), at: publicPath)

        // Borrow reference to the signer's FLOW token vault and withdraw provided amount
        let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(
                from: /storage/flowTokenVault
            ) ?? panic("Could not borrow reference to the owner's Vault!")
        self.sentVault <- vaultRef.withdraw(amount: amount) as! @FlowToken.Vault

        // Borrow reference to the new account's FLOW tokens receiver
        self.receiverRef = account.capabilities.borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)
            ?? panic("Could not borrow a Receiver reference to the FlowToken Vault in account ")
    }

    execute {
        // Deposit the withdrawn tokens in the new account's receiver
        self.receiverRef.deposit(from: <- self.sentVault)
    }
}

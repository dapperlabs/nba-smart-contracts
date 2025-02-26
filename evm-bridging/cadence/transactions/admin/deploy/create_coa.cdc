import "Burner"
import "FungibleToken"
import "FlowToken"
import "EVM"

/// Creates a COA and saves it in the signer's Flow account & passing the given value of Flow into FlowEVM
///
transaction(amount: UFix64) {
    let fundingVault: @FlowToken.Vault?
    let coa: &EVM.CadenceOwnedAccount

    prepare(signer: auth(BorrowValue, SaveValue, IssueStorageCapabilityController, PublishCapability) &Account) {
        /* --- Configure COA --- */
        //
        // Ensure there is not yet a CadenceOwnedAccount in the standard path
        let coaPath = /storage/evm
        if signer.storage.type(at: coaPath) != nil {
            panic(
                "Object already exists in signer's account at path=".concat(coaPath.toString())
                .concat(". Make sure the signing account does not already have a CadenceOwnedAccount.")
            )
        }
        // COA not found in standard path, create and publish a public **unentitled** capability
        signer.storage.save(<-EVM.createCadenceOwnedAccount(), to: coaPath)
        let coaCapability = signer.capabilities.storage.issue<&EVM.CadenceOwnedAccount>(coaPath)
        signer.capabilities.publish(coaCapability, at: /public/evm)

        // Borrow the CadenceOwnedAccount reference
        self.coa = signer.storage.borrow<&EVM.CadenceOwnedAccount>(
                from: coaPath
            ) ?? panic(
                "Could not find CadenceOwnedAccount (COA) in signer's account at path=".concat(coaPath.toString())
                .concat(". Make sure the signing account has initialized a COA at the expected path.")
            )

        /* --- Assign fundingVault --- */
        //
        if amount > 0.0 {
            // Reference the signer's FLOW vault & withdraw the funding amount
            let vault = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)
                ?? panic("Could not borrow a reference to the signer's FLOW vault from path=/storage/flowTokenVault")
            self.fundingVault <- vault.withdraw(amount: amount) as! @FlowToken.Vault
        } else {
            // No funding requested, so no need to withdraw from the vault
            self.fundingVault <- nil
        }
    }

    pre {
        self.fundingVault == nil || self.fundingVault?.balance ?? 0.0 == amount:
            "Mismatched funding vault acquired given requested amount=".concat(amount.toString())
    }

    execute {
        // Fund if necessary
        if self.fundingVault != nil || self.fundingVault?.balance ?? 0.0 > 0.0 {
            self.coa.deposit(from: <-self.fundingVault!)
        } else {
            Burner.burn(<-self.fundingVault)
        }
    }
}

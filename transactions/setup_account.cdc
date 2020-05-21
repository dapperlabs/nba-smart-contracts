import TopShot from 0x03
import FungibleToken from 0x04
import FlowToken from 0x05

// This is the transaction you would run from
// any account to set it up to use the fungible token and topshot Collection

transaction {

    prepare(acct: AuthAccount) {
        if acct.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault) == nil {
            let vault <- FlowToken.createEmptyVault() as! @FlowToken.Vault
            acct.save(<-vault, to: /storage/flowTokenVault)

            // Create a public capability to the stored Vault that only exposes
            // the `deposit` method through the `Receiver` interface
            //
            acct.link<&{FungibleToken.Receiver}>(
                /public/flowTokenReceiver,
                target: /storage/flowTokenVault
            )

            // Create a public capability to the stored Vault that only exposes
            // the `balance` field through the `Balance` interface
            //
            acct.link<&{FungibleToken.Balance}>(
                /public/flowTokenBalance,
                target: /storage/flowTokenVault
            )
        }

        if acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection) == nil {
            let collection <- TopShot.createEmptyCollection() as! @TopShot.Collection
            // Put a new Collection in storage
            acct.save(<-collection, to: /storage/MomentCollection)

            // create a public capability for the collection
            acct.link<&{TopShot.MomentCollectionPublic}>(/public/MomentCollection, target: /storage/MomentCollection)
        }
    }
}
 

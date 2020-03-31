import TopShot from 0x03
import FungibleToken, FlowToken from 0x01

// This is the transaction you would run from
// any account to set it up to use the fungible token and topshot Collection

transaction {

    prepare(acct: AuthAccount) {
        if acct.storage[FlowToken.Vault] == nil {
            let vault <- FlowToken.createEmptyVault()
            let oldVault <- acct.storage[FlowToken.Vault] <- vault
            destroy oldVault

            acct.published[&FlowToken.Vault{FungibleToken.Receiver}] = &acct.storage[FlowToken.Vault] as &FlowToken.Vault{FungibleToken.Receiver}
            acct.published[&FlowToken.Vault{FungibleToken.Balance}] = &acct.storage[FlowToken.Vault] as &FlowToken.Vault{FungibleToken.Balance}
        }

        if acct.storage[TopShot.Collection] == nil {
            let collection <- TopShot.createEmptyCollection()
            let oldCollection <- acct.storage[TopShot.Collection] <- collection
            destroy oldCollection

            acct.published[&TopShot.Collection{TopShot.MomentCollectionPublic}] = &acct.storage[TopShot.Collection] as &TopShot.Collection{TopShot.MomentCollectionPublic}
        }
    }
}
 
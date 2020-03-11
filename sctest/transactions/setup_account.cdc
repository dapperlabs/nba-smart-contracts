import TopShot from 0x02
import FungibleToken, FlowToken from 0x01

// This is the transaction you would run from
// any account to set it up to use the fungible token,
// topshot, and topshot market
transaction {

    prepare(acct: Account) {
        if acct.storage[FlowToken.Vault] == nil {
            let vault <- FlowToken.createEmptyVault()
            let oldVault <- acct.storage[FlowToken.Vault] <- vault
            destroy oldVault

            acct.published[&FungibleToken.Receiver] = &acct.storage[FlowToken.Vault] as &FungibleToken.Receiver
            acct.published[&FungibleToken.Balance] = &acct.storage[FlowToken.Vault] as &FungibleToken.Balance
        }

        if acct.storage[TopShot.Collection] == nil {
            let collection <- TopShot.createEmptyCollection()
            let oldCollection <- acct.storage[TopShot.Collection] <- collection
            destroy oldCollection

            acct.published[&TopShot.MomentCollectionPublic] = &acct.storage[TopShot.Collection] as &TopShot.MomentCollectionPublic
        }
    }
}
 
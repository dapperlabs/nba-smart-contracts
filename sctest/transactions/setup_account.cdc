import TopShot from 0x02
import Market from 0x03
import FungibleToken, FlowToken from 0x01

transaction {

    prepare(acct: Account) {
        if acct.storage[FlowToken.Vault] == nil {
            let vault <- FlowToken.createEmptyVault()
            let oldVault <- acct.storage[FlowToken.Vault] <- vault
            destroy oldVault

            acct.published[&FungibleToken.Receiver] = &acct.storage[FlowToken.Vault] as FungibleToken.Receiver
        }

        if acct.storage[TopShot.Collection] == nil {
            let collection <- TopShot.createEmptyCollection()
            let oldCollection <- acct.storage[TopShot.Collection] <- collection
            destroy oldCollection

            acct.published[&TopShot.MomentCollectionPublic] = &acct.storage[TopShot.Collection] as TopShot.MomentCollectionPublic
        }

        if acct.storage[Market.SaleCollection] == nil {
            let receiverRef = acct.published[&FungibleToken.Receiver] ?? panic("No receiver ref!")

            let sale <- Market.createSaleCollection(ownerVault: receiverRef)
            let oldSale <- acct.storage[Market.SaleCollection] <- sale
            destroy oldSale

            acct.published[&Market.SalePublic] = &acct.storage[Market.SaleCollection] as Market.SalePublic
        }
    }
}
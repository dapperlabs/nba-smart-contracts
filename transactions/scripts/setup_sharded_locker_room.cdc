import TopShot from 0xTOPSHOTADDRESS
import TopShotShardedCollection from 0xSHARDEDADDRESS
import NonFungibleToken from 0xNFTADDRESS

/*
    This transaction creates a capability on the minter account
    that links it to the locker room account.

    For example, on testnet:
    minter account = 0x70dff4d1005824db
    locker account = 0xd80d84b4b0a88782
*/


transaction() {

    prepare(minter: auth(Storage, Capabilities) &Account, locker: auth(Storage, Capabilities) &Account) {

        minter.storage.save(
            locker.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &TopShotShardedCollection.ShardedCollection>(/storage/TopShotShardedCollection),
            to: /storage/lockerTSShardedCollection2
        )

        minter.storage.save(
            locker.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &TopShot.Collection>(/storage/MomentCollection),
            to: /storage/lockerTSCollection2
        )        
    }
}
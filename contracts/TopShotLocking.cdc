import NonFungibleToken from 0xNFTADDRESS

access(all) contract TopShotLocking {

    // -----------------------------------------------------------------------
    // TopShotLocking contract Events
    // -----------------------------------------------------------------------

    // Emitted when a Moment is locked
    access(all) event MomentLocked(id: UInt64, duration: UFix64, expiryTimestamp: UFix64)

    // Emitted when a Moment is unlocked
    access(all) event MomentUnlocked(id: UInt64)

    // Dictionary of locked NFTs
    // TopShot nft resource id is the key
    // locked until timestamp is the value
    access(self) var lockedNFTs: {UInt64: UFix64}

    // Dictionary of NFTs overridden to be unlocked
    access(self) var unlockableNFTs: {UInt64: Bool} // nft resource id is the key

    // isLocked Returns a boolean indicating if an nft exists in the lockedNFTs dictionary
    //
    // Parameters: nftRef: A reference to the NFT resource
    //
    // Returns: true if NFT is locked
    access(all) view fun isLocked(nftRef: &{NonFungibleToken.NFT}): Bool {
        return self.lockedNFTs.containsKey(nftRef.id)
    }

    // getLockExpiry Returns the unix timestamp when an nft is unlockable
    //
    // Parameters: nftRef: A reference to the NFT resource
    //
    // Returns: unix timestamp
    access(all) view fun getLockExpiry(nftRef: &{NonFungibleToken.NFT}): UFix64 {
        if !self.lockedNFTs.containsKey(nftRef.id) {
            panic("NFT is not locked")
        }
        return self.lockedNFTs[nftRef.id]!
    }

    // lockNFT Takes an NFT resource and adds its unique identifier to the lockedNFTs dictionary
    //
    // Parameters: nft: NFT resource
    //             duration: number of seconds the NFT will be locked for
    //
    // Returns: the NFT resource
    access(all) fun lockNFT(nft: @{NonFungibleToken.NFT}, duration: UFix64): @{NonFungibleToken.NFT} {
        let TopShotNFTType: Type = CompositeType("A.TOPSHOTADDRESS.TopShot.NFT")!
        if !nft.isInstance(TopShotNFTType) {
            panic("NFT is not a TopShot NFT")
        }

        if self.lockedNFTs.containsKey(nft.id) {
            // already locked - short circuit and return the nft
            return <- nft
        }

        let expiryTimestamp = getCurrentBlock().timestamp + duration

        self.lockedNFTs[nft.id] = expiryTimestamp

        emit MomentLocked(id: nft.id, duration: duration, expiryTimestamp: expiryTimestamp)

        return <- nft
    }

    // unlockNFT Takes an NFT resource and removes it from the lockedNFTs dictionary
    //
    // Parameters: nft: NFT resource
    //
    // Returns: the NFT resource
    //
    // NFT must be eligible for unlocking by an admin
    access(all) fun unlockNFT(nft: @{NonFungibleToken.NFT}): @{NonFungibleToken.NFT} {
        if !self.lockedNFTs.containsKey(nft.id) {
            // nft is not locked, short circuit and return the nft
            return <- nft
        }

        let lockExpiryTimestamp: UFix64 = self.lockedNFTs[nft.id]!
        let isPastExpiry: Bool = getCurrentBlock().timestamp >= lockExpiryTimestamp

        let isUnlockableOverridden: Bool = self.unlockableNFTs.containsKey(nft.id)

        if !(isPastExpiry || isUnlockableOverridden) {
            panic("NFT is not eligible to be unlocked, expires at ".concat(lockExpiryTimestamp.toString()))
        }

        self.unlockableNFTs.remove(key: nft.id)
        self.lockedNFTs.remove(key: nft.id)

        emit MomentUnlocked(id: nft.id)

        return <- nft
    }

    // getIDs Returns the ids of all locked Top Shot NFT tokens
    //
    // Returns: array of ids
    //
    access(all) view fun getIDs(): [UInt64] {
        return self.lockedNFTs.keys
    }

    // getExpiry Returns the timestamp when a locked token is eligible for unlock
    //
    // Parameters: tokenID: the nft id of the locked token
    //
    // Returns: a unix timestamp in seconds
    //
    access(all) view fun getExpiry(tokenID: UInt64): UFix64? {
        return self.lockedNFTs[tokenID]
    }

    // getLockedNFTsLength Returns the count of locked tokens
    //
    // Returns: an integer containing the number of locked tokens
    //
    access(all) view fun getLockedNFTsLength(): Int {
        return self.lockedNFTs.length
    }

    // The path to the TopShotLocking Admin resource belonging to the Account
    // which the contract is deployed on
    access(all) view fun AdminStoragePath() : StoragePath { return /storage/TopShotLockingAdmin}

    // Admin is a special authorization resource that 
    // allows the owner to override the lock on a moment
    //
    access(all) resource Admin {
        // createNewAdmin creates a new Admin resource
        //
        access(all) fun createNewAdmin(): @Admin {
            return <-create Admin()
        }

        // markNFTUnlockable marks a given nft as being
        // unlockable, overridding the expiry timestamp
        // the nft owner will still need to send an unlock transaction to unlock
        //
        access(all) fun markNFTUnlockable(nftRef: &{NonFungibleToken.NFT}) {
            TopShotLocking.unlockableNFTs[nftRef.id] = true
        }

        access(all) fun unlockByID(id: UInt64) {
            if !TopShotLocking.lockedNFTs.containsKey(id) {
                // nft is not locked, do nothing
                return
            }
            TopShotLocking.lockedNFTs.remove(key: id)
            emit MomentUnlocked(id: id)
        }

        // admin may alter the expiry of a lock on an NFT
        access(all) fun setLockExpiryByID(id: UInt64, expiryTimestamp: UFix64) {
            if expiryTimestamp < getCurrentBlock().timestamp {
                panic("cannot set expiry in the past")
            }

            let duration = expiryTimestamp - getCurrentBlock().timestamp

            TopShotLocking.lockedNFTs[id] = expiryTimestamp

            emit MomentLocked(id: id, duration: duration, expiryTimestamp: expiryTimestamp)
        }

        // unlocks all NFTs
        access(all) fun unlockAll() {
            TopShotLocking.lockedNFTs = {}
            TopShotLocking.unlockableNFTs = {}
        }
    }

    // -----------------------------------------------------------------------
    // TopShotLocking initialization function
    // -----------------------------------------------------------------------
    //
    init() {
        self.lockedNFTs = {}
        self.unlockableNFTs = {}

        // Create a single admin resource
        let admin <- create Admin()

        // Store it in private account storage in `init` so only the admin can use it
        self.account.storage.save(<-admin, to: TopShotLocking.AdminStoragePath())
    }
}

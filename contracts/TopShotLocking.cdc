import NonFungibleToken from 0xNFTADDRESS

pub contract TopShotLocking {

    // -----------------------------------------------------------------------
    // TopShotLocking contract Events
    // -----------------------------------------------------------------------

    // Emitted when a Moment is locked
    pub event MomentLocked(id: UInt64, duration: UFix64, expiryTimestamp: UFix64)

    // Emitted when a Moment is unlocked
    pub event MomentUnlocked(id: UInt64)

    // Dictionary of locked NFTs
    // nft resource uuid is the key
    // locked until timestamp is the value
    access(self) var lockedNFTs: {UInt64: UFix64}

    // Dictionary of NFTs overridden to be unlocked
    access(self) var unlockableNFTs: {UInt64: Bool} // nft resource uuid is the key

    // isLocked Returns a boolean indicating if an nft exists in the lockedNFTs dictionary
    //
    // Parameters: nftRef: A reference to the NFT resource
    //
    // Returns: true if NFT is locked
    pub fun isLocked(nftRef: &NonFungibleToken.NFT): Bool {
        return self.lockedNFTs.containsKey(nftRef.uuid)
    }

    // getLockExpiry Returns the unix timestamp when an nft is unlockable
    //
    // Parameters: nftRef: A reference to the NFT resource
    //
    // Returns: unix timestamp
    pub fun getLockExpiry(nftRef: &NonFungibleToken.NFT): UFix64 {
        if !self.lockedNFTs.containsKey(nftRef.uuid) {
            panic("NFT is not locked")
        }
        return self.lockedNFTs[nftRef.uuid]!
    }

    // lockNFT Takes an NFT resource and adds its unique identifier to the lockedNFTs dictionary
    //
    // Parameters: nft: NFT resource
    //             duration: number of seconds the NFT will be locked for
    //
    // Returns: the NFT resource
    pub fun lockNFT(nft: @NonFungibleToken.NFT, duration: UFix64): @NonFungibleToken.NFT {
        let TopShotNFTType: Type = CompositeType("A.TOPSHOTADDRESS.TopShot.NFT")!
        if !nft.isInstance(TopShotNFTType) {
            panic("NFT is not a TopShot NFT")
        }

        let uuid = nft.uuid
        if self.lockedNFTs.containsKey(uuid) {
            // already locked - short circuit and return the nft
            return <- nft
        }

        let expiryTimestamp = getCurrentBlock().timestamp + duration

        self.lockedNFTs[uuid] = expiryTimestamp

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
    pub fun unlockNFT(nft: @NonFungibleToken.NFT): @NonFungibleToken.NFT {
        let uuid = nft.uuid
        if !self.lockedNFTs.containsKey(uuid) {
            // nft is not locked, short circuit and return the nft
            return <- nft
        }

        let lockExpiryTimestamp: UFix64 = self.lockedNFTs[uuid]!
        let isPastExpiry: Bool = getCurrentBlock().timestamp >= lockExpiryTimestamp

        let isUnlockableOverridden: Bool = self.unlockableNFTs.containsKey(uuid)

        if !(isPastExpiry || isUnlockableOverridden) {
            panic("NFT is not eligible to be unlocked, expires at ".concat(lockExpiryTimestamp.toString()))
        }

        self.unlockableNFTs.remove(key: uuid)
        self.lockedNFTs.remove(key: uuid)

        emit MomentUnlocked(id: nft.id)

        return <- nft
    }

    // Admin is a special authorization resource that 
    // allows the owner to override the lock on a moment
    //
    pub resource Admin {
        // createNewAdmin creates a new Admin resource
        //
        pub fun createNewAdmin(): @Admin {
            return <-create Admin()
        }

        // markNFTUnlockable marks a given nft as being
        // unlockable, overridding the expiry timestamp
        // the nft owner will still need to send an unlock transaction to unlock
        //
        pub fun markNFTUnlockable(nftRef: &NonFungibleToken.NFT) {
            TopShotLocking.unlockableNFTs[nftRef.uuid] = true
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
        self.account.save(<-admin, to: /storage/TopShotLockingAdmin)
    }
}

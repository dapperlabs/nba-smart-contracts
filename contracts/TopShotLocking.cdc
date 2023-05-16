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
    pub fun isLocked(nftRef: &NonFungibleToken.NFT): Bool {
        return self.lockedNFTs.containsKey(nftRef.id)
    }

    // getLockExpiry Returns the unix timestamp when an nft is unlockable
    //
    // Parameters: nftRef: A reference to the NFT resource
    //
    // Returns: unix timestamp
    pub fun getLockExpiry(nftRef: &NonFungibleToken.NFT): UFix64 {
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
    pub fun lockNFT(nft: @NonFungibleToken.NFT, duration: UFix64): @NonFungibleToken.NFT {
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
    pub fun unlockNFT(nft: @NonFungibleToken.NFT): @NonFungibleToken.NFT {
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
    pub fun getIDs(): [UInt64] {
        return self.lockedNFTs.keys
    }

    // getExpiry Returns the timestamp when a locked token is eligible for unlock
    //
    // Parameters: tokenID: the nft id of the locked token
    //
    // Returns: a unix timestamp in seconds
    //
    pub fun getExpiry(tokenID: UInt64): UFix64? {
        return self.lockedNFTs[tokenID]
    }

    // getLockedNFTsLength Returns the count of locked tokens
    //
    // Returns: an integer containing the number of locked tokens
    //
    pub fun getLockedNFTsLength(): Int {
        return self.lockedNFTs.length
    }

    // The path to the Subedition Admin resource belonging to the Account
    // which the contract is deployed on
    pub fun AdminStoragePath() : StoragePath { return /storage/TopShotLockingAdmin}

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
            TopShotLocking.unlockableNFTs[nftRef.id] = true
        }

        pub fun unlockByID(id: UInt64) {
            if !TopShotLocking.lockedNFTs.containsKey(id) {
                // nft is not locked, do nothing
                return
            }
            TopShotLocking.lockedNFTs.remove(key: id)
            emit MomentUnlocked(id: id)
        }

        // unlocks all NFTs
        pub fun unlockAll() {
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
        self.account.save(<-admin, to: TopShotLocking.AdminStoragePath())
    }
}

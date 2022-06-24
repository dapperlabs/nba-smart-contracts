import NonFungibleToken from 0xNFTADDRESS

pub contract TopShotLocking {

    // -----------------------------------------------------------------------
    // TopShotLocking contract Events
    // -----------------------------------------------------------------------

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

    // lockNFT Takes an NFT resource and adds its unique identifier to the lockedNFTs dictionary
    //
    // Parameters: nft: NFT resource
    //             expiryTimestamp: The unix timestamp in seconds after which the nft may be unlocked
    //
    // Returns: the NFT resource
    pub fun lockNFT(nft: @NonFungibleToken.NFT, expiryTimestamp: UFix64): @NonFungibleToken.NFT {
        let TopShotNFTType: Type = CompositeType("A.0xTOPSHOTADDRESS.TopShot.NFT")!
        if !nft.isInstance(TopShotNFTType) {
            panic("NFT is not a TopShot NFT")
        }

        let id = nft.uuid
        if self.lockedNFTs.containsKey(id) {
            // already locked - short circuit and return the nft
            return <- nft
        }

        self.lockedNFTs[id] = expiryTimestamp

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
        let id = nft.uuid
        if !self.lockedNFTs.containsKey(id) {
            // nft is not locked, short circuit and return the nft
            return <- nft
        }

        let lockExpiryTimestamp: UFix64 = self.lockedNFTs[id]!
        let isPastExpiry: Bool = getCurrentBlock().timestamp >= lockExpiryTimestamp

        let isUnlockableOverridden: Bool = self.unlockableNFTs.containsKey(id)

        if !(isPastExpiry || isUnlockableOverridden) {
            panic("NFT is not eligible to be unlocked, expires at ".concat(lockExpiryTimestamp.toString()))
        }

        self.unlockableNFTs.remove(key: id)
        self.lockedNFTs.remove(key: id)

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

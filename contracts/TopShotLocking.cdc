import NonFungibleToken from 0xNFTADDRESS

pub contract TopShotLocking {

    // -----------------------------------------------------------------------
    // TopShotLocking contract Events
    // -----------------------------------------------------------------------

    // Emitted when a moment is withdrawn from a Collection
    pub event Locked(id: UInt64, expiryTimestamp: UFix64)
    // Emitted when a moment is deposited into a Collection
    pub event Unlocked(id: UInt64)

    // Dictionary of locked NFTs
    // nft resource uuid is the key
    // locked until timestamp is the value
    pub var lockedNFTs: {UInt64: UFix64} 

    // Dictionary of NFTs overridden to be unlocked
    pub var unlockableNFTs: {UInt64: Bool} // nft resource uuid is the key

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
    //
    // Returns: the NFT resource
    pub fun lockNFT(nft: @NonFungibleToken.NFT, expiryTimestamp: UFix64): @NonFungibleToken.NFT {
        let id = nft.uuid

        if self.lockedNFTs.containsKey(id) {
            panic("NFT is already locked")
        }
        
        self.lockedNFTs[id] = expiryTimestamp

        emit Locked(id: id, expiryTimestamp: expiryTimestamp)

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
            panic("NFT is not locked")
        }

        let lockExpiryTimestamp: UFix64 = self.lockedNFTs[id]!
        let isPastExpiry: Bool = getCurrentBlock().timestamp >= lockExpiryTimestamp

        let isUnlockableOverridden: Bool = self.unlockableNFTs.containsKey(id)

        if !(isPastExpiry || isUnlockableOverridden) {
            panic("NFT is not eligible to be unlocked, expires at ".concat(lockExpiryTimestamp.toString()))
        }

        self.unlockableNFTs.remove(key: id)
        self.lockedNFTs.remove(key: id)

        emit Unlocked(id: id)

        return <- nft
    }

    // Admin is a special authorization resource that 
    // allows the owner to perform important functions to modify the 
    // various aspects of the Plays, Sets, and Moments
    //
    pub resource Admin {
        // createNewAdmin creates a new Admin resource
        //
        pub fun createNewAdmin(): @Admin {
            return <-create Admin()
        }

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

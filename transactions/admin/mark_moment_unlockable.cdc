import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

transaction(ownerAddress: Address, id: UInt64) {
    let adminRef: &TopShotLocking.Admin

    prepare(acct: auth(BorrowValue) &Account) {
        // Set TopShotLocking admin ref
        self.adminRef = acct.storage.borrow<&TopShotLocking.Admin>(from: /storage/TopShotLockingAdmin)
            ?? panic("Could not find reference to TopShotLocking Admin resource")
    }

    execute {
        // Set Top Shot NFT Owner collection ref
        let owner = getAccount(ownerAddress)

        let collectionRef = owner.capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)
            ?? panic("Could not reference owner's moment collection")

        // borrow the nft reference
        let nftRef = collectionRef.borrowNFT(id)!

        // mark the nft as unlockable
        self.adminRef.markNFTUnlockable(nftRef: nftRef)
    }
}

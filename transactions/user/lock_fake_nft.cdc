import TopShot from 0xFAKETOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

// This transaction attempts to send an NFT that is impersonating a TopShot NFT
// to the locking contract, it must fail

// Parameters
//
// id: the Flow ID of the TopShot moment
// duration: number of seconds that the moment will be locked for

transaction(id: UInt64, duration: UFix64) {
    prepare(acct: AuthAccount) {
        let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        let nft <- collectionRef.withdraw(withdrawID: id)

        let lockedNFT <- TopShotLocking.lockNFT(nft: <-nft, duration: duration)

        // destroy here to get rid of loss of resource error - should not actually get here
        destroy <- lockedNFT
    }
}

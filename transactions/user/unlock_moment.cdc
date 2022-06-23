import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS
import TopShotLocking from 0xTOPSHOTLOCKINGADDRESS

transaction(id: UInt64) {
    prepare(acct: AuthAccount) {
        let collectionRef = acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow from MomentCollection in storage")

        collectionRef.unlock(id: id)
    }
}

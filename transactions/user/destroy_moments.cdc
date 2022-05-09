import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS

// This transaction destroys a number of moments owned by a user

// Parameters
//
// momentIDs: an array of moment IDs of NFTs to be destroyed

transaction(momentIDs: [UInt64]) {

    let tokens: @NonFungibleToken.Collection
    
    prepare(acct: AuthAccount) {

        self.tokens <- acct.borrow<&TopShot.Collection>(from: /storage/MomentCollection)!.batchWithdraw(ids: momentIDs)
    }

    execute {
        // destroy the NFTs
        destroy self.tokens
    }
}
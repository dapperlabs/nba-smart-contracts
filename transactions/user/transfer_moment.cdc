import NonFungibleToken from 0xNFTADDRESS
import TopShot from 0xTOPSHOTADDRESS

// This transaction transfers a moment to a recipient

// This transaction is how a topshot user would transfer a moment
// from their account to another account
// The recipient must have a TopShot Collection object stored
// and a public MomentCollectionPublic capability stored at
// `/public/MomentCollection`

// Parameters:
//
// recipient: The Flow address of the account to receive the moment.
// withdrawID: The id of the moment to be transferred

transaction(recipient: Address, withdrawID: UInt64) {

    // local variable for storing the transferred token
    let transferToken: @{NonFungibleToken.NFT}
    
    prepare(acct: auth(BorrowValue) &Account) {

        // borrow a reference to the owner's collection
        let collectionRef = acct.storage.borrow<auth(NonFungibleToken.Withdraw) &TopShot.Collection>(from: /storage/MomentCollection)
            ?? panic("Could not borrow a reference to the stored Moment collection")
        
        // withdraw the NFT
        self.transferToken <- collectionRef.withdraw(withdrawID: withdrawID)
    }

    execute {
        
        // get the recipient's public account object
        let recipient = getAccount(recipient)

        // get the Collection reference for the receiver
        let receiverRef = recipient.capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)!

        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-self.transferToken)
    }
}
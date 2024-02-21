import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS

// This transaction mints DapperUtilityCoin (a Fungible Token) to self,
// then purchases a moment for sale from a seller
// then deposits bought moment to a recipient

// Parameters:
//
// sellerAddress: the Flow address of the account issuing the sale of a moment
// recipient: the Flow address who will receive the moment
// tokenID: the ID of the moment being purchased
// purchaseAmount: the amount for which the user is paying for the moment; must not be less than the moment's price

transaction(sellerAddress: Address, recipient: Address, tokenID: UInt64, purchaseAmount: UFix64) {

    // Local variable for the coin admin
    let ducRef: &DapperUtilityCoin.Minter

    prepare(signer: auth(Storage, Capabilities) &Account) {

        self.ducRef = signer.storage.borrow<&DapperUtilityCoin.Minter>(from: /storage/dapperUtilityCoinAdmin)
            ?? panic("Signer is not the token admin")
    }

    execute {

        let mintedVault <- self.ducRef.mintTokens(amount: purchaseAmount) as! @DapperUtilityCoin.Vault


        let seller = getAccount(sellerAddress)
        
        let topshotSaleCollection = seller.capabilities.borrow<&Market.SaleCollection>(/public/topshotSaleCollection)
            ?? panic("Could not borrow public sale reference")

        let boughtToken <- topshotSaleCollection.purchase(tokenID: tokenID, buyTokens: <-mintedVault)

        // get the recipient's public account object and borrow a reference to their moment receiver
        let recipient = getAccount(recipient).capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)
            ?? panic("Could not borrow a reference to the moment collection")

        // deposit the NFT in the receivers collection
        recipient.deposit(token: <-boughtToken)
    }
}
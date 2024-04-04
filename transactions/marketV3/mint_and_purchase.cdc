import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

transaction(sellerAddress: Address, recipient: Address, tokenID: UInt64, purchaseAmount: UFix64) {

    prepare(signer: auth(BorrowValue) &Account) {

        let tokenAdmin = signer
            .storage.borrow<&DapperUtilityCoin.Minter>(from: /storage/dapperUtilityCoinAdmin)
            ?? panic("Signer is not the token admin")


        let mintedVault <- tokenAdmin.mintTokens(amount: purchaseAmount) as! @DapperUtilityCoin.Vault


        let seller = getAccount(sellerAddress)
        let topshotSaleCollection = seller.capabilities.borrow<&TopShotMarketV3.SaleCollection>(TopShotMarketV3.marketPublicPath)
            ?? panic("Could not borrow public sale reference")

        let boughtToken <- topshotSaleCollection.purchase(tokenID: tokenID, buyTokens: <-mintedVault)

        // get the recipient's public account object and borrow a reference to their moment receiver
        let recipient = getAccount(recipient).capabilities.borrow<&TopShot.Collection>(/public/MomentCollection)
            ?? panic("Could not borrow a reference to the moment collection")

        // deposit the NFT in the receivers collection
        recipient.deposit(token: <-boughtToken)
    }
}

import FungibleToken from 0xFUNGIBLETOKENADDRESS
import DapperUtilityCoin from 0xDUCADDRESS
import TopShot from 0xTOPSHOTADDRESS
import Market from 0xMARKETADDRESS
import TopShotMarketV3 from 0xMARKETV3ADDRESS

transaction(sellerAddress: Address, recipient: Address, tokenID: UInt64, purchaseAmount: UFix64) {

    prepare(signer: AuthAccount) {

        let tokenAdmin = signer
            .borrow<&DapperUtilityCoin.Administrator>(from: /storage/dapperUtilityCoinAdmin) 
            ?? panic("Signer is not the token admin")

        let minter <- tokenAdmin.createNewMinter(allowedAmount: purchaseAmount)
        let mintedVault <- minter.mintTokens(amount: purchaseAmount) as! @DapperUtilityCoin.Vault

        destroy minter

        let seller = getAccount(sellerAddress)
        let topshotSaleCollection = seller.getCapability(TopShotMarketV3.marketPublicPath)
            .borrow<&{Market.SalePublic}>()
            ?? panic("Could not borrow public sale reference")

        let boughtToken <- topshotSaleCollection.purchase(tokenID: tokenID, buyTokens: <-mintedVault)

        // get the recipient's public account object and borrow a reference to their moment receiver
        let recipient = getAccount(recipient)
            .getCapability(/public/MomentCollection).borrow<&{TopShot.MomentCollectionPublic}>()
            ?? panic("Could not borrow a reference to the moment collection")

        // deposit the NFT in the receivers collection
        recipient.deposit(token: <-boughtToken)
    }
}

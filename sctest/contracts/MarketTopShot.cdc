/*

    MarketTopShot.cdc

    Description: Contract definitions for users to sell their moments

    Authors: Joshua Hannan joshua.hannan@dapperlabs.com
             Dieter Shirley dete@axiomzen.com

    Marketplace is where users can create a sale collectio that they
    store in their account storage. In the sale collection, 
    they can put their NFTs up for sale with a price and publish a 
    reference so that others can see the sale.

    If another user sees an NFT that they want to buy,
    they can send fungible tokens that equal or exceed the buy price
    to buy the NFT.  The NFT is transferred to them when
    they make the purchase.

    each user who wants to sell tokens will have a sale collection 
    instance in their account that holds the tokens that they are putting up for sale

    They can give a reference to this collection to the central contract
    that it can list the sales in a central place

*/

import FungibleToken, FlowToken from 0x01
import TopShot from 0x03

pub contract Market {

    pub event MomentListed(id: UInt64, price: UFix64, seller: Address?)
    pub event PriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)
    pub event TokenPurchased(id: UInt64, price: UFix64, seller: Address?)
    pub event SaleWithdrawn(id: UInt64, owner: Address?)
    pub event CutPercentageChanged(newPercent: UFix64, seller: Address?)

    // the reference that is used for depositing TopShot's cut of every sale
    access(contract) var TopShotVault: &{FungibleToken.Receiver}

    // The interface that user can publish to allow others too access their sale
    pub resource interface SalePublic {
        pub var prices: {UInt64: UFix64}
        pub var cutPercentage: UFix64
        pub fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault)
        pub fun getPrice(tokenID: UInt64): UFix64?
        pub fun getIDs(): [UInt64]
    }

    pub resource SaleCollection: SalePublic {

        // a dictionary of the NFTs that the user is putting up for sale
        pub var forSale: @TopShot.Collection

        // dictionary of the prices for each NFT by ID
        pub var prices: {UInt64: UFix64}

        // the fungible token vault of the owner of this sale
        // so that when someone buys a token, this resource can deposit
        // tokens in their account
        access(self) let ownerVault: &{FungibleToken.Receiver}

        // the reference that is used for depositing TopShot's cut of every sale
        access(self) let TopShotVault: &{FungibleToken.Receiver}

        // the percentage that is taken from every purchase for TopShot
        pub var cutPercentage: UFix64

        init (vault: &{FungibleToken.Receiver}, cutPercentage: UFix64) {
            self.forSale <- TopShot.createEmptyCollection()
            self.ownerVault = vault
            self.prices = {}
            self.TopShotVault = Market.TopShotVault
            self.cutPercentage = cutPercentage
        }

        // withdraw gives the owner the opportunity to remove a sale from the collection
        pub fun withdraw(tokenID: UInt64): @TopShot.NFT {
            // remove the price
            self.prices.remove(key: tokenID)
            // remove and return the token
            let token <- self.forSale.withdraw(withdrawID: tokenID)

            emit SaleWithdrawn(id: token.id, owner: self.owner?.address)

            return <-token
        }

        // listForSale lists an NFT for sale in this collection
        pub fun listForSale(token: @TopShot.NFT, price: UFix64) {
            let id: UInt64 = token.id

            self.prices[id] = price

            self.forSale.deposit(token: <-token)

            emit MomentListed(id: id, price: price, seller: self.owner?.address)
        }

        // changePrice changes the price of a token that is currently for sale
        pub fun changePrice(tokenID: UInt64, newPrice: UFix64) {
            pre {
                self.prices[tokenID] != nil: "Cannot change price for a token that doesnt exist."
            }
            self.prices[tokenID] = newPrice

            emit PriceChanged(id: tokenID, newPrice: newPrice, seller: self.owner?.address)
        }

        // changePercentage changes the cut percentage of a token that is currently for sale
        pub fun changePercentage(newPercent: UFix64) {
            self.cutPercentage = newPercent

            emit CutPercentageChanged(newPercent: newPercent, seller: self.owner?.address)
        }

        // purchase lets a user send tokens to purchase an NFT that is for sale
        pub fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault) {
            pre {
                self.forSale.ownedNFTs[tokenID] != nil && self.prices[tokenID] != nil:
                    "No token matching this ID for sale!"
                buyTokens.balance >= (self.prices[tokenID] ?? UFix64(0)):
                    "Not enough tokens to by the NFT!"
            }

            if let price = self.prices[tokenID] {
                self.prices[tokenID] = nil

                // take the cut of the tokens Top shot gets from the sent tokens
                let TopShotCut <- buyTokens.withdraw(amount: price - (price*self.cutPercentage)/(UFix64(100)))

                // deposit it into topshot's Vault
                self.TopShotVault.deposit(from: <-TopShotCut)
                
                // deposit the remaining tokens into the owners vault
                self.ownerVault.deposit(from: <-buyTokens)

                // deposit the NFT into the buyers collection
                recipient.deposit(token: <-self.withdraw(tokenID: tokenID))

                emit TokenPurchased(id: tokenID, price: price, seller: self.owner?.address)
            }
        }

        // idPrice returns the price of a specific token in the sale
        pub fun getPrice(tokenID: UInt64): UFix64? {
            let price = self.prices[tokenID]
            return price
        }

        // getIDs returns an array of token IDs that are for sale
        pub fun getIDs(): [UInt64] {
            return self.forSale.getIDs()
        }

        destroy() {
            destroy self.forSale
        }
    }

    // createCollection returns a new collection resource to the caller
    pub fun createSaleCollection(ownerVault: &{FungibleToken.Receiver}, cutPercentage: UFix64): @SaleCollection {
        return <- create SaleCollection(vault: ownerVault, cutPercentage: cutPercentage)
    }

    init() {
        let acct = getAccount(0x02)
        self.TopShotVault = acct.getCapability(/public/flowTokenReceiver)!
                                .borrow<&FlowToken.Vault{FungibleToken.Receiver}>()!
    }
}
 
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

import FungibleToken from 0x04
import ExampleToken from 0x05
import NonFungibleToken from 0x02
import TopShot from 0x03

pub contract Market {

    pub event MomentListed(id: UInt64, price: UFix64, seller: Address?)
    pub event PriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)
    pub event TokenPurchased(id: UInt64, price: UFix64, seller: Address?)
    pub event SaleWithdrawn(id: UInt64, owner: Address?)
    pub event CutPercentageChanged(newPercent: UFix64, seller: Address?)

    // The interface that user can publish to allow others too access their sale
    pub resource interface SalePublic {
        pub var prices: {UInt64: UFix64}
        pub var cutPercentage: UFix64
        pub fun purchase(tokenID: UInt64, recipient: &AnyResource{TopShot.MomentCollectionPublic}, buyTokens: @ExampleToken.Vault)
        pub fun getPrice(tokenID: UInt64): UFix64?
        pub fun getIDs(): [UInt64]
    }

    pub resource SaleCollection: SalePublic {

        // a dictionary of the NFTs that the user is putting up for sale
        access(self) var forSale: @TopShot.Collection

        // dictionary of the prices for each NFT by ID
        pub var prices: {UInt64: UFix64}

        // the fungible token vault of the owner of this sale
        // so that when someone buys a token, this resource can deposit
        // tokens in their account
        access(self) let ownerCapability: Capability

        // the reference that is used for depositing the beneficiary's cut of every sale
        access(self) let beneficiaryCapability: Capability

        // the percentage that is taken from every purchase for TopShot
        // This is a literal percentage
        // For example, if the percentage is 15%, cutPercentage = 0.15
        pub var cutPercentage: UFix64

        init (ownerCapability: Capability, beneficiaryCapability: Capability, cutPercentage: UFix64) {
            pre {
                ownerCapability.borrow<&{FungibleToken.Receiver}>() != nil: 
                    "Owner's Receiver Capability is invalid!"
                beneficiaryCapability.borrow<&{FungibleToken.Receiver}>() != nil: 
                    "Beneficiary's Receiver Capability is invalid!" 
            }
            
            self.forSale <- TopShot.createEmptyCollection() as! @TopShot.Collection
            self.ownerCapability = ownerCapability
            self.beneficiaryCapability = beneficiaryCapability
            self.prices = {}
            self.cutPercentage = cutPercentage
        }

        // withdraw gives the owner the opportunity to remove a sale from the collection
        pub fun withdraw(tokenID: UInt64): @TopShot.NFT {
            // remove the price
            self.prices.remove(key: tokenID)
            // remove and return the token
            let token <- self.forSale.withdraw(withdrawID: tokenID) as! @TopShot.NFT

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
        pub fun purchase(tokenID: UInt64, recipient: &AnyResource{TopShot.MomentCollectionPublic}, buyTokens: @ExampleToken.Vault) {
            pre {
                self.forSale.ownedNFTs[tokenID] != nil && self.prices[tokenID] != nil:
                    "No token matching this ID for sale!"
                buyTokens.balance == (self.prices[tokenID] ?? UFix64(0)):
                    "Not enough tokens to by the NFT!"
            }

            let price = self.prices[tokenID]!

            self.prices[tokenID] = nil

            // take the cut of the tokens Top shot gets from the sent tokens
            let TopShotCut <- buyTokens.withdraw(amount: price*self.cutPercentage)

            // deposit it into topshot's Vault
            self.beneficiaryCapability.borrow<&{FungibleToken.Receiver}>()!
                .deposit(from: <-TopShotCut)
            
            // deposit the remaining tokens into the owners vault
            self.ownerCapability.borrow<&{FungibleToken.Receiver}>()!
                .deposit(from: <-buyTokens)

            // deposit the NFT into the buyers collection
            recipient.deposit(token: <-self.withdraw(tokenID: tokenID))

            emit TokenPurchased(id: tokenID, price: price, seller: self.owner?.address)
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

        // borrowMoment Returns a borrowed reference to a Moment in the collection
        // so that the caller can read data from it
        pub fun borrowMoment(id: UInt64): &TopShot.NFT {
            post {
                result.id == id: "The ID of the reference is incorrect"
            }
            let ref = self.forSale.borrowMoment(id: id)
            return ref
        }

        destroy() {
            destroy self.forSale
        }
    }

    // createCollection returns a new collection resource to the caller
    pub fun createSaleCollection(ownerCapability: Capability, beneficiaryCapability: Capability, cutPercentage: UFix64): @SaleCollection {
        return <- create SaleCollection(ownerCapability: ownerCapability, beneficiaryCapability: beneficiaryCapability, cutPercentage: cutPercentage)
    }
}
 
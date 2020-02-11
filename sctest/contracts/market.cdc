import FlowToken from 0x0000000000000000000000000000000000000001
import TopShot from 0x0000000000000000000000000000000000000002

// Marketplace is where users can put their NFTs up for sale with a price
// if another user sees an NFT that they want to buy,
// they can send fungible tokens that equal or exceed the buy price
// to buy the NFT.  The NFT is transferred to them when
// they make the purchase


// each user who wants to sell tokens will have a sale collection 
// instance in their account that holds the tokens that they are putting up for sale

// They will give a reference to this collection to the central contract
// that it can use to list tokens


access(all) contract Market {

    access(all) event ForSale(id: UInt64, price: UInt256)
    access(all) event PriceChanged(id: UInt64, newPrice: UInt256)
    access(all) event TokenPurchased(id: UInt64, price: UInt256)
    access(all) event SaleWithdrawn(id: UInt64)

    // the reference that is used for depositing TopShot's cut of every sale
    access(account) var TopShotVault: &FlowToken.Vault
    
    // the percentage that is taken from every purchase for TopShot
    access(account) var cutPercentage: UInt8

    // The collection of sale references that are included in the marketplace
    access(all) var saleReferences: {Int: &SalePublic}

    access(all) var numSales: Int

    // The interface that user can publish to allow others too access their sale
    access(all) resource interface SalePublic {
        access(all) var prices: {UInt64: UInt256}
        access(self) let cutPercentage: UInt8
        access(all) fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault)
        access(all) fun idPrice(tokenID: UInt64): UInt256?
        access(all) fun getIDs(): [UInt64]
    }

    access(all) resource SaleCollection {

        // a dictionary of the NFTs that the user is putting up for sale
        access(all) var forSale: @{UInt64: TopShot.NFT}

        // dictionary of the prices for each NFT by ID
        access(all) var prices: {UInt64: UInt256}

        // the fungible token vault of the owner of this sale
        // so that when someone buys a token, this resource can deposit
        // tokens in their account
        access(account) let ownerVault: &FlowToken.Vault

        // the reference that is used for depositing TopShot's cut of every sale
        access(self) let TopShotVault: &FlowToken.Vault

        // the percentage that is taken from every purchase for TopShot
        access(self) let cutPercentage: UInt8

        init (vault: &FlowToken.Vault) {
            self.forSale <- {}
            self.ownerVault = vault
            self.prices = {}
            self.TopShotVault = Market.TopShotVault
            self.cutPercentage = Market.cutPercentage
        }

        // withdraw gives the owner the opportunity to remove a sale from the collection
        access(all) fun withdraw(tokenID: UInt64): @TopShot.NFT {
            // remove the price
            self.prices.remove(key: tokenID)
            // remove and return the token
            let token <- self.forSale.remove(key: tokenID) ?? panic("missing NFT")

            emit SaleWithdrawn(id: tokenID)

            return <-token
        }

        // listForSale lists an NFT for sale in this collection
        access(all) fun listForSale(token: @TopShot.NFT, price: UInt256) {
            let id: UInt64 = token.id

            self.prices[id] = price

            let oldToken <- self.forSale[id] <- token

            emit ForSale(id: id, price: price)

            destroy oldToken
        }

        // changePrice changes the price of a token that is currently for sale
        access(all) fun changePrice(tokenID: UInt64, newPrice: UInt256) {
            pre {
                self.prices[tokenID] != nil: "Cannot change price for a token that doesnt exist."
            }
            self.prices[tokenID] = newPrice

            emit PriceChanged(id: tokenID, newPrice: newPrice)
        }

        // purchase lets a user send tokens to purchase an NFT that is for sale
        access(all) fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault) {
            pre {
                self.forSale[tokenID] != nil && self.prices[tokenID] != nil:
                    "No token matching this ID for sale!"
                buyTokens.balance >= (self.prices[tokenID] ?? UInt256(0)):
                    "Not enough tokens to by the NFT!"
            }

            if let price = self.prices[tokenID] {
                self.prices[tokenID] = nil

                // take the cut of the tokens Top shot gets from the sent tokens
                let TopShotCut = buyTokens.withdraw(amount: price - (price*self.cutPercentage)/100)

                // deposit it into topshot's Vault
                self.TopShotVault.deposit(from: <-TopShotCut)
                
                // deposit the remaining tokens into the owners vault
                self.ownerVault.deposit(from: <-buyTokens)

                // deposit the NFT into the buyers collection
                recipient.deposit(token: <-self.withdraw(tokenID: tokenID))

                emit TokenPurchased(id: tokenID, price: price)
            }
        }

        // idPrice returns the price of a specific token in the sale
        access(all) fun idPrice(tokenID: UInt64): UInt256? {
            let price = self.prices[tokenID]
            return price
        }

        // getIDs returns an array of token IDs that are for sale
        access(all) fun getIDs(): [UInt64] {
            return self.forSale.keys
        }

        destroy() {
            destroy self.forSale
        }
    }

    // createCollection returns a new collection resource to the caller
    access(all) fun createSaleCollection(ownerVault: &FlowToken.Vault): @SaleCollection {
        return <- create SaleCollection(vault: ownerVault)
    }

    // These next three functions may or may not be needed but serve as a
    // preliminary way for the contract to keep track of sales that are
    // listed in the marketplace
    access(account) fun addSale(reference: &SalePublic): Int {
        pre {
            reference.getIDs().length != 0: "Cannot add an empty sale!"
        }

        self.saleReferences[self.numSales] = reference
        self.numSales = self.numSales + 1

        return self.numSales
    }

    access(account) fun removeSale(id: Int) {
        self.saleReferences.remove(key: id)
        self.numSales = self.numSales + 1
    }

    access(all) fun getSaleReference(id: Int): &SalePublic? {
        return self.saleReferences[id]
    }

    init() {
        self.TopShotVault = self.account.storage[&FlowToken.Vault]
        self.cutPercentage = 5
        self.saleReferences = {}
        self.numSales = 1
    }
}

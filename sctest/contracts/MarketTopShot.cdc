import FungibleToken, FlowToken from 0x01
import TopShot from 0x03

// Marketplace is where users can put their NFTs up for sale with a price
// if another user sees an NFT that they want to buy,
// they can send fungible tokens that equal or exceed the buy price
// to buy the NFT.  The NFT is transferred to them when
// they make the purchase

// each user who wants to sell tokens will have a sale collection 
// instance in their account that holds the tokens that they are putting up for sale

// They will give a reference to this collection to the central contract
// that it can use to list tokens

pub contract Market {

    pub event ForSale(id: UInt64, price: UInt64)
    pub event PriceChanged(id: UInt64, newPrice: UInt64)
    pub event TokenPurchased(id: UInt64, price: UInt64)
    pub event SaleWithdrawn(id: UInt64)
    pub event CutPercentageChanged(newPercent: UInt64)

    // the reference that is used for depositing TopShot's cut of every sale
    access(contract) var TopShotVault: &FungibleToken.Receiver

    // The collection of sale references that are included in the marketplace
    pub var saleReferences: {Int: &SalePublic}

    // The number of sales that have been listed in this contract
    // used as an index to the central sale listings
    pub var numSales: Int

    // The interface that user can publish to allow others too access their sale
    pub resource interface SalePublic {
        pub var prices: {UInt64: UInt64}
        pub var cutPercentage: UInt64
        pub fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault)
        pub fun idPrice(tokenID: UInt64): UInt64?
        pub fun getIDs(): [UInt64]
    }

    pub resource SaleCollection: SalePublic {

        // a dictionary of the NFTs that the user is putting up for sale
        pub var forSale: @TopShot.Collection

        // dictionary of the prices for each NFT by ID
        pub var prices: {UInt64: UInt64}

        // the fungible token vault of the owner of this sale
        // so that when someone buys a token, this resource can deposit
        // tokens in their account
        access(self) let ownerVault: &FungibleToken.Receiver

        // the reference that is used for depositing TopShot's cut of every sale
        access(self) let TopShotVault: &FungibleToken.Receiver

        // the percentage that is taken from every purchase for TopShot
        pub var cutPercentage: UInt64

        init (vault: &FungibleToken.Receiver, cutPercentage: UInt64) {
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

            emit SaleWithdrawn(id: token.id)

            return <-token
        }

        // listForSale lists an NFT for sale in this collection
        pub fun listForSale(token: @TopShot.NFT, price: UInt64) {
            let id: UInt64 = token.id

            self.prices[id] = price

            self.forSale.deposit(token: <-token)

            emit ForSale(id: id, price: price)
        }

        // changePrice changes the price of a token that is currently for sale
        pub fun changePrice(tokenID: UInt64, newPrice: UInt64) {
            pre {
                self.prices[tokenID] != nil: "Cannot change price for a token that doesnt exist."
            }
            self.prices[tokenID] = newPrice

            emit PriceChanged(id: tokenID, newPrice: newPrice)
        }

        // changePercentage changes the cut percentage of a token that is currently for sale
        pub fun changePercentage(newPercent: UInt64) {
            self.cutPercentage = newPercent

            emit CutPercentageChanged(newPercent: newPercent)
        }

        // purchase lets a user send tokens to purchase an NFT that is for sale
        pub fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault) {
            pre {
                self.forSale.ownedNFTs[tokenID] != nil && self.prices[tokenID] != nil:
                    "No token matching this ID for sale!"
                buyTokens.balance >= (self.prices[tokenID] ?? UInt64(0)):
                    "Not enough tokens to by the NFT!"
            }

            if let price = self.prices[tokenID] {
                self.prices[tokenID] = nil

                // take the cut of the tokens Top shot gets from the sent tokens
                let TopShotCut <- buyTokens.withdraw(amount: price - (price*self.cutPercentage)/(UInt64(100)))

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
        pub fun idPrice(tokenID: UInt64): UInt64? {
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
    pub fun createSaleCollection(ownerVault: &FungibleToken.Receiver, cutPercentage: UInt64): @SaleCollection {
        return <- create SaleCollection(vault: ownerVault, cutPercentage: cutPercentage)
    }

    // These next three functions may or may not be needed but serve as a
    // preliminary way for the contract to keep track of sales that are
    // listed in the marketplace
    access(account) fun addSale(reference: &SalePublic): Int {
        pre {
            reference.getIDs().length != 0: "Cannot add an empty sale!"
        }

        self.numSales = self.numSales + 1

        self.saleReferences[self.numSales] = reference

        return self.numSales
    }

    access(account) fun removeSale(id: Int) {
        self.saleReferences.remove(key: id)
    }

    pub fun getSaleReference(id: Int): &SalePublic {
        return self.saleReferences[id] ?? panic("No reference!")
    }

    init() {
        let acct = getAccount(0x02)
        self.TopShotVault = acct.published[&FungibleToken.Receiver] ?? panic("No vault!")
        self.saleReferences = {}
        self.numSales = 0
    }
}
 
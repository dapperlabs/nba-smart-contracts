import FungibleToken, FlowToken from 0x01
import TopShot from 0x02

// Marketplace is where users can put their NFTs up for sale with a price
// if another user sees an NFT that they want to buy,
// they can send fungible tokens that equal or exceed the buy price
// to buy the NFT.  The NFT is transferred to them when
// they make the purchase

// This version of the Sale contract defines a resource that only can
// hole one NFT for sale at a time, when the user wants to start a new
// sale, a new resource is created, and when the sale completes, the 
// resource can be safely destroyed

// This way, users can list separate cutPercentages for different sales
// and give out separate references to different sales.

// In the future, once references and downcasting are implemented, we will
// be able to have generic NFTs and cut Receivers in 

// They will give a reference to this collection to the central contract
// that it can use to list tokens


access(all) contract SingleSaleMarket {

    access(all) event ForSale(id: UInt64, price: UInt64, cut: UInt64)
    access(all) event PriceChanged(newPrice: UInt64)
    access(all) event TokenPurchased(id: UInt64, price: UInt64)
    access(all) event Withdraw(id: UInt64)
    access(all) event CutPercentageChanged(newPercent: UInt64)

    // The collection of sale references that are included in the marketplace
    access(all) var saleReferences: {Int: &SalePublic}

    // The number of sales in this marketplace
    access(all) var numSales: Int

    // The interface that user can publish to allow others too access their sale
    access(all) resource interface SalePublic {
        access(all) var price: UInt64
        access(all) var cutPercentage: UInt64
        access(all) fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault)
        access(all) fun getID(): UInt64?
        access(all) fun getMoldID(id: UInt64): UInt32?
        access(all) fun getQuality(id: UInt64): Int?
        access(all) fun getPlaceInQuality(id: UInt64): UInt32?
        access(all) fun getMetaData(id: UInt64): {String: String}?
    }

    access(all) resource Sale: SalePublic {

        // a dictionary of the NFTs that the user is putting up for sale
        access(all) var forSale: @TopShot.NFT?

        // price of the NFT for sale
        access(all) var price: UInt64

        // the fungible token vault of the owner of this sale
        // so that when someone buys a token, this resource can deposit
        // tokens in their account
        access(account) let ownerVault: &FungibleToken.Receiver

        // the reference that is used for depositing TopShot's cut of every sale
        access(self) let TopShotVault: &FungibleToken.Receiver

        // the percentage that is taken from every purchase for TopShot
        access(all) var cutPercentage: UInt64

        init (forSale: @TopShot.NFT, 
              vault: &FungibleToken.Receiver, 
              price: UInt64, 
              cutVault: &FungibleToken.Receiver, 
              cut: UInt64) {
            let id = forSale.id
            self.forSale <- forSale
            self.ownerVault = vault
            self.price = price
            self.TopShotVault = cutVault
            self.cutPercentage = cut

            emit ForSale(id: id, price: price, cut: cut)
        }

        // changePrice changes the price of the token that is currently for sale
        access(all) fun changePrice(newPrice: UInt64) {
            self.price = newPrice

            emit PriceChanged(newPrice: newPrice)
        }

        // changePercentage changes the cut percentage of a token that is currently for sale
        access(all) fun changePercentage(newPercent: UInt64) {
            self.cutPercentage = newPercent

            emit CutPercentageChanged(newPercent: newPercent)
        }

        // purchase lets a user send tokens to purchase an NFT that is for sale
        access(all) fun purchase(tokenID: UInt64, recipient: &TopShot.Collection, buyTokens: @FlowToken.Vault) {
            pre {
                buyTokens.balance >= self.price:
                    "Not enough tokens to by the NFT!"
            }

            // take the cut of the tokens Top shot gets from the sent tokens
            let TopShotCut <- buyTokens.withdraw(amount: self.price - (self.price*self.cutPercentage)/(UInt64(100)))

            // deposit it into topshot's Vault
            self.TopShotVault.deposit(from: <-TopShotCut)
            
            // deposit the remaining tokens into the owners vault
            self.ownerVault.deposit(from: <-buyTokens)

            let nft <- self.forSale ?? panic("No NFT!")

            // deposit the NFT into the buyers collection
            recipient.deposit(token: <-nft)

            emit TokenPurchased(id: tokenID, price: self.price)
        }

        // idPrice returns the price of a specific token in the sale
        access(all) fun getPrice(): UInt64 {
            return self.price
        }

        // getIDs returns an array of token IDs that are for sale
        access(all) fun getID(): UInt64? {
            return self.forSale?.id
        }

        // getMoldID gets the mold ID of the moment that is for sale
        access(all) fun getMoldID(id: UInt64): UInt32? {
            return self.forSale?.moldID
        }

        // getQuality gets the quality of the moment that is for sale
        access(all) fun getQuality(id: UInt64): Int? {
            return self.forSale?.quality
        }

        // getPlaceInQuality gets the place in the quality count
        // of the moment that is for sale
        access(all) fun getPlaceInQuality(id: UInt64): UInt32? {
            return self.forSale?.placeInQuality
        }

        // getMetaData gets the metadata for the moment that is for sale
        access(all) fun getMetaData(id: UInt64): {String: String}? {
            if let moldID = self.getMoldID(id: id) {
                return TopShot.molds[moldID]?.metadata
            } else {
                return nil
            }
        }

        destroy() {
            destroy self.forSale
        }
    }

    // createCollection returns a new collection resource to the caller
    access(all) fun createSaleCollection(forSale: @TopShot.NFT, vault: &FungibleToken.Receiver, price: UInt64, cutVault: &FungibleToken.Receiver, cut: UInt64): @Sale {
        return <- create Sale(forSale: <-forSale, vault: vault, price: price, cutVault: cutVault, cut: cut)
    }

    // These next three functions may or may not be needed but serve as a
    // preliminary way for the contract to keep track of sales that are
    // listed in the marketplace
    access(account) fun addSale(reference: &SalePublic): Int {

        self.numSales = self.numSales + 1

        self.saleReferences[self.numSales] = reference

        return self.numSales
    }

    access(account) fun removeSale(id: Int) {
        self.saleReferences.remove(key: id)
    }

    access(all) fun getSaleReference(id: Int): &SalePublic {
        return self.saleReferences[id] ?? panic("No reference!")
    }

    init() {
        self.saleReferences = {}
        self.numSales = 0
    }
}
 
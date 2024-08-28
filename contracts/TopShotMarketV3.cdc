/**

    TopShotMarketV3.cdc

    Description: Contract definitions for users to sell their moments

    Marketplace is where users can create a sale collection that they
    store in their account storage. In the sale collection, 
    they can put their NFTs up for sale with a price and publish a 
    reference so that others can see the sale.

    If another user sees an NFT that they want to buy,
    they can send fungible tokens that equal the buy price
    to buy the NFT.  The NFT is transferred to them when
    they make the purchase.

    Each user who wants to sell tokens will have a sale collection 
    instance in their account that contains price information 
    for each node in their collection. The sale holds a capability that 
    links to their main moment collection.

    They can give a reference to this collection to a central contract
    so that it can list the sales in a central place

    When a user creates a sale, they will supply four arguments:
    - A TopShot.Collection capability that allows their sale to withdraw
      a moment when it is purchased.
    - A FungibleToken.Receiver capability as the place where the payment for the token goes.
    - A FungibleToken.Receiver capability specifying a beneficiary, where a cut of the purchase gets sent. 
    - A cut percentage, specifying how much the beneficiary will recieve.
    
    The cut percentage can be set to zero if the user desires and they 
    will receive the entirety of the purchase. TopShot will initialize sales 
    for users with the TopShot admin vault as the vault where cuts get 
    deposited to.
**/

import FungibleToken from "FungibleToken"
import NonFungibleToken from "NonFungibleToken"
import TopShot from "TopShot"
import Market from "Market"
import DapperUtilityCoin from "DapperUtilityCoin"
import TopShotLocking from "TopShotLocking"
import MetadataViews from "MetadataViews"

access(all) contract TopShotMarketV3 {

    access(all) entitlement Create
    access(all) entitlement Cancel
    access(all) entitlement Update

    // -----------------------------------------------------------------------
    // TopShot Market contract Event definitions
    // -----------------------------------------------------------------------

    /// emitted when a TopShot moment is listed for sale
    access(all) event MomentListed(id: UInt64, price: UFix64, seller: Address?)
    /// emitted when the price of a listed moment has changed
    access(all) event MomentPriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)
    /// emitted when a token is purchased from the market
    access(all) event MomentPurchased(id: UInt64, price: UFix64, seller: Address?, momentName: String, momentDescription: String, momentThumbnailURL: String)
    /// emitted when a moment has been withdrawn from the sale
    access(all) event MomentWithdrawn(id: UInt64, owner: Address?)

    /// Path where the `SaleCollection` is stored
    access(all) let marketStoragePath: StoragePath

    /// Path where the public capability for the `SaleCollection` is
    access(all) let marketPublicPath: PublicPath

    /// SaleCollection
    ///
    /// This is the main resource that token sellers will store in their account
    /// to manage the NFTs that they are selling. The SaleCollection
    /// holds a TopShot Collection resource to store the moments that are for sale.
    /// The SaleCollection also keeps track of the price of each token.
    /// 
    /// When a token is purchased, a cut is taken from the tokens
    /// and sent to the beneficiary, then the rest are sent to the seller.
    ///
    /// The seller chooses who the beneficiary is and what percentage
    /// of the tokens gets taken from the purchase
    access(all) resource SaleCollection: Market.SalePublic {

        /// A collection of the moments that the user has for sale
        access(self) var ownerCollection: Capability<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>

        /// Capability to point at the V1 sale collection
        access(self) var marketV1Capability: Capability<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>?

        /// Dictionary of the low low prices for each NFT by ID
        access(self) var prices: {UInt64: UFix64}

        /// The fungible token vault of the seller
        /// so that when someone buys a token, the tokens are deposited
        /// to this Vault
        access(self) var ownerCapability: Capability<&{FungibleToken.Receiver}>

        /// The capability that is used for depositing 
        /// the beneficiary's cut of every sale
        access(self) var beneficiaryCapability: Capability<&{FungibleToken.Receiver}>

        /// The percentage that is taken from every purchase for the beneficiary
        /// For example, if the percentage is 15%, cutPercentage = 0.15
        access(all) var cutPercentage: UFix64

        init (ownerCollection: Capability<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>,
              ownerCapability: Capability<&{FungibleToken.Receiver}>,
              beneficiaryCapability: Capability<&{FungibleToken.Receiver}>,
              cutPercentage: UFix64,
              marketV1Capability: Capability<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>?) {
            pre {
                // Check that the owner's moment collection capability is correct
                ownerCollection.check(): 
                    "Owner's Moment Collection Capability is invalid!"

                // Check that both capabilities are for fungible token Vault receivers
                ownerCapability.check(): 
                    "Owner's Receiver Capability is invalid!"
                beneficiaryCapability.check(): 
                    "Beneficiary's Receiver Capability is invalid!" 

                // Make sure the V1 sale collection capability is valid
                marketV1Capability == nil || marketV1Capability!.check():
                    "V1 Market Capability is invalid"
            }
            
            // create an empty collection to store the moments that are for sale
            self.ownerCollection = ownerCollection
            self.ownerCapability = ownerCapability
            self.beneficiaryCapability = beneficiaryCapability
            // prices are initially empty because there are no moments for sale
            self.prices = {}
            self.cutPercentage = cutPercentage
            self.marketV1Capability = marketV1Capability
        }

        /// listForSale lists an NFT for sale in this sale collection
        /// at the specified price
        ///
        /// Parameters: tokenID: The id of the NFT to be put up for sale
        ///             price: The price of the NFT
        access(Create) fun listForSale(tokenID: UInt64, price: UFix64) {
            pre {
                self.ownerCollection.borrow()!.borrowMoment(id: tokenID) != nil:
                    "Moment does not exist in the owner's collection"

                !TopShotLocking.isLocked(nftRef: self.ownerCollection.borrow()!.borrowNFT(tokenID)!):
                    "Moment is locked"
            }

            // Set the token's price
            self.prices[tokenID] = price

            emit MomentListed(id: tokenID, price: price, seller: self.owner?.address)
        }

        /// cancelSale cancels a moment sale and clears its price
        ///
        /// Parameters: tokenID: the ID of the token to withdraw from the sale
        ///
        access(Cancel) fun cancelSale(tokenID: UInt64) {
            
            // First check this version of the sale
            if self.prices[tokenID] != nil {
                // Remove the price from the prices dictionary
                self.prices.remove(key: tokenID)

                // Set prices to nil for the withdrawn ID
                self.prices[tokenID] = nil
                
                // Emit the event for withdrawing a moment from the Sale
                emit MomentWithdrawn(id: tokenID, owner: self.owner?.address)

            // If not found in this SaleCollection, check V1
            } else if let v1Market = self.marketV1Capability {
                let v1MarketRef = v1Market.borrow()!

                if v1MarketRef.getPrice(tokenID: tokenID) != nil {
                    // withdraw the moment from the v1 collection
                    let token <- v1MarketRef.withdraw(tokenID: tokenID)

                    // borrow a reference to the main top shot collection
                    let ownerCollectionRef = self.ownerCollection.borrow()
                        ?? panic("Could not borrow owner collection reference")

                    // deposit the withdrawn moment into the main collection
                    ownerCollectionRef.deposit(token: <-token)
                }
            }
        }

        /// purchase lets a user send tokens to purchase an NFT that is for sale
        /// the purchased NFT is returned to the transaction context that called it
        ///
        /// Parameters: tokenID: the ID of the NFT to purchase
        ///             buyTokens: the fungible tokens that are used to buy the NFT
        ///
        /// Returns: @TopShot.NFT: the purchased NFT
        access(all) fun purchase(tokenID: UInt64, buyTokens: @DapperUtilityCoin.Vault): @TopShot.NFT {

            // First check this sale collection for the NFT
            if self.prices[tokenID] != nil {
                assert(
                    buyTokens.balance == self.prices[tokenID]!,
                    message: "Not enough tokens to buy the NFT!"
                )

                // Read the price for the token
                let price = self.prices[tokenID]!

                // Set the price for the token to nil
                self.prices[tokenID] = nil

                // Take the cut of the tokens that the beneficiary gets from the sent tokens
                let beneficiaryCut <- buyTokens.withdraw(amount: price*self.cutPercentage)

                // Deposit it into the beneficiary's Vault
                self.beneficiaryCapability.borrow()!
                    .deposit(from: <-beneficiaryCut)
                
                // Deposit the remaining tokens into the owners vault
                self.ownerCapability.borrow()!
                    .deposit(from: <-buyTokens)

                // Return the purchased token
                let boughtMoment <- self.ownerCollection.borrow()!.withdraw(withdrawID: tokenID) as! @TopShot.NFT

                let momentDisplay = boughtMoment.resolveView(Type<MetadataViews.Display>())! as! MetadataViews.Display
                emit MomentPurchased(id: tokenID, price: price, seller: self.owner?.address, momentName: momentDisplay.name, momentDescription: momentDisplay.description, momentThumbnailURL: momentDisplay.thumbnail.uri())

                return <-boughtMoment

            // If not found in this SaleCollection, check V1
            } else if let v1Market = self.marketV1Capability {
                let v1MarketRef = v1Market.borrow()!

                return <-v1MarketRef.purchase(tokenID: tokenID, buyTokens: <-buyTokens)
            } 
            
            // Refactored to avoid dead code to resolve
            // https://github.com/dapperlabs/nba-smart-contracts/issues/165
            panic("No token matching this ID for sale!")

        }

        /// changeOwnerReceiver updates the capability for the sellers fungible token Vault
        ///
        /// Parameters: newOwnerCapability: The new fungible token capability for the account 
        ///                                 who received tokens for purchases
        access(Update) fun changeOwnerReceiver(_ newOwnerCapability: Capability<&{FungibleToken.Receiver}>) {
            pre {
                newOwnerCapability.borrow() != nil: 
                    "Owner's Receiver Capability is invalid!"
            }
            self.ownerCapability = newOwnerCapability
        }

        /// changeBeneficiaryReceiver updates the capability for the beneficiary of the cut of the sale
        ///
        /// Parameters: newBeneficiaryCapability the new capability for the beneficiary of the cut of the sale
        ///
        access(Update) fun changeBeneficiaryReceiver(_ newBeneficiaryCapability: Capability<&{FungibleToken.Receiver}>) {
            pre {
                newBeneficiaryCapability.borrow() != nil: 
                    "Beneficiary's Receiver Capability is invalid!" 
            }
            self.beneficiaryCapability = newBeneficiaryCapability
        }

        /// getPrice returns the price of a specific token in the sale
        /// 
        /// Parameters: tokenID: The ID of the NFT whose price to get
        ///
        /// Returns: UFix64: The price of the token
        access(all) view fun getPrice(tokenID: UInt64): UFix64? {
            if let price = self.prices[tokenID] {
                return price
            } else if let marketV1 = self.marketV1Capability {
                return marketV1.borrow()!.getPrice(tokenID: tokenID)
            }
            return nil
        }

        /// getIDs returns an array of token IDs that are for sale
        access(all) view fun getIDs(): [UInt64] {
            let v3Keys = self.prices.keys

            // Add any V1 SaleCollection IDs if they exist
            if let marketV1 = self.marketV1Capability {
                let v1Keys = marketV1.borrow()!.getIDs()
                return v1Keys.concat(v3Keys)
            } else {
                return v3Keys
            }
        }

        /// borrowMoment Returns a borrowed reference to a Moment for sale
        /// so that the caller can read data from it
        ///
        /// Parameters: id: The ID of the moment to borrow a reference to
        ///
        /// Returns: &TopShot.NFT? Optional reference to a moment for sale 
        ///                        so that the caller can read its data
        ///
        access(all) view fun borrowMoment(id: UInt64): &TopShot.NFT? {
            // first check this collection
            if self.prices[id] != nil {
                let ref = self.ownerCollection.borrow()!.borrowMoment(id: id)
                return ref
            } else {
                // If it wasn't found here, check the V1 SaleCollection
                if let marketV1 = self.marketV1Capability {
                    return marketV1.borrow()!.borrowMoment(id: id)
                }
                return nil
            }
        }
    }

    /// createCollection returns a new collection resource to the caller
    access(all) fun createSaleCollection(ownerCollection: Capability<auth(NonFungibleToken.Withdraw, NonFungibleToken.Update) &TopShot.Collection>,
                                 ownerCapability: Capability<&{FungibleToken.Receiver}>,
                                 beneficiaryCapability: Capability<&{FungibleToken.Receiver}>,
                                 cutPercentage: UFix64,
                                 marketV1Capability: Capability<auth(Market.Create, NonFungibleToken.Withdraw, Market.Update) &Market.SaleCollection>?): @SaleCollection {

        return <- create SaleCollection(ownerCollection: ownerCollection,
                                        ownerCapability: ownerCapability,
                                        beneficiaryCapability: beneficiaryCapability,
                                        cutPercentage: cutPercentage,
                                        marketV1Capability: marketV1Capability)
    }

    init() {
        self.marketStoragePath = /storage/topshotSale3Collection
        self.marketPublicPath = /public/topshotSalev3Collection
    }
}
 
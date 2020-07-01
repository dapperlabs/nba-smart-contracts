/*
    MarketTopShot.cdc

    Description: Contract definitions for users to sell their moments

    Authors: Joshua Hannan joshua.hannan@dapperlabs.com
             Dieter Shirley dete@axiomzen.com
*/

import TopShot from 0x179b6b1cb6755e31

pub contract TopShotTrading {

    // -----------------------------------------------------------------------
    // TopShot Market contract Event definitions
    // -----------------------------------------------------------------------

    // emitted when a TopShot moment is listed for sale
    pub event MomentListed(id: UInt64, seller: Address?)
    // emitted when the price of a listed moment has changed
    pub event MomentRequirementChanged(id: UInt64, seller: Address?)
    // emitted when a token is purchased from the market
    pub event MomentTraded(id: UInt64, seller: Address?)
    // emitted when a moment has been removed from this trade
    pub event MomentWithdrawn(id: UInt64, owner: Address?)

    // TradePublic 
    //
    // The interface that a user can publish their trades as 
    // to allow others to access their trades
    pub resource interface TradePublic {
        pub fun completeTrade(tokenID: UInt64, tokenToTrade: @TopShot.NFT): @TopShot.NFT {
            post {
                result.id == tokenID: "The ID of the withdrawn token must be the same as the requested ID"
            }
        }
        pub fun getTradesOffered(): {UInt64: TradeRequirement}
        pub fun getRequirements(tokenID: UInt64): TradeRequirement?
        pub fun getIDs(): [UInt64]?
        pub fun borrowMoment(id: UInt64): &TopShot.NFT? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.id == id): 
                    "Cannot borrow Moment reference: The ID of the returned reference is incorrect"
            }
        }
    }

    pub struct TradeRequirement {
        pub(set) var idMin: UInt64
        pub(set) var idMax: UInt64

        pub(set) var setIDMin: UInt32
        pub(set) var setIDMax: UInt32

        pub(set) var playIDMin: UInt32
        pub(set) var playIDMax: UInt32

        pub(set) var serialNumMin: UInt32
        pub(set) var serialNumMax: UInt32

        init(idMin: UInt64, idMax: UInt64,
             setIDMin: UInt32, setIDMax: UInt32,
             playIDMin: UInt32, playIDMax: UInt32,
             serialNumMin: UInt32, serialNumMax: UInt32) {
            
            pre {
                idMin <= idMax: "IDs range must be valid"
                setIDMin <= setIDMax: "SetID range must be valid"
                playIDMin <= playIDMax: "PlayID range must be valid"
                serialNumMin <= serialNumMax: "Serial number range must be valid"
            }
            
            self.idMin = idMin
            self.idMax = idMax
            self.setIDMin = setIDMin
            self.setIDMax = setIDMax
            self.playIDMin = playIDMin
            self.playIDMax = playIDMax
            self.serialNumMin = serialNumMin
            self.serialNumMax = serialNumMax
        }
    }

    // TradingCenter
    //
    pub resource TradingCenter: TradePublic {

        // The user's main moment collection
        access(self) var ownerCollection: Capability

        // Dictionary of the prices for each NFT by ID
        access(self) var requirements: {UInt64: TradeRequirement}

        init (collectionCapability: Capability) {
            pre {
                // Check that the capabilities are for topshot collections
                collectionCapability.borrow<&TopShot.Collection>() != nil: 
                    "Owner's collection Capability is invalid!"
            }
            
            self.ownerCollection = collectionCapability
            self.requirements = {}
        }

        // listForTrade lists an NFT for sale in this sale collection
        // with the specified trade requirements
        pub fun listForTrade(tokenID: UInt64,
                             idMin: UInt64, idMax: UInt64,
                             setIDMin: UInt32, setIDMax: UInt32,
                             playIDMin: UInt32, playIDMax: UInt32,
                             serialNumMin: UInt32, serialNumMax: UInt32) 
        {
            pre {
                self.ownerCollection.borrow<&TopShot.Collection>()!.borrowMoment(id: tokenID) != nil: "Cannot list a moment that you don't own"
            }

            // Set the token's price
            self.requirements[tokenID] = TradeRequirement(idMin: idMin, idMax: idMax,
                                                          setIDMin: setIDMin, setIDMax: setIDMax,
                                                          playIDMin: playIDMin, playIDMax: playIDMax,
                                                          serialNumMin: serialNumMin, serialNumMax: serialNumMax)

            emit MomentListed(id: tokenID, seller: self.owner?.address)
        }

        // 
        pub fun deListMoment(tokenID: UInt64) {

            // Remove the price from the prices dictionary
            self.requirements.remove(key: tokenID)

            // set prices to nil for the withdrawn ID
            self.requirements[tokenID] = nil
            
            // Emit the event for withdrawing a moment from the Sale
            emit MomentWithdrawn(id: tokenID, owner: self.owner?.address)
        }

        // completeTrade lets a user send a NFT to complete a trade that is listed
        // The moment that is sent must meet the requirements that the lister
        // specified
        pub fun completeTrade(tokenID: UInt64, tokenToTrade: @TopShot.NFT): @TopShot.NFT {
            pre {
                self.requirements[tokenID] != nil:
                    "No token matching this ID for trade!"
                tokenToTrade.id >= self.requirements[tokenID]!.idMin && tokenToTrade.id <= self.requirements[tokenID]!.idMax
                tokenToTrade.data.setID >= self.requirements[tokenID]!.setIDMin && tokenToTrade.data.setID <= self.requirements[tokenID]!.setIDMax
                tokenToTrade.data.playID >= self.requirements[tokenID]!.playIDMin && tokenToTrade.data.playID <= self.requirements[tokenID]!.playIDMax
            }

            let requirements = self.requirements[tokenID]

            self.requirements[tokenID] = nil

            emit MomentTraded(id: tokenID, seller: self.owner?.address)

            let collectionRef = self.ownerCollection.borrow<&TopShot.Collection>()!

            collectionRef.deposit(token: <-tokenToTrade)

            // return the purchased token
            let tradedNFT <-collectionRef.withdraw(withdrawID: tokenID) as! @TopShot.NFT

            return <-tradedNFT
        }

        pub fun getTradesOffered(): {UInt64: TradeRequirement} {
            return self.requirements
        }

        // getIDs returns an array of token IDs that are for trade
        pub fun getIDs(): [UInt64]? {
            return self.requirements.keys
        }

        // getPrice returns the price of a specific token in the sale
        pub fun getRequirements(tokenID: UInt64): TradeRequirement? {
            return self.requirements[tokenID]
        }

        // borrowMoment Returns a borrowed reference to a Moment in the collection
        // so that the caller can read data from it
        pub fun borrowMoment(id: UInt64): &TopShot.NFT? {
            let collectionRef = self.ownerCollection.borrow<&TopShot.Collection>()!
            let ref = collectionRef.borrowMoment(id: id)

            return ref
        }
    }

    // createTradingCenter returns a new trading center resource to the caller
    pub fun createTradingCenter(ownerCollection: Capability): @TradingCenter {
        return <- create TradingCenter(collectionCapability: ownerCollection)
    }
}
 
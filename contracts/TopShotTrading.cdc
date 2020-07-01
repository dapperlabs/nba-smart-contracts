/*
    TopShotTrading.cdc

    Description: Contract definitions for users to trade their moments

    Authors: Joshua Hannan joshua.hannan@dapperlabs.com

    TopShotTrading defines a resource object similar to the market that
    allows a user to put their moments up for trading.

    When a user specifies that they want to trade a moment, they say 
    what their requirements are for a trade by indicating which range of 
    moment ID, set ID, play ID, and serial number they want. 

    They specify a min and a max for each one, in case they are ok with
    multiple moments that meet the same criteria
*/

import TopShot from 0x179b6b1cb6755e31

pub contract TopShotTrading {

    // -----------------------------------------------------------------------
    // TopShot Trading contract Event definitions
    // -----------------------------------------------------------------------

    // emitted when a TopShot moment is listed for trading
    pub event MomentListed(id: UInt64, seller: Address?)
    // emitted when the requirements of a listed moment has changed
    pub event MomentRequirementChanged(id: UInt64, seller: Address?)
    // emitted when a token is traded from a trading center
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

    // Indicates what kind of moment the user wants to trade for
    // The user specifies a min and max for what they are comfortable
    // with each ID. Set them as the same to only choose one id, set, play
    // or serial number
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

        // Dictionary of the requirements for each NFT trade by ID
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

            // Set the token's requirements
            self.requirements[tokenID] = TradeRequirement(idMin: idMin, idMax: idMax,
                                                          setIDMin: setIDMin, setIDMax: setIDMax,
                                                          playIDMin: playIDMin, playIDMax: playIDMax,
                                                          serialNumMin: serialNumMin, serialNumMax: serialNumMax)

            emit MomentListed(id: tokenID, seller: self.owner?.address)
        }

        // remove the moment from trading eligibility
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
                // Make sure the metadata for the moment matches the trading requirements
                tokenToTrade.id >= self.requirements[tokenID]!.idMin && tokenToTrade.id <= self.requirements[tokenID]!.idMax
                tokenToTrade.data.setID >= self.requirements[tokenID]!.setIDMin && tokenToTrade.data.setID <= self.requirements[tokenID]!.setIDMax
                tokenToTrade.data.playID >= self.requirements[tokenID]!.playIDMin && tokenToTrade.data.playID <= self.requirements[tokenID]!.playIDMax
            }

            // get the requirements
            let requirements = self.requirements[tokenID]

            // remove the record of the trade offer
            self.requirements[tokenID] = nil

            emit MomentTraded(id: tokenID, seller: self.owner?.address)

            // get a reference to the owner's collection
            let collectionRef = self.ownerCollection.borrow<&TopShot.Collection>()!

            // deposit the new token into the owner's account
            collectionRef.deposit(token: <-tokenToTrade)

            // withdraw the offered token from the original owner's collection
            let tradedNFT <-collectionRef.withdraw(withdrawID: tokenID) as! @TopShot.NFT

            // return the token
            return <-tradedNFT
        }

        // getTradesOffered returns the dicionary of all the offered trades in the object
        pub fun getTradesOffered(): {UInt64: TradeRequirement} {
            return self.requirements
        }

        // getIDs returns an array of token IDs that are for trade
        pub fun getIDs(): [UInt64]? {
            return self.requirements.keys
        }

        // getRequirementes returns the trade requirements of a specific token up for trade
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
 
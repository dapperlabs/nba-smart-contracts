// first draft for how NBA might mint and manage molds and collections
// for the top shot NFTs

// The Topshot contract is where all the molds are stored.  All the moments
// that are minted will be able to access data for the molds they reference
// that are stored in the topshot contract

// The molds have their own mold IDs that are separate from the moment IDs
// All Moment IDs share an ID space so that none can have the same ID

// When moments are minted, they are returned by the minter.  The transaction
// has to handle the moment after that

// The top shot account will also have its own moment collection it can use to 
// hold its own moments

import NonFungibleToken from 0x01

pub contract TopShot: NonFungibleToken {

    pub struct Mold {
        // the unique ID that the mold has
        pub let id: UInt32

        // Stores all the metadata about the mold as a string mapping
        pub let metadata: {String: String}

        // the number of moments that can be minted from each quality for this mold
        pub let qualityCounts: {Int: UInt32}

        // the number of moments that have been minted from each quality mold
        // cannot be greater than the corresponding qualityCounts entry
        pub var numLeft: {Int: UInt32}

        init(id: UInt32, metadata: {String: String}, qualityCounts: {Int: UInt32}) {
            self.id = id
            self.metadata = metadata
            self.qualityCounts = qualityCounts
            self.numLeft = qualityCounts
        }
    }

    pub resource NFT: NonFungibleToken.INFT {
        // global unique moment ID
        pub let id: UInt64

        // shows metadata that is only associated with a specific NFT, and not a mold
        pub var metadata: {String:String}

        // quality identifier. Will soon be an enum
        pub var quality: Int

        // Tells which number of the Quality this moment is
        pub let placeInQuality: UInt32

        // the ID of the mold that the moment references
        pub var moldID: UInt32

        init(newID: UInt64, moldID: UInt32, quality: Int, place: UInt32) {
            pre {
                quality > 0 && quality <= 5: "Quality identifier must be 1-5!"
            }
            self.id = newID
            self.moldID = moldID
            self.quality = quality
            self.placeInQuality = place
            self.metadata = {}
        }

        pub fun getMomentMetadataField(field: String): String? {
            return TopShot.getMoldMetadataField(moldID: self.moldID, field: field)
        }

        pub fun getMomentMetadata(): {String:String}? {
            return TopShot.getMoldMetadata(moldID: self.moldID)
        }
    }

    // variable size dictionary of Mold conforming tokens
    // Mold is a struct type with an `UInt64` ID field
    pub var molds: {UInt32: Mold}

    // the ID that is used to cast molds. Every time a mold is cast,
    // moldID is assigned to the new mold's ID and then is incremented by 1.
    pub var moldID: UInt32

    // the total number of Top shot moment NFTs in existence
    // Is also used as moment IDs just like moldID
    pub var totalSupply: UInt64

    // getMoldMetadata gets a specific metadata field of a mold that is stored in this collection
    pub fun getMoldMetadataField(moldID: UInt32, field: String): String? {
        let moldOpt = self.molds[moldID]

        if let mold = moldOpt {
            return mold.metadata[field]
        } else {
            return nil
        }
    }

    pub fun getMoldMetadata(moldID: UInt32): {String:String}? {
        let moldOpt = self.molds[moldID]

        if let mold = moldOpt {
            return mold.metadata
        } else {
            return nil
        }
    }

    // getNumMomentsLeftInQuality get the number of moments left of a certain quality
    // for the specified mold ID
    pub fun getNumMomentsLeftInQuality(id: UInt32, quality: Int): UInt32 {
        if let mold = self.molds[id] {
            let numLeft = mold.numLeft[quality] ?? panic("missing numLeft!")
            return numLeft
        } else {
            return 0
        }
    }

    // getNumMintedInQuality returns the number of moments that have been minted of 
    // a certain mold ID and quality
    pub fun getNumMintedInQuality(id: UInt32, quality: Int): UInt32 {
        if let mold = self.molds[id] {
            let numLeft = mold.numLeft[quality] ?? panic("missing numLeft!")
            let qualityCount = mold.qualityCounts[quality] ?? panic("missing quality count!")
            return qualityCount - numLeft
        } else {
            return 0
        }
    }

    // getQualityTotal returns the total number of moments of a certain quality
    // that are allowed to be minted
    pub fun getQualityTotal(id: UInt32, quality: Int): UInt32 {
        if let mold = self.molds[id] {
            let qualityCount = mold.qualityCounts[quality] ?? panic("missing quality count!")
            return qualityCount
        } else {
            return 0
        }
    }

    // This is the interface that users can cast their moment Collection as
    // to allow others to deposit moments into their collection
    pub resource interface MomentCollectionPublic {
        pub var ownedNFTs: @{UInt64: NFT}
        // deposit deposits a token into the collection
        pub fun deposit(token: @NFT)
        pub fun batchDeposit(tokens: @Collection)
        // idExists checks to see if a Moment with the given ID exists in the collection
        pub fun idExists(id: UInt64): Bool
        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64]
        // getMetaData returns the metadata associated with a specific moment
        pub fun getMetaData(id: UInt64, field: String): String 
        // getMoldMetaDataField returns a field associated with a mold
        pub fun getMoldMetadataField(id: UInt64, field: String): String?
        // getMoldMetadata returns all the metadata associated with a mold
        pub fun getMoldMetadata(id: UInt64): {String:String}?
    }

    pub resource Collection: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.Metadata, MomentCollectionPublic { 
        // dictionary of Moment conforming tokens
        // Moment is a resource type with an `Int` ID field
        pub var ownedNFTs: @{UInt64: NFT}

        init() {
            self.ownedNFTs <- {}
        }

        // withdraw removes an Moment from the collection and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing Moment")
            
            return <-token
        }

        // withdraws a multiple tokens and returns them as a Collection
        pub fun batchWithdraw(ids: [UInt64]): @Collection {
            var i = 0
            var batchCollection: @Collection <- create Collection()

            while i < ids.length {
                batchCollection.deposit(token: <-self.withdraw(withdrawID: ids[i]))

                i = i + 1
            }
            return <-batchCollection
        }

        // deposit takes a Moment and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @NFT) {
            // add the new token to the dictionary
            let oldToken <- self.ownedNFTs[token.id] <- token
            destroy oldToken
        }

        // takes a Collection object as an argument
        // and deposits each contained NFT into this collection
        pub fun batchDeposit(tokens: @Collection) {
            var i = 0
            let keys = tokens.getIDs()

            while i < keys.length {
                self.deposit(token: <-tokens.withdraw(withdrawID: keys[i]))

                i = i + 1
            }
            destroy tokens
        }

        // idExists checks to see if a Moment with the given ID exists in the collection
        pub fun idExists(id: UInt64): Bool {
            if self.ownedNFTs[id] != nil {
                return true
            }

            return false
        }

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        // Gets metadata associated with the specific Moment NFT
        //
        pub fun getMetaData(id: UInt64, field: String): String {
            let token <- self.ownedNFTs.remove(key: id) ?? panic("No NFT!")
            
            let dataOpt = token.metadata[field]

            let oldToken <- self.ownedNFTs[id] <- token
            destroy oldToken

            if let data = dataOpt {
                return data
            } else {
                return "None"
            }
        }

        // getMomentMetadataField gets a specific mold metadata field of a moment in the collection
        //
        pub fun getMoldMetadataField(id: UInt64, field: String): String? {
            let moment <- self.ownedNFTs[id] <- nil

            let field = moment?.getMomentMetadataField(field: field) ?? panic("moment doesn't exist!")

            let oldMoment <- self.ownedNFTs[id] <- moment
            destroy oldMoment

            return field
        }

        // getMomentMetadata gets all the metadata of a certain moment in the collection
        pub fun getMoldMetadata(id: UInt64): {String:String}? {
            let moment <- self.ownedNFTs[id] <- nil

            let metadata = moment?.getMomentMetadata() ?? panic("moment doesn't exist!")

            let oldMoment <- self.ownedNFTs[id] <- moment
            destroy oldMoment
            
            return metadata
        }

        destroy() {
            destroy self.ownedNFTs
        }
    }

    pub fun createEmptyCollection(): @Collection {
        return <-create Collection()
    }

    // MoldCaster is a resource that the user who has admin access to the Topshot 
    // contract will store in their account 
    // this ensures that they are the only ones who can cast molds and mint moments
    pub resource Admin {
        // castMold casts a mold struct and stores it in the dictionary
        // for the molds
        // the mold ID must be unused
        // returns the ID the new mold
        pub fun castMold(metadata: {String: String}, qualityCounts: {Int: UInt32}): UInt32 {
            pre {
                qualityCounts.length > 8: "Wrong number of qualities!"
                metadata.length != 0: "Wrong amount of metadata!"
            }
            var newMold = Mold(id: TopShot.moldID, metadata: metadata, qualityCounts: qualityCounts)

            TopShot.molds[TopShot.moldID] = newMold

            // increment the ID so that it isn't used again
            TopShot.moldID = TopShot.moldID + UInt32(1)

            return TopShot.moldID - UInt32(1)
        }

        // mintMoment mints a new moment and returns the newly minted moment
        pub fun mintMoment(moldID: UInt32, quality: Int): @NFT {
            pre {
                // check to see if any more moments of this quality are allowed to be minted
                TopShot.getNumMomentsLeftInQuality(id: moldID, quality: quality) > UInt32(0): "All the moments of this quality have been minted!"
            }

            // update the number left in the quality that are allowed to be minted
            let mold = TopShot.molds[moldID] ?? panic("invalid mold ID")
            var numLeft = mold.numLeft[quality] ?? panic("invalid quality ID")
            mold.numLeft[quality] = numLeft - UInt32(1)
            TopShot.molds[moldID] = mold

            // gets this moment's place in the moments for this quality
            let placeInQuality = TopShot.getNumMintedInQuality(id: moldID, quality: quality)

            // mint the new moment
            let newMoment: @NFT <- create NFT(newID: TopShot.totalSupply, 
                                                    moldID: moldID, 
                                                    quality: quality, 
                                                    place: placeInQuality)


            TopShot.totalSupply = TopShot.totalSupply + UInt64(1)

            return <-newMoment
        }

        // batchMintMoment mints an arbitrary quantity of moments and returns all of them in
        // a new moment Collection
        pub fun batchMintMoment(moldID: UInt32, quality: Int, quantity: UInt64): @Collection {
            let newCollection <- create Collection()

            var i: UInt64 = 0
            while i < quantity {
                newCollection.deposit(token: <-self.mintMoment(moldID: moldID, quality: quality))
            }

            return <-newCollection
        }

        pub fun createAdmin(): @Admin {
            return <-create Admin()
        }
    }

    init() {
        // initialize the fields
        self.molds = {}
        self.moldID = 0
        self.totalSupply = 0

        // Create a new collection
        let oldCollection <- self.account.storage[Collection] <- create Collection()
        destroy oldCollection

        self.account.storage[&Collection] = &self.account.storage[Collection] as Collection
        self.account.published[&MomentCollectionPublic] = &self.account.storage[Collection] as MomentCollectionPublic

        // Create a new Admin resource and store it in account storage
        let oldAdmin <- self.account.storage[Admin] <- create Admin()
        destroy oldAdmin

        // Create a private reference to the Admin resource and store it in private account storage
        self.account.storage[&Admin] = &self.account.storage[Admin] as Admin
    }

}
 

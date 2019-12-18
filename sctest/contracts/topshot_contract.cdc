// first draft for how NBA might mint and manage molds and collections
// for the top shot NFTs

// The NBA top shot account has a top level collection that holds all the molds
// the molds are structs that contain moment metadata

// The mold collection is where all casting of molds happens.  It is also
// where all the molds are stored.
// The Mold Collection will not store moments

// the top shot account also has a resource that is used for minting moments, MomentFactory
// when it mints a moment, it gives the moment a reference to the mold collection
// that holds the molds so a moment can always be able to fetch its metadata 
// by making calls to the mold collection

// The molds have their own mold IDs that are separate from the moment IDs
// All Moment IDs share an ID space so that none can have the same ID

// When moments are minted, they are 
// sent directly to the account that will own the Moment
// this account could be the top shot moment collection or a user's collection

// The top shot account will also have its own moment collection it can use to 
// hold its own moments

// TODO: Want to make some generic way of reading a field from a moment that 
// is inherited from a mold. so we don't have to make a getter for each field in the 
// MoldCollection.  Could also store all the fields as a dictionary keyed by the name
// of each field.
// Eventually want to be able to do something like:
// collectionRef.moments[2].name
// and have it use the reference internally to get the name from the mold

pub contract interface TopShotPublic {
    pub struct Mold {}
    pub var molds: {Int: Mold}
    pub var moldID: Int
    pub var momentID: Int
    pub fun getMoldMetadataField(moldID: Int, field: String): String?
    pub fun getNumMomentsLeftInQuality(id: Int, quality: Int): Int {
        pre {
            quality > 0 && quality <= 5: "Quality needs to be 1-5"
            id > 0: "ID needs to be positive!"
        }
    }
    pub fun getNumMintedInQuality(id: Int, quality: Int): Int {
        pre {
            quality > 0 && quality <= 5: "Quality needs to be 1-5"
            id > 0: "ID needs to be positive!"
        }
    }
}

pub contract TopShot: TopShotPublic {
    pub struct Mold {
        pub let id: Int  // the unique ID that the mold has

        // Stores all the metadata about the mold as a string mapping
        pub let metadata: {String: String}

        pub let qualityCounts: {Int: Int}  // the number of moments that can be minted from each quality for this mold

        pub var numLeft: {Int: Int}      // the number of moments that have been minted from each quality mold
                                            // cannot be greater than the corresponding qualityCounts

        init(id: Int, metadata: {String: String} , qualityCounts: {Int: Int}) {
            self.id = id
            self.metadata = metadata
            self.qualityCounts = qualityCounts
            self.numLeft = qualityCounts
        }

        // called when a moment is minted
        pub fun updateNumLeft(quality: Int) {
            var numLeft = self.numLeft[quality] ?? panic("missing quality count!")
            
            numLeft = numLeft - 1
            
            self.numLeft[quality] = numLeft
        }
    }

    pub resource Moment {
        // global unique moment ID
        pub let id: Int

        // quality identifier. Will soon be an enum
        pub var quality: Int

        // Tells which number of the Quality this moment is
        pub let placeInQuality: Int

        // the ID of the mold that the moment references
        pub var moldID: Int

        // reference to the NBA Mold Collection that holds the mold
        //pub let contractReference: &TopShotPublic

        init(newID: Int, moldID: Int, quality: Int, place: Int) { //, reference: &TopShotPublic) {
            pre {
                newID > 0: "MomentID must be a positive integer!"
                moldID > 0: "MoldID must be a positive integer!"
                quality > 0 && quality <= 5: "Quality identifier must be 1-5!"
            }
            self.id = newID
            self.moldID = moldID
            self.quality = quality
            self.placeInQuality = place
            //self.contractReference = reference
        }
    }

    // variable size dictionary of Mold conforming tokens
    // Mold is a struct type with an `Int` ID field
    pub var molds: {Int: Mold}

    // the ID that is used to cast molds
    pub var moldID: Int

    // the ID that is used to mint unique moments
    pub var momentID: Int

    // castMold casts a mold struct and stores it in the dictionary
    // for the molds
    // the mold ID must be unused
    pub fun castMold(metadata: {String: String}, qualityCounts: {Int: Int}) {
        pre {
            qualityCounts.length == 5: "Wrong number of qualities!"
            metadata.length != 0: "Wrong amount of metadata!"
        }
        var newMold = Mold(id: self.moldID, metadata: metadata, qualityCounts: qualityCounts)

        self.molds[self.moldID] = newMold

        // increment the ID so that it isn't used again
        self.moldID = self.moldID + 1
    }

    // getMoldMetadata gets a specific metadata field of a mold that is stored in this collection
    pub fun getMoldMetadataField(moldID: Int, field: String): String? {
        let moldOpt = self.molds[moldID]

        if let mold = moldOpt {
            return mold.metadata[field]
        } else {
            return nil
        }
    }

    // getNumMomentsLeftInQuality get the number of moments left of a certain quality
    // for the specified mold ID
    pub fun getNumMomentsLeftInQuality(id: Int, quality: Int): Int {
        if let mold = self.molds[id] {
            let numLeft = mold.numLeft[quality] ?? panic("missing numLeft!")
            return numLeft
        } else {
            return 0
        }
    }

    // getNumMintedInQuality returns the number of moments that have been minted of 
    // a certain mold ID and quality
    pub fun getNumMintedInQuality(id: Int, quality: Int): Int {
        if let mold = self.molds[id] {
            let numLeft = mold.numLeft[quality] ?? panic("missing numLeft!")
            let qualityCount = mold.qualityCounts[quality] ?? panic("missing quality count!")
            return qualityCount - numLeft
        } else {
            return -1
        }
    }

    // getQualityTotal returns the total number of moments of a certain quality
    // that are allowed to be minted
    pub fun getQualityTotal(id: Int, quality: Int): Int {
        if let mold = self.molds[id] {
            let qualityCount = mold.qualityCounts[quality] ?? panic("missing quality count!")
            return qualityCount
        } else {
            return -1
        }
    }

    // mintMoment mints a moment NFT based off of a mold that is stored in the collection
    // the moment ID must be unused
    pub fun mintMoment(moldID: Int, quality: Int, recipient: &MomentReceiver) {
        pre {
            // check to see if any more moments of this quality are allowed to be minted
            self.getNumMomentsLeftInQuality(id: moldID, quality: quality) <= 0: "All the moments of this quality have been minted!"
        }

        // update the number left in the quality that are allowed to be minted
        self.molds[moldID]?.updateNumLeft(quality: quality)

        // gets this moment's place in the moments for this quality
        let placeInQuality = self.getNumMintedInQuality(id: moldID, quality: quality)

        // mint the new moment
        var newMoment: @Moment <- create Moment(newID: self.momentID, 
                                                moldID: moldID, 
                                                quality: quality, 
                                                place: placeInQuality)
                                                //reference: &self)
        
        // deposit the moment in the owner's account
        recipient.deposit(token: <-newMoment)

        // update the moment IDs so they can't be reused
        self.momentID = self.momentID + 1
    }

    pub resource interface MomentReceiver {
        pub var moments: @{Int: Moment}
        // deposit deposits a token into the collection
        pub fun deposit(token: @Moment)
        // idExists checks to see if a Moment with the given ID exists in the collection
        pub fun idExists(tokenID: Int): Bool
        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [Int]
    }

    pub resource MomentCollection: MomentReceiver { 
        // dictionary of Moment conforming tokens
        // Moment is a resource type with an `Int` ID field
        pub var moments: @{Int: Moment}

        init() {
            self.moments <- {}
        }

        // withdraw removes an Moment from the collection and moves it to the caller
        pub fun withdraw(tokenID: Int): @Moment {
            pre {
                tokenID > 0: "Token ID must be positive!"
            }
            let token <- self.moments.remove(key: tokenID) ?? panic("missing Moment")
            
            return <-token
        }

        // deposit takes a Moment and adds it to the collections dictionary
        // and adds the ID to the id array
        pub fun deposit(token: @Moment) {
            // add the new token to the dictionary
            let oldToken <- self.moments[token.id] <- token
            destroy oldToken
        }

        // transfer takes a reference to another user's Moment collection,
        // takes the Moment out of this collection, and deposits it
        // in the reference's collection
        pub fun transfer(recipient: &MomentReceiver, tokenID: Int) {
            // remove the token from the dictionary get the token from the optional
            let token <- self.withdraw(tokenID: tokenID)

            // deposit it in the recipient's account
            recipient.deposit(token: <-token)
        }

        // idExists checks to see if a Moment with the given ID exists in the collection
        pub fun idExists(tokenID: Int): Bool {
            if self.moments[tokenID] != nil {
                return true
            }

            return false
        }

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [Int] {
            return self.moments.keys
        }

        destroy() {
            destroy self.moments
        }
    }

    init() {
        self.molds = {}
        self.moldID = 1
        self.momentID = 1

        let oldMomentCollection <- self.account.storage[MomentCollection] <- create MomentCollection()
        destroy oldMomentCollection

        self.account.storage[&MomentCollection] = &self.account.storage[MomentCollection] as MomentCollection
        self.account.published[&MomentReceiver] = &self.account.storage[MomentCollection] as MomentReceiver
    }

}
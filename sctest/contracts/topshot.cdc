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

// Note: All state changing functions will panic if an invalid argument is
// provided, if a mistake happens, or if certain states aren't allowed.  Functions
// that don't modify state will simply return 0 or nil and those cases need
// to be handled by the caller.

//import NonFungibleToken from 0x01

access(all) contract TopShot { //: NonFungibleToken {

    access(all) event MoldCasted(id: UInt32, qualityCounts: [UInt32])
    access(all) event MomentMinted(id: UInt64, moldID: UInt32)
    access(all) event ContractInitialized()
    access(all) event Withdraw(id: UInt64)
    access(all) event Deposit(id: UInt64)

    access(all) struct Mold {
        // the unique ID that the mold has
        access(all) let id: UInt32

        // Stores all the metadata about the mold as a string mapping
        access(all) let metadata: {String: String}

        // the number of moments that can be minted from each quality for this mold
        access(all) let qualityCounts: {Int: UInt32}

        // the number of moments that have been minted from each quality mold
        // cannot be greater than the corresponding qualityCounts entry
        access(account) var numLeft: {Int: UInt32}

        // shows if a certain quality of this mold can be minted or not
        access(account) var canBeMinted: {Int: Bool}

        init(id: UInt32, metadata: {String: String}, counts: [UInt32]) {
            pre {
                counts.length == 8: "Wrong number of qualities!"
                metadata.length != 0: "Wrong amount of metadata!"
            }
            self.id = id
            self.metadata = metadata
            self.qualityCounts = {1: counts[0], 2: counts[1], 3: counts[2], 4: counts[3], 5: counts[4], 6: counts[5], 7: counts[6], 8: counts[7]}
            self.numLeft = self.qualityCounts
            self.canBeMinted = {1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true}
        }
    }

    access(all) resource NFT { //: NonFungibleToken.INFT {
        // global unique moment ID
        access(all) let id: UInt64

        // shows metadata that is only associated with a specific NFT, and not a mold
        access(all) var metadata: {String:String}

        // quality identifier. Will soon be an enum
        access(all) var quality: Int

        // Tells which number of the Quality this moment is
        access(all) let placeInQuality: UInt32

        // the ID of the mold that the moment references
        access(all) var moldID: UInt32

        init(newID: UInt64, moldID: UInt32, quality: Int, place: UInt32) {
            pre {
                quality > 0 && quality <= 8: "Quality identifier must be 1-5!"
            }
            self.id = newID
            self.moldID = moldID
            self.quality = quality
            self.placeInQuality = place
            self.metadata = {}
        }
    }

    // variable size dictionary of Mold conforming tokens
    // Mold is a struct type with an `UInt64` ID field
    access(all) var molds: {UInt32: Mold}

    // the ID that is used to cast molds. Every time a mold is cast,
    // moldID is assigned to the new mold's ID and then is incremented by 1.
    access(all) var moldID: UInt32

    // the total number of Top shot moment NFTs in existence
    // Is also used as moment IDs for minting just like moldID
    access(all) var totalSupply: UInt64

    // getNumMomentsLeftInQuality get the number of moments left of a certain quality
    // for the specified mold ID
    access(all) fun getNumMomentsLeftInQuality(id: UInt32, quality: Int): UInt32 {
        if let mold = self.molds[id] {
            if let numLeft = mold.numLeft[quality] {
                return numLeft
            }
            else {
                return 0
            }
        } else {
            return 0
        }
    }

    // getNumMintedInQuality returns the number of moments that have been minted of 
    // a certain mold ID and quality
    // All the `return 0` lines are situations when the caller provided an incorrect
    // paramter value so it returns 0 to show there have been none minted.
    access(all) fun getNumMintedInQuality(id: UInt32, quality: Int): UInt32 {
        if let mold = self.molds[id] {
            if let numLeft = mold.numLeft[quality] {
                if let qualityCount = mold.qualityCounts[quality] {
                    return qualityCount - numLeft
                } else {
                    return 0
                }
            } else {
                return 0
            }
        } else {
            return 0
        }
    }

    // mintingAllowed Returns a boolean that indicates if minting is allowed
    // for a certain mold and quality.
    //
    access(all) fun mintingAllowed(id: UInt32, quality: Int): Bool {
        if let mold = self.molds[id] {
            if let canBeMinted = mold.canBeMinted[quality] {
                return canBeMinted
            } else {
                return false
            }
        } else {
            return false
        }
    }

    // This is the interface that users can cast their moment Collection as
    // to allow others to deposit moments into their collection
    access(all) resource interface MomentCollectionPublic {
        access(all) fun deposit(token: @NFT)
        access(all) fun batchDeposit(tokens: @Collection)
        access(all) fun getIDs(): [UInt64]
        access(all) fun getMoldID(id: UInt64): UInt32
        access(all) fun getQuality(id: UInt64): Int
        access(all) fun getPlaceInQuality(id: UInt64): UInt32
        access(all) fun getMetaData(id: UInt64): {String: String}
    }

    access(all) resource Collection: MomentCollectionPublic { //: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.Metadata, MomentCollectionPublic { 
        // Dictionary of Moment conforming tokens
        // NFT is a resource type with a UInt64 ID field
        access(account) var ownedNFTs: @{UInt64: NFT}

        init() {
            self.ownedNFTs <- {}
        }

        // withdraw removes an Moment from the collection and moves it to the caller
        access(all) fun withdraw(withdrawID: UInt64): @NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing Moment")

            emit Withdraw(id: token.id)
            
            return <-token
        }

        // batchWithdraw withdraws multiple tokens and returns them as a Collection
        access(all) fun batchWithdraw(ids: [UInt64]): @Collection {
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
        access(all) fun deposit(token: @NFT) {
            let id = token.id
            // add the new token to the dictionary
            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id)

            destroy oldToken
        }

        // batchDeposit takes a Collection object as an argument
        // and deposits each contained NFT into this collection
        access(all) fun batchDeposit(tokens: @Collection) {
            var i = 0
            let keys = tokens.getIDs()

            while i < keys.length {
                self.deposit(token: <-tokens.withdraw(withdrawID: keys[i]))

                i = i + 1
            }
            destroy tokens
        }

        // getIDs returns an array of the IDs that are in the collection
        access(all) fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        access(all) fun getMoldID(id: UInt64): UInt32 {
            return self.ownedNFTs[id]?.moldID ?? panic("No moment!")
        }

        access(all) fun getQuality(id: UInt64): Int {
            return self.ownedNFTs[id]?.quality ?? panic("No moment!")
        }

        access(all) fun getPlaceInQuality(id: UInt64): UInt32 {
            return self.ownedNFTs[id]?.placeInQuality ?? panic("No moment!")
        }

        access(all) fun getMetaData(id: UInt64): {String: String} {
            return TopShot.molds[self.getMoldID(id: id)]?.metadata ?? panic("No mold!")
        }

        destroy() {
            destroy self.ownedNFTs
        }
    }

    access(all) fun createEmptyCollection(): @Collection {
        return <-create Collection()
    }

    // Admin is a resource that the user who has admin access to the Topshot 
    // contract will store in their account 
    // this ensures that they are the only ones who can cast molds and mint moments
    access(all) resource Admin {
        // castMold casts a mold struct and stores it in the dictionary for the molds
        // the mold ID must be unused
        // returns the ID of the new mold
        access(all) fun castMold(metadata: {String: String}, qualityCounts: [UInt32]): UInt32 {
            pre {
                qualityCounts.length == 8: "Quality Counts must have eight elementS"
            }
            // Create the new Mold
            var newMold = Mold(id: TopShot.moldID, metadata: metadata, counts: qualityCounts)

            // Store it in the contract storage
            TopShot.molds[TopShot.moldID] = newMold

            // increment the ID so that it isn't used again
            TopShot.moldID = TopShot.moldID + UInt32(1)

            emit MoldCasted(id: TopShot.moldID - UInt32(1), qualityCounts: qualityCounts)

            return TopShot.moldID - UInt32(1)
        }

        // sets minting allowed to false
        // cannot be reversed
        access(all) fun disallowMinting(moldID: UInt32, quality: Int) {
            pre {
                quality > 0 && quality <= 8: "Quality must be an integer between 1 and 5"
            }
            if let mold = TopShot.molds[moldID] {
                mold.canBeMinted[quality] = false
                TopShot.molds[moldID] = mold
            } else {
                panic("Incorrect mold specified!")
            }
        }

        // mintMoment mints a new moment and returns the newly minted moment
        access(all) fun mintMoment(moldID: UInt32, quality: Int): @NFT {
            pre {
                // check to see if any more moments of this quality are allowed to be minted
                TopShot.getNumMomentsLeftInQuality(id: moldID, quality: quality) > UInt32(0): "All the moments of this quality have been minted!"
                TopShot.mintingAllowed(id: moldID, quality: quality): "This mold quality does not allow new minting!"
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

            emit MomentMinted(id: TopShot.totalSupply, moldID: moldID)

            TopShot.totalSupply = TopShot.totalSupply + UInt64(1)

            return <-newMoment
        }

        // batchMintMoment mints an arbitrary quantity of moments all of the same ID
        // and quality and returns all of them in a new moment Collection
        access(all) fun batchMintMoment(moldID: UInt32, quality: Int, quantity: UInt64): @Collection {
            let newCollection <- create Collection()

            var i: UInt64 = 0
            while i < quantity {
                newCollection.deposit(token: <-self.mintMoment(moldID: moldID, quality: quality))
            }

            return <-newCollection
        }

        // Creates a new admin resource that can be transferred to another account
        access(all) fun createAdmin(): @Admin {
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

        // Create a private reference to the Collection and store it in private account storage
        self.account.storage[&Collection] = &self.account.storage[Collection] as &Collection

        // Create a safe, public reference to the Collection and store it in public reference storage
        self.account.published[&MomentCollectionPublic] = &self.account.storage[Collection] as &MomentCollectionPublic

        // Create a new Admin resource and store it in account storage
        let oldAdmin <- self.account.storage[Admin] <- create Admin()
        destroy oldAdmin

        // Create a private reference to the Admin resource and store it in private account storage
        self.account.storage[&Admin] = &self.account.storage[Admin] as &Admin

        emit ContractInitialized()
    }

}
 
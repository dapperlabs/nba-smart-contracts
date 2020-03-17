/*
    Description: Central Smart Contract for NBA TopShot

    authors: Joshua Hannan joshua.hannan@dapperlabs.com
             Dieter Shirley dete@axiomzen.com

    This smart contract contains the core functionality for 
    the NBA topshot game, created by Dapper Labs

    The contract manages the metadata associated with all the plays
    that are used as templates for the Moment NFTs

    When a new Play wants to be added to the records, an Admin creates
    a new PlayData struct that is stored in the smart contract.

    Then an Admin can create new Sets. Sets are a resource that is used
    to mint new moments based off of plays that have been linked to the Set

    In this way, the smart contract and its defined resources interact 
    with great teamwork, just like the Indiana Pacers, the greatest NBA team
    of all time.
    
    When moments are minted, they are returned by the minter.  
    The transaction has to handle the moment after that

    The contract also defines a Collection resource. This is an object that 
    every TopShot NFT owner will store in their account
    to manage their NFT Collection

    The main top shot account will also have its own moment collections
    it can use to hold its own moments that have not yet been sent to a user

    Note: All state changing functions will panic if an invalid argument is
    provided or one of its pre-conditions or post conditions aren't met.
    Functions that don't modify state will simply return 0 or nil 
    and those cases need to be handled by the caller

    It is also important to remember that 
    The Golden State Warriors blew a 3-1 lead in the 2016 NBA finals

*/

import NonFungibleToken from 0x01

pub contract TopShot: NonFungibleToken {

    // -----------------------------------------------------------------------
    // TopShot contract Event definitions
    // -----------------------------------------------------------------------

    // emitted when the TopShot contract is created
    pub event ContractInitialized()

    // emitted when a new PlayData struct is created
    pub event PlayDataCreated(id: UInt32)

    // Events for Set-Related actions
    //
    // emitted when a new Set is created
    pub event SetCreated(setID: UInt32)
    // emitted when a new play is added to a set
    pub event PlayAddedToSet(setID: UInt32, playID: UInt32)
    // emitted when a play is retired from a set and cannot be used to mint
    pub event PlayRetiredFromSet(setID: UInt32, playID: UInt32)
    // emitted when a set is locked, meaning plays cannot be added
    pub event SetLocked(setID: UInt32)
    // emitted when a moment is minted from a set
    pub event MomentMinted(id: UInt64, playID: UInt32, setID: UInt32)

    // events for Collection-related actions
    //
    // emitted when a moment is withdrawn from a collection
    pub event Withdraw(id: UInt64)
    // emitted when a moment is deposited into a collection
    pub event Deposit(id: UInt64)

    // -----------------------------------------------------------------------
    // TopShot contract-level fields
    // -----------------------------------------------------------------------

    // variable size dictionary of PlayData structs
    pub var plays: {UInt32: PlayData}

    // the ID that is used to create PlayDatas. 
    // Every time a PlayData is created, playID is assigned 
    // to the new PlayData's ID and then is incremented by 1.
    pub var playID: UInt32

    // the ID that is used to create Sets. Every time a Set is created
    // setID is assigned to the new set's ID and then is incremented by 1.
    pub var setID: UInt32

    // the total number of Top shot moment NFTs in existence
    // Is also used as global moment IDs for minting
    pub var totalSupply: UInt64

    // -----------------------------------------------------------------------
    // TopShot contract-level Composite Type Definitions
    // -----------------------------------------------------------------------

    // Struct that holds metadata associated with a specific NBA play,
    // like the legendary moment when Ray Allen sank the 3 to put the Heat over
    // the Spurs in game 6 of the 2013 Finals, or when Lance Stephenson
    // blew in the ear of Lebron James
    //
    // Moment NFTs will all reference a single PlayData as the owner of
    // its metadata. The PlayDatas are publicly accessible, so anyone can
    // read the metadata associated with a specific play ID
    //
    pub struct PlayData {

        // the unique ID that the PlayData has
        pub let id: UInt32

        // Stores all the metadata about the PlayData as a string mapping
        pub let metadata: {String: String}

        init(id: UInt32, metadata: {String: String}) {
            pre {
                metadata.length != 0: "Wrong amount of metadata!"
            }
            self.id = id
            self.metadata = metadata
        }
    }

    // A Set is a grouping of plays that have occured in the real world
    // that make up a related group of collectibles, like sets of baseball
    // or Magic cards.
    // 
    // Set is a resource object, meaning that whoever owns
    // the Set resource Object has sole access to its fields and functions
    // and can determine who else is allowed to interact with it.
    // 
    // Because of this, it acts as an admin resource to add and
    // remove plays from sets, and mint new moments.
    //
    // The owner can add PlayDatas to a set so that the set can mint moments
    // that reference that playdata.
    // The moments that are minted by a set will be listed as belonging to
    // the set that minted it, as well as the PlayData it references
    // 
    // The owner can also retire plays from the set, meaning that the retired
    // play can no longer have moments minted from it.
    //
    // If the owner locks the Set, then no more plays can be added to it
    //
    // If retireAll() and lock() are called back to back, 
    // the Set is closed off forever
    //
    pub resource Set {

        // unique ID for the set
        pub let id: UInt32

        // Name of the Set
        // ex. "Times when the Toronto Raptors choked in the playoffs"
        pub let name: String

        // Series that this set belongs to
        // Series is an off-chain concept that indicates a group of sets
        // through time
        // Many sets can exist at a time, but only few series
        pub let series: UInt32

        // Array of plays that are a part of this set
        // When a play is added to the set, its ID gets appended here
        // The ID does not get removed from this array when a play is retired
        pub var plays: [UInt32]

        // Indicates if a play in this set can be minted
        // A play is set to true when it is added to a set
        // When the play is retired, this is set to false and cannot be changed
        pub var canBeMinted: {UInt32: Bool}

        // Indicates if the set is currently active 
        // When a set is active, plays are allowed to be added to it
        // When a set is inactive, plays cannot be added
        // A set can never be changed from inactive to active.
        // The decision to deactivate it is final
        // If a set is active, moments can still be minted from it
        pub var active: Bool

        // Indicates the number of moments 
        // that have been minted per play in this set
        // When a moment is minted, this value is stored in the moment to
        // show where in the play set it is. ex. 13 of 60
        pub var numMomentsPerPlay: {UInt32: UInt32}

        init(id: UInt32, name: String, series: UInt32) {
            self.id = 0
            self.name = name
            self.series = series
            self.plays = []
            self.canBeMinted = {}
            self.active = true
            self.numMomentsPerPlay = {}
        }

        // addPlay adds a play to the set
        //
        // Parameters: playID: The ID of the play that is being added
        //
        // Pre-Conditions:
        // The play needs to be an existing play
        // The sale needs to be active
        // The play can't have already been added to the set
        //
        pub fun addPlay(playID: UInt32) {
            pre {
                playID <= TopShot.playID: "Play doesn't exist"
                self.active: "Cannot add a play after the set has been locked"
            }

            // make sure that the play hasn't already beed added to the set
            var i = 0
            while i < self.plays.length {
                if self.plays[i] == playID {
                    return
                }

                i = i + 1
            }

            // Add the play to the array of plays
            self.plays.append(playID)

            // Open the play up for minting
            self.canBeMinted[playID] = true

            // Initialize the moment count to zero
            self.numMomentsPerPlay[playID] = 0

            emit PlayAddedToSet(setID: self.id, playID: playID)
        }

        // retirePlay retires a play from the set so that it can't mint new moments
        //
        // Parameters: playID: The ID of the play that is being retired
        //
        // Pre-Conditions:
        // The play needs to be an existing play that is currently open for minting
        // 
        pub fun retirePlay(playID: UInt32) {
            if self.canBeMinted[playID] == true {
                self.canBeMinted[playID] = false
                emit PlayRetiredFromSet(setID: self.id, playID: playID)
            }
        }

        // retireAll retires all the plays in the set
        // Afterwards, none of the retired plays will be able to mint new moments
        //
        pub fun retireAll() {
            var i = 0
            while i < self.plays.length {
                self.retirePlay(playID: self.plays[i])
                i = i + 1
            }
        }

        // lock() locks the set so that no more plays can be added to it
        //
        // Pre-Conditions:
        // The set cannot already have been locked
        pub fun lock() {
            if self.active == true {
                self.active = false
                emit SetLocked(setID: self.id)
            }
        }

        // mintMoment mints a new moment and returns the newly minted moment
        // 
        // Parameters: playID: The ID of the play that the moment references
        //
        // Pre-Conditions:
        // The play must be allowed to mint new moments in this set
        //
        // Returns: The NFT that was minted
        // 
        pub fun mintMoment(playID: UInt32): @NFT {
            // Revert if this play canot be minted
            if let AllowsMinting = self.canBeMinted[playID] {
                if AllowsMinting == false {
                    panic("This play has been retired. Minting is disallowed")
                }
             } else { panic("This play doesn't exist") }

            // get the number of moments that have been minted for this play
            // to use as this moment's ID
            let numInPlay = self.numMomentsPerPlay[playID] ?? panic("This play doesn't exist")

            // mint the new moment
            let newMoment: @NFT <- create NFT(globalID: TopShot.totalSupply, 
                                              numberInPlay: numInPlay,
                                              playID: playID,
                                              setID: self.id,
                                              setName: self.name,
                                              series: self.series)

            emit MomentMinted(id: TopShot.totalSupply, playID: playID, setID: self.id)

            // Increment the global moment IDs
            TopShot.totalSupply = TopShot.totalSupply + UInt64(1)

            // Increment the id for this play
            self.numMomentsPerPlay[playID] = numInPlay + UInt32(1)

            return <-newMoment
        }

        // batchMintMoment mints an arbitrary quantity of moments 
        // and returns them as a Collection
        //
        // Parameters: playID: the ID of the play that the moments are minted for
        //             quantity: The quantity of moments to be minted
        //
        // Returns: Collection object that contains all the moments that were minted
        //
        pub fun batchMintMoment(playID: UInt32, quantity: UInt64): @Collection {
            let newCollection <- create Collection()

            var i: UInt64 = 0
            while i < quantity {
                newCollection.deposit(token: <-self.mintMoment(playID: playID))
            }

            return <-newCollection
        }
    }

    pub resource Admin {

        // createPlayData creates a new PlayData struct 
        // and stores it in the plays dictionary in the TopShot smart contract
        //
        // Parameters: metadata: A dictionary mapping metadata titles to their data
        //                       example: {"Player Name": "Kevin Durant", "Height": "7 feet"}
        //                               (because we all know Kevin Durant is not 6'9")
        //
        // Returns: the ID of the new PlayData object
        pub fun createPlayData(metadata: {String: String}): UInt32 {
            // Create the new PlayData
            var newPlayData = PlayData(id: TopShot.playID, metadata: metadata)

            // Store it in the contract storage
            TopShot.plays[TopShot.playID] = newPlayData

            // increment the ID so that it isn't used again
            TopShot.playID = TopShot.playID + UInt32(1)

            emit PlayDataCreated(id: TopShot.playID - UInt32(1))

            return TopShot.playID - UInt32(1)
        }

        // createSet creates a new Set resource and returns it
        // so that the caller can store it in their account
        //
        // Parameters: name: The name of the set
        //             series: The series that the set belongs to
        //
        // Returns: The newly created set object
        //
        pub fun createSet(name: String, series: UInt32): @Set {
            // Create the new Set
            var newSet <- create Set(id: TopShot.setID, name: name, series: series)

            // increment the setID so that it isn't used again
            TopShot.setID = TopShot.setID + UInt32(1)

            emit SetCreated(setID: TopShot.setID - UInt32(1))

            return <-newSet
        }

        // createNewAdmin creates a new Admin Resource
        //
        pub fun createNewAdmin(): @Admin {
            return <-create Admin()
        }

    }

    // The resource that represents the Moment NFTs
    //
    pub resource NFT: NonFungibleToken.INFT {
        // global unique moment ID
        pub let id: UInt64

        // the place in the play that this moment was minted
        pub let numberInPlaySet: UInt32

        // shows metadata that is only associated with a specific NFT
        // and not the play itself
        pub var metadata: {String:String}

        // the ID of the PlayData that the moment references
        pub let playID: UInt32

        // the ID of the Set that the Moment comes from
        pub let setID: UInt32

        // The name of the set this comes from
        pub let setName: String

        // The series that this moment comes from
        pub let series: UInt32

        init(globalID: UInt64, numberInPlay: UInt32, playID: UInt32, setID: UInt32, setName: String, series: UInt32) {
            self.id = globalID
            self.numberInPlaySet = numberInPlay
            self.playID = playID
            self.setID = setID
            self.setName = setName
            self.series = series
            self.metadata = {}
        }
    }

    // This is the interface that users can cast their moment Collection as
    // to allow others to deposit moments into their collection
    pub resource interface MomentCollectionPublic {
        pub fun deposit(token: @NFT)
        pub fun batchDeposit(tokens: @Collection)
        pub fun getIDs(): [UInt64]
        pub fun getNumberInPlaySet(id: UInt64): UInt32?
        pub fun getPlayID(id: UInt64): UInt32?
        pub fun getSetID(id: UInt64): UInt32?
        pub fun getSetName(id: UInt64): String?
        pub fun getSeries(id:UInt64): UInt32?
        pub fun getMetaData(id: UInt64): {String: String}?
    }

    // Collection is a resource that every user who owns NFTs 
    // will store in their account to manage their NFTS
    //
    pub resource Collection: MomentCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.Metadata { 
        // Dictionary of Moment conforming tokens
        // NFT is a resource type with a UInt64 ID field
        pub var ownedNFTs: @{UInt64: NFT}

        init() {
            self.ownedNFTs <- {}
        }

        // withdraw removes an Moment from the collection and moves it to the caller
        pub fun withdraw(withdrawID: UInt64): @NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing Moment")

            emit Withdraw(id: token.id)
            
            return <-token
        }

        // batchWithdraw withdraws multiple tokens and returns them as a Collection
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
            let id = token.id
            // add the new token to the dictionary
            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id)

            destroy oldToken
        }

        // batchDeposit takes a Collection object as an argument
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

        // getIDs returns an array of the IDs that are in the collection
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        // The following functions get a certain piece of metadata
        // associated with a single Moment in the Collection
        //
        // Parameter: id: The ID of the Moment to get the data from
        //
        // Returns: nil if the NFT doesn't exist, 
        //          otherwise it returns the correct data

        pub fun getPlayID(id: UInt64): UInt32? {
            return self.ownedNFTs[id]?.playID
        }

        pub fun getNumberInPlaySet(id: UInt64): UInt32? {
            return self.ownedNFTs[id]?.numberInPlaySet
        }

        pub fun getSetID(id: UInt64): UInt32? {
            return self.ownedNFTs[id]?.setID
        }

        pub fun getSetName(id: UInt64): String? {
            return self.ownedNFTs[id]?.setName
        }

        pub fun getSeries(id:UInt64): UInt32? {
            return self.ownedNFTs[id]?.series
        }

        pub fun getMetaData(id: UInt64): {String: String}? {
            if let PlayDataID = self.getPlayID(id: id) {
                return TopShot.plays[PlayDataID]?.metadata
            } else {
                return nil
            }
        }

        // If a transaction destroys the Collection object,
        // All the NFTs contained within are also destroyed
        // Kind of like when Damien Lillard destroys the hopes and
        // dreams of the entire city of Houston
        //
        destroy() {
            destroy self.ownedNFTs
        }
    }

    // -----------------------------------------------------------------------
    // TopShot contract-level function definitions
    // -----------------------------------------------------------------------

    // createEmptyCollection creates a new, empty Collection object so that
    // a user can store it in their account storage.
    // Once they have a Collection in their storage, they are able to receive
    // Moments in transactions
    //
    pub fun createEmptyCollection(): @Collection {
        return <-create Collection()
    }

    // -----------------------------------------------------------------------
    // TopShot initialization function
    // -----------------------------------------------------------------------
    //
    init() {
        // initialize the fields
        self.plays = {}
        self.playID = 0
        self.setID = 0
        self.totalSupply = 0

        // Create a new collection
        let oldCollection <- self.account.storage[Collection] <- create Collection()
        destroy oldCollection

        // Create a safe, public reference to the Collection 
        // and store it in public reference storage
        self.account.published[&MomentCollectionPublic] = &self.account.storage[Collection] as &MomentCollectionPublic

        // Create a new Admin resource and store it
        let oldAdmin <- self.account.storage[Admin] <- create Admin()
        destroy oldAdmin

        emit ContractInitialized()
    }
}
 
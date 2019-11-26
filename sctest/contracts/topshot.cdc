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

pub struct Mold {
    pub let id: Int  // the unique ID that the mold has

    pub let name: String   // the name of the moment

    pub let rarityCounts: {String: Int}  // the number of moments that can be minted from each rarity for this mold

    pub var numLeft: {String: Int}      // the number of moments that have been minted from each rarity mold
                                          // cannot be greater than the corresponding rarityCounts

    init(id: Int, name: String, rarityCounts: {String: Int}) {
        self.id = id
        self.name = name
        self.rarityCounts = rarityCounts
        self.numLeft = rarityCounts
        // and then subtract from it when a new one is minted
        // and not let any more be minted when it gets to zero
        // when someone wants to get how many have been minted, just  
        // subtract numMinted from rarityCounts
        // this way, the dictionary keys would always match for the two
    }

    // called when a moment is minted
    pub fun updateNumLeft(rarity: String) {
        var numLeft = self.numLeft[rarity] ?? panic("missing rarity count!")
        
        numLeft = numLeft - 1
        
        self.numLeft[rarity] = numLeft
    }
}

pub resource MoldCollection {

    // variable size dictionary of Mold conforming tokens
    // Mold is a struct type with an `Int` ID field
    pub var molds: {Int: Mold}

    // the ID that is used to cast molds
    pub var moldID: Int

    init() {
        self.molds = {}
        self.moldID = 1
    }

    // castMold casts a mold struct and stores it in the dictionary
    // for the molds
    // the mold ID must be unused
    pub fun castMold(name: String, rarityCounts: {String: Int}) {
        var newMold: Mold = Mold(id: self.moldID, name: name, rarityCounts: rarityCounts)

        self.molds[self.moldID] = newMold

        // increment the ID so that it isn't used again
        self.moldID = self.moldID + 1

    }

    // getMoldName gets the name of a mold that is stored in this collection
    pub fun getMoldName(moldID: Int): String {
        return self.molds[moldID]?.name ?? panic("missing mold name!")
    }

    // getNumMomentsLeftInRarity get the number of moments left of a certain rarity
    // for the specified mold ID
    pub fun getNumMomentsLeftInRarity(id: Int, rarity: String): Int {
        let mold = self.molds[id] ?? panic("missing mold!")

        let numLeft = mold.numLeft[rarity] ?? panic("missing numLeft!")

        return numLeft
    }

    // getNumMintedInRarity returns the number of moments that have been minted of 
    // a certain mold ID and rarity
    pub fun getNumMintedInRarity(id: Int, rarity: String): Int {
        let mold = self.molds[id] ?? panic("missing mold!")

        let numLeft = mold.numLeft[rarity] ?? panic("missing numLeft!")
        let rarityCount = mold.rarityCounts[rarity] ?? panic("missing rarity count!")

        return rarityCount - numLeft
    }

}


pub resource Moment {
    pub let id: Int

    pub var strength: Int

    pub var rarity: String

    // Tells which number of the Rarity this moment is
    pub let placeInRarity: Int

    // the ID of the mold that the moment references
    pub var moldID: Int

    // reference to the NBA Mold Collection that holds the mold
    pub let moldReference: &MoldCollection

    init(newID: Int, str: Int, moldID: Int, rarity: String, place: Int, reference: &MoldCollection) {
        self.id = newID
        self.strength = str
        self.moldID = moldID
        self.rarity = rarity
        self.placeInRarity = place
        self.moldReference = reference
    }

}

pub resource MomentCollection { 
    // dictionary of Moment conforming tokens
    // Moment is a resource type with an `Int` ID field
    pub var moments: <-{Int: Moment}

    init() {
        self.moments = {}
    }

    // withdraw removes an Moment from the collection and moves it to the caller
    pub fun withdraw(tokenID: Int): <-Moment {
        let token <- self.moments.remove(key: tokenID) ?? panic("missing Moment")
            
        return <-token
    }

    // deposit takes a Moment and adds it to the collections dictionary
    // and adds the ID to the id array
    pub fun deposit(token: <-Moment): Void {
        let id: Int = token.id
        
        var newToken: <-Moment? <- token

        // add the new token to the dictionary
        let oldToken <- self.moments[id] <- newToken

        destroy oldToken
    }

    // transfer takes a reference to another user's Moment collection,
    // takes the Moment out of this collection, and deposits it
    // in the reference's collection
    pub fun transfer(recipient: &MomentCollection, tokenID: Int): Void {

        // remove the token from the dictionary get the token from the optional
        let token <- self.withdraw(tokenID: tokenID)

        // deposit it in the recipient's account
        recipient.deposit(token: <-token)
    }

    // idExists checks to see if a Moment with the given ID exists in the collection
    pub fun idExists(tokenID: Int): Bool {
        if (self.moments[tokenID] != nil) {
            return true
        }

        return false
    }

    destroy() {
        destroy self.moments
    }

    // getIDs returns an array of the IDs that are in the collection
    pub fun getIDs(): [Int] {
        return self.moments.keys
    }
}


pub resource MomentFactory {

    // the ID that is used to mint moments
    pub var MomentID: Int

    // reference to this mold collection that can be used
    // to initialize the moments 
    pub var moldReference: &MoldCollection

    init(moldRef: &MoldCollection) {
        self.MomentID = 1
        self.moldReference = moldRef
    }

    // mintMoment mints a moment NFT based off of a mold that is stored in the collection
    // the moment ID must be unused
    pub fun mintMoment(moldID: Int, rarity: String, recipient: &MomentCollection) {

        // check to see if any more moments are allowed to be minted in this rarity
        let numLeft = self.moldReference.getNumMomentsLeftInRarity(id: moldID, rarity: rarity)
        if numLeft <= 0 { panic("All the moments of this rarity have been minted!") }

        // update the number left in the rarity that are allowed to be minted
        self.moldReference.molds[moldID]?.updateNumLeft(rarity: rarity)

        // gets this moment's place in the moments for this rarity
        let placeInRarity = self.moldReference.getNumMintedInRarity(id: moldID, rarity: rarity)

        // mint the new moment
        var newMoment: <-Moment <- create Moment(newID: self.MomentID, 
                                                str: 1, 
                                                moldID: moldID, 
                                                rarity: rarity, 
                                                place: placeInRarity, 
                                                reference: self.moldReference)
        
        // deposit the moment in the owner's account
        recipient.deposit(token: <-newMoment)

        // update the moment IDs so they can't be reused
        self.MomentID = self.MomentID + 1
    }

    pub fun createMomentCollection(): <-MomentCollection {
        return <-create MomentCollection()
    }
}

pub fun createMoldCollection(): <-MoldCollection {
    return <-create MoldCollection()
}

pub fun createMomentFactory(ref: &MoldCollection): <-MomentFactory {
    return <-create MomentFactory(moldRef: ref)
}
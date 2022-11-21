// TopShotTier
//
// TopShot NFTs weren't launched with a way to embed tier (rarity) in them.
// Because cadence doesn't permit adding new fields to structs/resources, 
// a new contract to maintain that mapping needs to be made, and historical data needs to be added.
// These mappings boil down to two possibilities:
// 1. a unique override for a specific Moment ID
// 2. A combination of setID + playID
pub contract TopShotTier {
    pub let AdminStoragePath: StoragePath
    pub let MappingStoragePath: StoragePath
    pub let MappingPublicPath: PublicPath

    pub let mappings: TierMapping

    pub resource interface TierMappingPublic {
        pub fun resolveRarity(id: UInt64, playID: UInt64, setID: UInt64): String
    }

    // TierMapping
    // Maintains mappings for resolve moment rarity.
    pub resource TierMapping: TierMappingPublic {
        // For most moments, the mapping of "setID+playID" is used to get tier. 
        // A moment can also have an override which will be present on that moment's id as a string.
        pub let mappings: {String: String}

        // resolveRarity
        // Returns the rarity of a TopShot Moment.
        // First we will check if there are any explicit overrides for a particular moment,
        // then we check for a unique grouping of setID + playID.
        pub fun resolveRarity(id: UInt64, playID: UInt64, setID: UInt64): String {
            if self.mappings[id.toString()] != nil {
                return self.mappings[id.toString()]!
            }

            let tier = self.mappings[setID.toString().concat("+").concat(playID.toString())]
            return tier != nil ? tier! : ""
        }

        pub fun addMapping(setID: UInt64, playID: UInt64, tier: String) {
            self.mappings[setID.toString().concat("+").concat(playID.toString())] = tier
        }

        pub fun addOverride(id: UInt64, tier: String) {
            self.mappings[id.toString()] = tier
        }

        init() {
            self.mappings = {}
        }
    }

    pub fun resolveRarity(id: UInt64, playID: UInt64, setID: UInt64): String {
        return self.account.getCapability<&{TierMappingPublic}>(TopShotTier.MappingPublicPath).borrow()!.resolveRarity(id: id, playID: playID, setID: setID)
    }

    init() {
        self.MappingStoragePath = /storage/TopShotTierMapping
        self.MappingPublicPath = /public/TopShotTierMapping

        self.mappings = TierMapping()
        self.account.save(<- create TierMetadataAdmin(), to: self.AdminStoragePath)
        self.account.link<&TierMapping{TierMappingPublic}>(self.MappingPublicPath, target: self.MappingStoragePath)
    }
}
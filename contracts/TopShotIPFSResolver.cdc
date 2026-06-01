access(all) contract TopShotIPFSResolver {

    access(all) event CIDSet(setID: UInt32, playID: UInt32, subeditionID: UInt32, mediaType: String, cid: String)
    access(all) event GatewayUpdated(gateway: String)

    access(all) let numBuckets: UInt64

    // Each shard is a resource stored at its own path so writes
    // borrow a reference and mutate in place — no dict copy.
    access(all) resource Shard {
        access(all) let entries: {String: {String: String}}

        init() {
            self.entries = {}
        }

        access(contract) fun setCID(key: String, mediaType: String, cid: String) {
            var inner = self.entries[key] ?? {}
            inner[mediaType] = cid
            self.entries[key] = inner
        }
    }

    // IPFS gateway base URL (e.g., "https://ipfs.dapperlabs.com/ipfs/")
    access(contract) var gateway: String

    access(all) view fun getGateway(): String { return self.gateway }

    access(all) let AdminStoragePath: StoragePath

    access(all) resource Admin {
        access(all) fun setCID(setID: UInt32, playID: UInt32, subeditionID: UInt32, mediaType: String, cid: String) {
            let key = TopShotIPFSResolver.buildKey(setID: setID, playID: playID, subeditionID: subeditionID)
            let bucket = TopShotIPFSResolver.getBucket(setID: setID, playID: playID, subeditionID: subeditionID)
            let shard = TopShotIPFSResolver.borrowShard(bucket)
            shard.setCID(key: key, mediaType: mediaType, cid: cid)
            emit CIDSet(setID: setID, playID: playID, subeditionID: subeditionID, mediaType: mediaType, cid: cid)
        }

        access(all) fun setGateway(gateway: String) {
            TopShotIPFSResolver.gateway = gateway
            emit GatewayUpdated(gateway: gateway)
        }

        access(all) fun clearAllCIDs() {
            var i: UInt64 = 0
            while i < TopShotIPFSResolver.numBuckets {
                let path = TopShotIPFSResolver.shardPath(i)
                if let existing <- TopShotIPFSResolver.account.storage.load<@Shard>(from: path) {
                    destroy existing
                }
                i = i + 1
            }
        }
    }

    access(all) view fun getCIDs(setID: UInt32, playID: UInt32, subeditionID: UInt32): {String: String}? {
        let key = self.buildKey(setID: setID, playID: playID, subeditionID: subeditionID)
        let bucket = self.getBucket(setID: setID, playID: playID, subeditionID: subeditionID)
        let path = self.shardPath(bucket)
        if let shard = self.account.storage.borrow<&Shard>(from: path) {
            if let entry = shard.entries[key] {
                // Return a copy, not a reference
                let result: {String: String} = {}
                for mediaType in entry.keys {
                    result[mediaType] = entry[mediaType]
                }
                return result
            }
        }
        return nil
    }

    access(contract) fun borrowShard(_ bucket: UInt64): &Shard {
        let path = self.shardPath(bucket)
        if let shard = self.account.storage.borrow<&Shard>(from: path) {
            return shard
        }
        self.account.storage.save(<- create Shard(), to: path)
        return self.account.storage.borrow<&Shard>(from: path)!
    }

    access(all) view fun shardPath(_ bucket: UInt64): StoragePath {
        return StoragePath(identifier: "TopShotIPFSShard_".concat(bucket.toString()))!
    }

    access(all) view fun getBucket(setID: UInt32, playID: UInt32, subeditionID: UInt32): UInt64 {
        return (UInt64(setID) + UInt64(playID) + UInt64(subeditionID)) % self.numBuckets
    }

    access(all) view fun buildKey(setID: UInt32, playID: UInt32, subeditionID: UInt32): String {
        return setID.toString()
            .concat("_")
            .concat(playID.toString())
            .concat("_")
            .concat(subeditionID.toString())
    }

    init() {
        self.numBuckets = 31
        self.gateway = "https://ipfs.dapperlabs.com/ipfs/"
        self.AdminStoragePath = /storage/TopShotIPFSResolverAdmin
        self.account.storage.save<@Admin>(<- create Admin(), to: self.AdminStoragePath)
    }
}

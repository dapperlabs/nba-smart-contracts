access(all) contract TopShotIPFSResolver {

    access(all) event CIDSet(setID: UInt32, playID: UInt32, subeditionID: UInt32, mediaType: String, cid: String)
    access(all) event GatewayUpdated(gateway: String)

    // Number of buckets to split CID entries across
    access(all) let numBuckets: UInt64

    // Sharded storage: bucket index -> composite key -> { mediaType: CID }
    // Let's us store 32 * 100k = 32M CID entries
    access(contract) let shards: {UInt64: {String: {String: String}}}

    // IPFS gateway base URL (e.g., "https://ipfs.dapperlabs.com/ipfs/")
    access(all) var gateway: String

    access(all) let AdminStoragePath: StoragePath

    access(all) resource Admin {
        access(all) fun setCID(setID: UInt32, playID: UInt32, subeditionID: UInt32, mediaType: String, cid: String) {
            let key = TopShotIPFSResolver.buildKey(setID: setID, playID: playID, subeditionID: subeditionID)
            let bucket = TopShotIPFSResolver.getBucket(setID: setID, playID: playID, subeditionID: subeditionID)
            var shard = TopShotIPFSResolver.shards[bucket] ?? {}
            var inner = shard[key] ?? {}
            inner[mediaType] = cid
            shard[key] = inner
            TopShotIPFSResolver.shards[bucket] = shard
            emit CIDSet(setID: setID, playID: playID, subeditionID: subeditionID, mediaType: mediaType, cid: cid)
        }

        access(all) fun setGateway(gateway: String) {
            TopShotIPFSResolver.gateway = gateway
            emit GatewayUpdated(gateway: gateway)
        }
    }

    access(all) view fun getCIDs(setID: UInt32, playID: UInt32, subeditionID: UInt32): {String: String}? {
        let key = self.buildKey(setID: setID, playID: playID, subeditionID: subeditionID)
        let bucket = self.getBucket(setID: setID, playID: playID, subeditionID: subeditionID)
        let shard = self.shards[bucket] ?? {}
        return shard[key]
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
        self.numBuckets = 10
        self.shards = {}
        self.gateway = "https://ipfs.dapperlabs.com/ipfs/"
        self.AdminStoragePath = /storage/TopShotIPFSResolverAdmin
        self.account.storage.save<@Admin>(<- create Admin(), to: self.AdminStoragePath)
    }
}

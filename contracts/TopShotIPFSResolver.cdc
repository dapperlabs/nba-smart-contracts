access(all) contract TopShotIPFSResolver {

    access(all) event CIDSet(setID: UInt32, playID: UInt32, subeditionID: UInt32, mediaType: String, cid: String)
    access(all) event GatewayUpdated(gateway: String)

    // Composite key "setID_playID_subeditionID" → { mediaType: CID }
    access(contract) let cids: {String: {String: String}}

    // IPFS gateway base URL (e.g., "https://ipfs.dapperlabs.com/ipfs/")
    access(all) var gateway: String

    access(all) let AdminStoragePath: StoragePath

    access(all) resource Admin {
        access(all) fun setCID(setID: UInt32, playID: UInt32, subeditionID: UInt32, mediaType: String, cid: String) {
            let key = TopShotIPFSResolver.buildKey(setID: setID, playID: playID, subeditionID: subeditionID)
            var inner = TopShotIPFSResolver.cids[key] ?? {}
            inner[mediaType] = cid
            TopShotIPFSResolver.cids[key] = inner
            emit CIDSet(setID: setID, playID: playID, subeditionID: subeditionID, mediaType: mediaType, cid: cid)
        }

        access(all) fun setGateway(gateway: String) {
            TopShotIPFSResolver.gateway = gateway
            emit GatewayUpdated(gateway: gateway)
        }
    }

    access(all) view fun getCIDs(setID: UInt32, playID: UInt32, subeditionID: UInt32): {String: String}? {
        let key = self.buildKey(setID: setID, playID: playID, subeditionID: subeditionID)
        return self.cids[key]
    }

    access(all) view fun buildKey(setID: UInt32, playID: UInt32, subeditionID: UInt32): String {
        return setID.toString()
            .concat("_")
            .concat(playID.toString())
            .concat("_")
            .concat(subeditionID.toString())
    }

    init() {
        self.cids = {}
        self.gateway = "https://ipfs.dapperlabs.com/ipfs/"
        self.AdminStoragePath = /storage/TopShotIPFSResolverAdmin
        self.account.storage.save<@Admin>(<- create Admin(), to: self.AdminStoragePath)
    }
}

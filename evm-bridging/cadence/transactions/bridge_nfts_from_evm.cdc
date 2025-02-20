import "FungibleToken"
import "NonFungibleToken"
import "ViewResolver"
import "MetadataViews"
import "FlowToken"
import "ScopedFTProviders"
import "EVM"
import "FlowEVMBridge"
import "FlowEVMBridgeConfig"
import "FlowEVMBridgeUtils"
import "CrossVMMetadataViews"

/// Bridges NFTs with provided IDs from EVM to Cadence, unwrapping them first if applicable.
///
/// @param nftIdentifier: The identifier of the NFT to unwrap and bridge (e.g., 'A.0b2a3299cc857e29.TopShot.NFT')
/// @param nftIDs: The IDs of the NFTs to bridge to Cadence from EVM (must be held in the signer's COA before submission)
///
transaction(
    nftIdentifier: String,
    nftIDs: [UInt256]
) {
    let nftType: Type
    let collection: &{NonFungibleToken.Collection}
    let scopedProvider: @ScopedFTProviders.ScopedFTProvider
    let coa: auth(EVM.Bridge, EVM.Call) &EVM.CadenceOwnedAccount
    let viewResolver: &{ViewResolver}

    prepare(signer: auth(BorrowValue, CopyValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {
        // Borrow a reference to the signer's COA
        self.coa = signer.storage.borrow<auth(EVM.Bridge, EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("Could not borrow COA signer's account at path /storage/evm")

        // Get NFT type, address, and name from the provided identifier
        self.nftType = CompositeType(nftIdentifier)
            ?? panic("Could not construct NFT type from identifier: ".concat(nftIdentifier))
        let nftContractAddress = FlowEVMBridgeUtils.getContractAddress(fromType: self.nftType)
            ?? panic("Could not get contract address from identifier: ".concat(nftIdentifier))
        let nftContractName = FlowEVMBridgeUtils.getContractName(fromType: self.nftType)
            ?? panic("Could not get contract name from identifier: ".concat(nftIdentifier))

        // Borrow a reference to the NFT collection, configuring if necessary
        self.viewResolver = getAccount(nftContractAddress).contracts.borrow<&{ViewResolver}>(name: nftContractName)
            ?? panic("Could not borrow ViewResolver from NFT contract with name "
                .concat(nftContractName).concat(" and address ")
                .concat(nftContractAddress.toString()))
        let collectionData = self.viewResolver.resolveContractView(
                resourceType: self.nftType,
                viewType: Type<MetadataViews.NFTCollectionData>()
            ) as! MetadataViews.NFTCollectionData?
            ?? panic("Could not resolve NFTCollectionData view for NFT type ".concat(self.nftType.identifier))
        if signer.storage.borrow<&{NonFungibleToken.Collection}>(from: collectionData.storagePath) == nil {
            signer.storage.save(<-collectionData.createEmptyCollection(), to: collectionData.storagePath)
            signer.capabilities.unpublish(collectionData.publicPath)
            let collectionCap = signer.capabilities.storage.issue<&{NonFungibleToken.Collection}>(collectionData.storagePath)
            signer.capabilities.publish(collectionCap, at: collectionData.publicPath)
        }
        self.collection = signer.storage.borrow<&{NonFungibleToken.Collection}>(from: collectionData.storagePath)
            ?? panic("Could not borrow a NonFungibleToken Collection from the signer's storage path "
                    .concat(collectionData.storagePath.toString()))

        // Set a cap on the withdrawable bridge fee
        var approxFee = FlowEVMBridgeUtils.calculateBridgeFee(
                bytes: 400_000 // 400 kB as upper bound on movable storage used in a single transaction
            ) + (FlowEVMBridgeConfig.baseFee * UFix64(nftIDs.length))
        // Issue and store bridge-dedicated Provider Capability in storage if necessary
        if signer.storage.type(at: FlowEVMBridgeConfig.providerCapabilityStoragePath) == nil {
            let providerCap = signer.capabilities.storage.issue<auth(FungibleToken.Withdraw) &{FungibleToken.Provider}>(
                /storage/flowTokenVault
            )
            signer.storage.save(providerCap, to: FlowEVMBridgeConfig.providerCapabilityStoragePath)
        }

        // Copy the stored Provider capability and create a ScopedFTProvider
        let providerCapCopy = signer.storage.copy<Capability<auth(FungibleToken.Withdraw) &{FungibleToken.Provider}>>(
                from: FlowEVMBridgeConfig.providerCapabilityStoragePath
            ) ?? panic("Invalid FungibleToken Provider Capability found in storage at path "
                .concat(FlowEVMBridgeConfig.providerCapabilityStoragePath.toString()))
        let providerFilter = ScopedFTProviders.AllowanceFilter(approxFee)
        self.scopedProvider <- ScopedFTProviders.createScopedFTProvider(
                provider: providerCapCopy,
                filters: [ providerFilter ],
                expiration: getCurrentBlock().timestamp + 1.0
            )
    }

    execute {
        // Unwrap NFTs if applicable
        unwrapNFTsIfApplicable(self.coa,
            nftIDs: nftIDs,
            nftType: self.nftType,
            viewResolver: self.viewResolver
        )

        // Bridge each NFT from the signer's COA in EVM to the signer's collection in Cadence
        for id in nftIDs {
            // Execute the bridge to Cadence for the current ID
            let nft: @{NonFungibleToken.NFT}  <- self.coa.withdrawNFT(
                type: self.nftType,
                id: id,
                feeProvider: &self.scopedProvider as auth(FungibleToken.Withdraw) &{FungibleToken.Provider}
            )

            // Ensure the bridged nft is the correct type
            assert(
                nft.getType() == self.nftType,
                message: "Bridged nft type mismatch - requested: ".concat(self.nftType.identifier)
                    .concat(", received: ").concat(nft.getType().identifier)
            )

            // Deposit the bridged NFT into the signer's collection
            self.collection.deposit(token: <-nft)
        }

        // Destroy the ScopedFTProvider
        destroy self.scopedProvider
    }
}

/// Unwraps NFTs from a project's custom ERC721 wrapper contract to bridged NFTs on EVM, if applicable.
/// Enables projects to use their own ERC721 contract while leveraging the bridge's underlying contract,
/// until direct custom contract support is added to the bridge.
///
/// @param coa: The COA of the signer
/// @param nftIDs: The IDs of the NFTs to wrap
/// @param nftType: The type of the NFTs to wrap
/// @param viewResolver: The ViewResolver of the NFT contract
///
access(all) fun unwrapNFTsIfApplicable(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    nftIDs: [UInt256],
    nftType: Type,
    viewResolver: &{ViewResolver}
) {
    // Get the project-defined ERC721 address if it exists
    if let crossVMPointer = viewResolver.resolveContractView(
            resourceType: nftType,
            viewType: Type<CrossVMMetadataViews.EVMPointer>()
    ) as! CrossVMMetadataViews.EVMPointer? {
        // Get the underlying ERC721 address if it exists
        if let underlyingAddress = getUnderlyingERC721Address(coa, crossVMPointer.evmContractAddress) {
            for id in nftIDs {
                // Unwrap NFT if it is wrapped
                if isNFTWrapped(coa,
                    nftID: id,
                    underlying: underlyingAddress,
                    wrapper: crossVMPointer.evmContractAddress
                ) {
                    let res = mustCall(coa, crossVMPointer.evmContractAddress,
                        functionSig: "withdrawTo(address,uint256[])",
                        args: [coa.address(), [id]]
                    )
                    let decodedRes = EVM.decodeABI(types: [Type<Bool>()], data: res.data)
                    assert(decodedRes.length == 1, message: "Invalid response length")
                    assert(decodedRes[0] as! Bool, message: "Failed to unwrap NFT")
                }
            }
        }
    }
}

/// Calls a function on an EVM contract from provided coa
///
access(all) fun mustCall(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ contractAddr: EVM.EVMAddress,
    functionSig: String,
    args: [AnyStruct]
): EVM.Result {
    let res = coa.call(
        to: contractAddr,
        data: EVM.encodeABIWithSignature(functionSig, args),
        gasLimit: 4_000_000,
        value: EVM.Balance(attoflow: 0)
    )

    assert(res.status == EVM.Status.successful,
        message: "Failed to call '".concat(functionSig)
            .concat("\n\t error code: ").concat(res.errorCode.toString())
            .concat("\n\t error message: ").concat(res.errorMessage)
            .concat("\n\t gas used: ").concat(res.gasUsed.toString())
            .concat("\n\t args count: ").concat(args.length.toString())
            .concat("\n\t caller address: 0x").concat(coa.address().toString())
            .concat("\n\t contract address: 0x").concat(contractAddr.toString())
    )

    return res
}

/// Gets the underlying ERC721 address if it exists (i.e. if the ERC721 is a wrapper)
///
access(all) fun getUnderlyingERC721Address(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ wrapperAddress: EVM.EVMAddress
): EVM.EVMAddress? {
    let res = coa.call(
        to: wrapperAddress,
        data: EVM.encodeABIWithSignature("underlying()", []),
        gasLimit: 100_000,
        value: EVM.Balance(attoflow: 0)
    )

    // If the call fails, return nil
    if res.status != EVM.Status.successful || res.data.length == 0 {
        return nil
    }

    // Decode and return the underlying ERC721 address
    let decodedResult = EVM.decodeABI(
        types: [Type<EVM.EVMAddress>()],
        data: res.data
    )
    assert(decodedResult.length == 1, message: "Invalid response length")
    return decodedResult[0] as! EVM.EVMAddress
}

/// Checks if the provided NFT is wrapped in the underlying ERC721 contract
///
access(all) fun isNFTWrapped(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    nftID: UInt256,
    underlying: EVM.EVMAddress,
    wrapper: EVM.EVMAddress
): Bool {
    let res = coa.call(
        to: underlying,
        data: EVM.encodeABIWithSignature("ownerOf(uint256)", [nftID]),
        gasLimit: 100_000,
        value: EVM.Balance(attoflow: 0)
    )

    // If the call fails, return false
    if res.status != EVM.Status.successful || res.data.length == 0{
        return false
    }

    // Decode and compare the addresses
    let decodedResult = EVM.decodeABI(
        types: [Type<EVM.EVMAddress>()],
        data: res.data
    )
    assert(decodedResult.length == 1, message: "Invalid response length")
    let owner = decodedResult[0] as! EVM.EVMAddress
    return owner.toString() == wrapper.toString()
}

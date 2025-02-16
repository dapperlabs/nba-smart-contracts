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

/// Bridges NFTs with provided IDs from Cadence to EVM, wrapping them in a wrapper ERC721 if applicable.
///
/// @param nftIdentifier: The identifier of the NFT to wrap and bridge (e.g., 'A.877931736ee77cff.TopShot.NFT')
/// @param nftIDs: Array of IDs of the NFTs to wrap
/// @param recipientEvmAddressIfnotCoa: The EVM address
///
transaction(
    nftIdentifier: String,
    nftIDs: [UInt64],
    recipientEvmAddressIfnotCoa: String?
) {
    let nftType: Type
    let collection: auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}
    let coa: auth(EVM.Bridge, EVM.Call) &EVM.CadenceOwnedAccount
    let scopedProvider: @ScopedFTProviders.ScopedFTProvider
    let viewResolver: &{ViewResolver}

    prepare(signer: auth(CopyValue, BorrowValue, IssueStorageCapabilityController, PublishCapability, SaveValue) &Account) {
        // Retrieve or create COA in signer's account
        if let coa = signer.storage.borrow<auth(EVM.Call, EVM.Bridge) &EVM.CadenceOwnedAccount>(from: /storage/evm) {
            self.coa = coa
        } else {
            signer.storage.save<@EVM.CadenceOwnedAccount>(<- EVM.createCadenceOwnedAccount(), to: /storage/evm)
            signer.capabilities.publish(
                signer.capabilities.storage.issue<&EVM.CadenceOwnedAccount>(/storage/evm),
                at: /public/evm
            )
            self.coa = signer.storage.borrow<auth(EVM.Call, EVM.Bridge) &EVM.CadenceOwnedAccount>(from: /storage/evm)!
        }

        // Get NFT collection info
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
        self.collection = signer.storage.borrow<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
                from: collectionData.storagePath
            ) ?? panic("Could not borrow a NonFungibleToken Collection from the signer's storage path "
                .concat(collectionData.storagePath.toString()))

        // Withdraw the requested NFT & set a cap on the withdrawable bridge fee
        var approxFee = FlowEVMBridgeUtils.calculateBridgeFee(
                bytes: 400_000 // 400 kB as upper bound on movable storage used in a single transaction
            ) + (FlowEVMBridgeConfig.baseFee * UFix64(nftIDs.length))

        // Check that the NFT is onboarded
        let requiresOnboarding = FlowEVMBridge.typeRequiresOnboarding(self.nftType)
            ?? panic("Bridge does not support the requested asset type ".concat(nftIdentifier))
        assert(!requiresOnboarding, message: "NFT must be onboarded before bridging and wrapping")

        // Issue and store bridge-dedicated Provider Capability in storage if necessary
        if signer.storage.type(at: FlowEVMBridgeConfig.providerCapabilityStoragePath) == nil {
            let providerCap = signer.capabilities.storage.issue<auth(FungibleToken.Withdraw) &{FungibleToken.Provider}>(
                /storage/flowTokenVault
            )
            signer.storage.save(providerCap, to: FlowEVMBridgeConfig.providerCapabilityStoragePath)
        }

        // Copy the stored Provider capability and create a ScopedFTProvider
        let providerCapCopy = signer.storage.copy<Capability<auth(FungibleToken.Withdraw) &{FungibleToken.Provider}>>(
            from: FlowEVMBridgeConfig.providerCapabilityStoragePath)
                ?? panic("Invalid FungibleToken Provider Capability found in storage at path "
                    .concat(FlowEVMBridgeConfig.providerCapabilityStoragePath.toString()))
        let providerFilter = ScopedFTProviders.AllowanceFilter(approxFee)
        self.scopedProvider <- ScopedFTProviders.createScopedFTProvider(
            provider: providerCapCopy,
            filters: [ providerFilter ],
            expiration: getCurrentBlock().timestamp + 1.0
        )
    }

    execute {
        // Iterate over requested IDs and bridge each NFT to the signer's COA in EVM
        for id in nftIDs {
            // Withdraw the NFT & ensure it's the correct type
            let nft <- self.collection.withdraw(withdrawID: id)
            assert(
                nft.getType() == self.nftType,
                message: "Bridged nft type mismatch - requested: ".concat(self.nftType.identifier)
                    .concat(", received: ").concat(nft.getType().identifier)
            )
            // Execute the bridge to EVM for the current ID
            self.coa.depositNFT(
                nft: <- nft,
                feeProvider: &self.scopedProvider as auth(FungibleToken.Withdraw) &{FungibleToken.Provider}
            )
        }

        // Destroy the ScopedFTProvider
        destroy self.scopedProvider

        // Wrap NFTs if applicable
        wrapAndTransferNFTsIfApplicable(self.coa,
            nftIDs: nftIDs,
            nftType: self.nftType,
            viewResolver: self.viewResolver,
            recipientIfNotCoa: recipientEvmAddressIfnotCoa != nil: EVM.addressFromString(recipientEvmAddressIfnotCoa!) ? nil
        )
    }
}

/// Wraps and transfers bridged NFTs into a project's custom ERC721 wrapper contract on EVM, if applicable.
/// Enables projects to use their own ERC721 contract while leveraging the bridge's underlying contract,
/// until direct custom contract support is added to the bridge.
///
/// @param coa: The COA of the signer
/// @param nftIDs: The IDs of the NFTs to wrap
/// @param nftType: The type of the NFTs to wrap
/// @param viewResolver: The ViewResolver of the NFT contract
/// @param recipientIfNotCoa: The EVM address to transfer the wrapped NFTs to, nil if the NFTs should stay in signer's COA
///
access(all) fun wrapAndTransferNFTsIfApplicable(
    _ coa: auth(EVM.Call, EVM.Bridge) &EVM.CadenceOwnedAccount,
    nftIDs: [UInt64],
    nftType: Type,
    viewResolver: &{ViewResolver},
    recipientIfNotCoa: EVM.EVMAddress?
) {
    // Get the project-defined ERC721 address if it exists
    if let crossVMPointer = viewResolver.resolveContractView(
            resourceType: nftType,
            viewType: Type<CrossVMMetadataViews.EVMPointer>()
    ) as! CrossVMMetadataViews.EVMPointer? {
        // Get the underlying ERC721 address if it exists
        if let underlyingAddress = getUnderlyingERC721Address(coa, crossVMPointer.evmContractAddress) {
            // Wrap NFTs if underlying ERC721 address matches bridge's associated address for NFT type
            if underlyingAddress.equals(FlowEVMBridgeConfig.getEVMAddressAssociated(with: nftType)!) {
                // Approve contract to withdraw underlying NFTs from signer's coa
                mustCall(coa, underlyingAddress,
                    functionSig: "setApprovalForAll(address,bool)",
                    args: [crossVMPointer.evmContractAddress, true]
                )

                // Wrap NFTs with provided IDs, and check if the call was successful
                let res = mustCall(coa, crossVMPointer.evmContractAddress,
                    functionSig: "depositFor(address,uint256[])",
                    args: [coa.address(), nftIDs]
                )
                let decodedRes = EVM.decodeABI(types: [Type<Bool>()], data: res.data)
                assert(decodedRes.length == 1, message: "Invalid response length")
                assert(decodedRes[0] as! Bool, message: "Failed to wrap NFTs")

                // Transfer NFTs to recipient if provided
                if let to = recipientIfNotCoa {
                    mustTransferNFTs(coa, crossVMPointer.evmContractAddress, nftIDs: nftIDs, to: to)
                }

                // Revoke approval for contract to withdraw underlying NFTs from signer's coa
                mustCall(coa, underlyingAddress,
                    functionSig: "setApprovalForAll(address,bool)",
                    args: [crossVMPointer.evmContractAddress, false]
                )
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
): EVM.EVMResult {
    let res = coa.call(
        to: contractAddr,
        data: EVM.encodeABIWithSignature(functionSig, args),
        gasLimit: 400_000,
        value: EVM.Balance(attoflow: 0)
    )

    assert(res.status == EVM.Status.successful,
        message: "Failed to call '".concat(functionSig).concat("'\n\t\t error code: ")
            .concat(res.errorCode.toString()).concat("\n\t\t message: ")
            .concat(res.errorMessage)
    )

    return res
}

/// Transfers NFTs from the provided COA to the provided EVM address
///
access(all) fun mustTransferNFTs(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ erc721Address: EVM.EVMAddress,
    nftIDs: [UInt64],
    to: EVM.EVMAddress
) {
    for id in nftIDs {
        assert(isOwner(coa, erc721Address, id, coa.address()), message: "NFT not owned by signer's COA")
        mustCall(coa, erc721Address,
            functionSig: "safeTransferFrom(address,address,uint256)",
            args: [coa.address(), to, id]
        )
        assert(isOwner(coa, erc721Address, id, to), message: "NFT not transferred to recipient")
    }
}

/// Checks if the provided NFT is owned by the provided EVM address
///
access(all) fun isOwner(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ erc721Address: EVM.EVMAddress,
    _ nftID: UInt64,
    _ ownerToCheck: EVM.EVMAddress
): Bool {
    let res = coa.call(
        to: erc721Address,
        data: EVM.encodeABIWithSignature("ownerOf(uint256)", [nftID]),
        gasLimit: 100_000,
        value: EVM.Balance(attoflow: 0)
    )
    assert(res.status == EVM.Status.successful, message: "Call to ERC721.ownerOf(uint256) failed")
    let decodedRes = EVM.decodeABI(types: [Type<EVM.EVMAddress>()], data: res.data)
    if decodedRes.length == 1 {
        let actualOwner = decodedRes[0] as! EVM.EVMAddress
        return actualOwner.equals(ownerToCheck)
    }
    return false
}

/// Gets the underlying ERC721 address
///
access(all) fun getUnderlyingERC721Address(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ wrapperAddress: EVM.EVMAddress
): EVM.EVMAddress {
    let res = coa.call(
        to: wrapperAddress,
        data: EVM.encodeABIWithSignature("underlying()", []),
        gasLimit: 100_000,
        value: EVM.Balance(attoflow: 0)
    )

    assert(res.status == EVM.Status.successful,
        message: "Failed to call 'underlying()'\n\t\t error code: "
            .concat(res.errorCode.toString()).concat("\n\t\t message: ")
            .concat(res.errorMessage)
    )
    let decodedResult = EVM.decodeABI(types: [Type<EVM.EVMAddress>()], data: res.data)
    assert(decodedResult.length == 1, message: "Invalid response length")

    return decodedResult[0] as! EVM.EVMAddress
}

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

/// Bridges wrapped NFTs with provided IDs from EVM to Cadence
///
/// @param wrapperERC721Address: The EVM address of the wrapper ERC721 NFT
/// @param nftIDs: The ERC721 ids of the NFTs to bridge to Cadence from EVM
///
transaction(wrapperERC721Address: String, nftIDs: [UInt256]) {

    let nftType: Type
    let collection: &{NonFungibleToken.Collection}
    let scopedProvider: @ScopedFTProviders.ScopedFTProvider
    let coa: auth(EVM.Bridge, EVM.Call) &EVM.CadenceOwnedAccount

    prepare(signer: auth(BorrowValue, CopyValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {
        /* --- Reference the signer's CadenceOwnedAccount --- */
        //
        // Borrow a reference to the signer's COA
        self.coa = signer.storage.borrow<auth(EVM.Bridge, EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("Could not borrow COA signer's account at path /storage/evm")

        // Get NFT identifier from EVM contract
        let nftIdentifier = getNFTIdentifier(self.coa, EVM.addressFromString(wrapperERC721Address))

        /* --- Construct the NFT type --- */
        //
        // Construct the NFT type from the provided identifier
        self.nftType = CompositeType(nftIdentifier)
            ?? panic("Could not construct NFT type from identifier: ".concat(nftIdentifier))
        // Parse the NFT identifier into its components
        let nftContractAddress = FlowEVMBridgeUtils.getContractAddress(fromType: self.nftType)
            ?? panic("Could not get contract address from identifier: ".concat(nftIdentifier))
        let nftContractName = FlowEVMBridgeUtils.getContractName(fromType: self.nftType)
            ?? panic("Could not get contract name from identifier: ".concat(nftIdentifier))

        /* --- Reference the signer's NFT Collection --- */
        //
        // Borrow a reference to the NFT collection, configuring if necessary
        let viewResolver = getAccount(nftContractAddress).contracts.borrow<&{ViewResolver}>(name: nftContractName)
            ?? panic("Could not borrow ViewResolver from NFT contract with name "
                .concat(nftContractName).concat(" and address ")
                .concat(nftContractAddress.toString()))
        let collectionData = viewResolver.resolveContractView(
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

        /* --- Configure a ScopedFTProvider --- */
        //
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
        // Get contract addresses
        let wrapperAddress = EVM.addressFromString(wrapperERC721Address)
        let underlyingAddress = getUnderlyingERC721Address(self.coa, wrapperAddress)

        // Approve contract to withdraw underlying NFTs from signer's coa
        call(self.coa, wrapperAddress,
            functionSig: "setApprovalForAll(address,bool)",
            args: [underlyingAddress, true]
        )

        // Unwrap NFTs with provided IDs
        call(self.coa, wrapperAddress,
            functionSig: "withdrawTo(address,uint256[])",
            args: [self.coa.address(), nftIDs]
        )

        // Iterate over the provided nftIDs
        for id in nftIDs {
            // Execute the bridge
            let nft: @{NonFungibleToken.NFT} <- self.coa.withdrawNFT(
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


/// Gets the NFT identifier from the EVM contract
///
access(all) fun getNFTIdentifier(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ wrapperAddress: EVM.EVMAddress
): String {
    let res = coa.call(
        to: wrapperAddress,
        data: EVM.encodeABIWithSignature("getCadenceIdentifier()", []),
        gasLimit: 100_000,
        value: EVM.Balance(attoflow: 0)
    )

    assert(res.status == EVM.Status.successful, message: "Call to get underlying ERC721 address failed")
    let decodedResult = EVM.decodeABI(types: [Type<String>()], data: res.data)
    assert(decodedResult.length == 1, message: "Invalid response length")

    return decodedResult[0] as! String
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

    assert(res.status == EVM.Status.successful, message: "Call to get underlying ERC721 address failed")
    let decodedResult = EVM.decodeABI(types: [Type<EVM.EVMAddress>()], data: res.data)
    assert(decodedResult.length == 1, message: "Invalid response length")

    return decodedResult[0] as! EVM.EVMAddress
}

/// Calls a function on an EVM contract from provided coa
///
access(all) fun call(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ contractAddr: EVM.EVMAddress,
    functionSig: String,
    args: [AnyStruct],
) {
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
}

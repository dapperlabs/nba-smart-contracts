import "EVM"

/// Unwraps NFTs with provided IDs
///
/// @param wrapperERC721Address: EVM address of the wrapper ERC721 NFT
/// @param nftIDs: Array of IDs of the NFTs to unwrap
///
transaction(
    wrapperERC721Address: String,
    nftIDs: [UInt256]
) {
    // Cadence-owned account
    let coa: auth(EVM.Call) &EVM.CadenceOwnedAccount

    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow a reference to the signer's COA
        self.coa = signer.storage.borrow<auth(EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("No COA found in signer's account")
    }

    execute {
        // Get contract addresses
        let wrapperAddress = EVM.addressFromString(wrapperERC721Address)
        let underlyingAddress = getUnderlyingERC721Address(self.coa, wrapperAddress)

        // Approve contract to withdraw underlying NFTs from signer's coa
        mustCall(self.coa, wrapperAddress,
            functionSig: "setApprovalForAll(address,bool)",
            args: [underlyingAddress, true]
        )

        // Unwrap NFTs with provided IDs
        mustCall(self.coa, wrapperAddress,
            functionSig: "withdrawTo(address,uint256[])",
            args: [self.coa.address(), nftIDs]
        )
    }
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
access(all) fun mustCall(
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

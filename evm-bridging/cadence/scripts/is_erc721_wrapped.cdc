import "EVM"

/// Returns true if the NFT is wrapped in the underlying ERC721 contract
///
/// @param flowAccountWithCoa - A Flow account with a COA
/// @param nftID - The ID of the NFT
/// @param underlying - The address of the underlying ERC721 contract
/// @param wrapper - The address of the wrapper ERC721 contract
///
access(all) fun main(flowAccountWithCoa: Address, nftID: UInt64, underlying: String, wrapper: String): Bool {
    let coa = getAuthAccount<auth(BorrowValue) &Account>(flowAccountWithCoa)
        .storage.borrow<auth(EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("No COA found in signer's account")

    return isNFTWrapped(coa, nftID: nftID, underlying: EVM.addressFromString(underlying), wrapper: EVM.addressFromString(wrapper))
}

/// Checks if the provided NFT is wrapped in the underlying ERC721 contract
///
access(all) fun isNFTWrapped(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    nftID: UInt64,
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
    if res.status != EVM.Status.successful {
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

import "EVM"

/// Returns the hex encoded address of the underlying ERC721 contract
///
access(all) fun main(flowNftAddress: Address, wrapperERC721Address: String): String? {
    let coa = getAuthAccount<auth(BorrowValue) &Account>(flowNftAddress)
        .storage.borrow<auth(EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("No COA found in signer's account")

    return getUnderlyingERC721Address(coa,
        EVM.addressFromString(wrapperERC721Address)
    ).toString()
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

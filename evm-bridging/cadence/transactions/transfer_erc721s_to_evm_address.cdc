import "Burner"
import "FungibleToken"
import "FlowToken"
import "EVM"

/// Transfers ERC721s from the signer's COA to an EVM address. All tokens must be defined in the same
/// contract and will be sent to the defined recipient.
///
/// @param erc721Address - The EVM address of the ERC721 contract
/// @param toEVMAddress - The EVM address to transfer the NFT to
/// @param nftIDs - The IDs of the NFTs to transfer
///
transaction(
    erc721Address: String,
    toEVMAddress: String,
    nftIDs: [UInt256]
) {
    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow COA from signer's account storage
        let coa = signer.storage.borrow<auth(EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("Could not find coa in signer's account.")

        // Transfer NFTs from signer's COA to provided EVM address
        mustTransferNFTs(coa, EVM.addressFromString(erc721Address),
            nftIDs: nftIDs,
            to: EVM.addressFromString(toEVMAddress),
        )
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

/// Transfers NFTs from the provided COA to the provided EVM address
///
access(all) fun mustTransferNFTs(
    _ coa: auth(EVM.Call) &EVM.CadenceOwnedAccount,
    _ erc721Address: EVM.EVMAddress,
    nftIDs: [UInt256],
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
    _ nftID: UInt256,
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

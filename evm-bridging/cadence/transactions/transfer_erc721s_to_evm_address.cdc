import "Burner"
import "FungibleToken"
import "FlowToken"
import "EVM"

/// Transfers an NFT from the signer's COA to an EVM address
///
/// @param erc721Address - The address of the ERC721 contract
/// @param toEVMAddress - The address to transfer the NFT to
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

        // Parse addresses
        let erc721 = EVM.addressFromString(erc721Address)
        let to = EVM.addressFromString(toEVMAddress)

        // Transfer NFTs from signer's COA to provided EVM address
        for nftID in nftIDs {
            mustCall(coa, erc721,
                functionSig: "safeTransferFrom(address,address,uint256)",
                args: [coa.address(), to, nftID]
            )
        }
    }
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

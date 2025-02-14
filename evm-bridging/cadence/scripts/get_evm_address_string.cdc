import "EVM"

/// Returns the hex encoded address of the COA in the given Flow address
///
access(all) fun main(flowAddress: Address): String? {
    return getAuthAccount<auth(BorrowValue) &Account>(flowAddress)
        .storage.borrow<&EVM.CadenceOwnedAccount>(from: /storage/evm)
        ?.address()
        ?.toString()
        ?? nil
}

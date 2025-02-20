import "EVM"

transaction(evmContractAddress: String) {
    let coa: auth(EVM.Call) &EVM.CadenceOwnedAccount

    prepare(signer: auth(SaveValue, BorrowValue, Capabilities) &Account) {
        // Retrieve or create COA
        if let coa = signer.storage.borrow<auth(EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm) {
            self.coa = coa
        } else {
            signer.storage.save<@EVM.CadenceOwnedAccount>(<- EVM.createCadenceOwnedAccount(), to: /storage/evm)
            signer.capabilities.publish(
                signer.capabilities.storage.issue<&EVM.CadenceOwnedAccount>(/storage/evm),
                at: /public/evm
            )
            self.coa = signer.storage.borrow<auth(EVM.Call) &EVM.CadenceOwnedAccount>(from: /storage/evm)!
        }
    }

    execute {
        // Test array with various values including edge cases (min and max uint64 values)
        let testArray: [UInt64] = [0, 1, 999999, 18446744073709551615]

        // Encode and call
        let res = self.coa.call(
            to: EVM.addressFromString(evmContractAddress),
            data: EVM.encodeABIWithSignature(
                "testArrayEncoding(uint256[])",
                [testArray]
            ),
            gasLimit: 100_000,
            value: EVM.Balance(attoflow: 0)
        )

        assert(res.status == EVM.Status.successful, message: "Call failed")

        // Decode and verify
        let decoded = EVM.decodeABI(types: [Type<[UInt64]>()], data: res.data)
        let returnedArray = decoded[0] as! [UInt64]

        // Compare arrays
        assert(testArray.length == returnedArray.length, message: "Array length mismatch")
        for i, value in testArray {
            assert(value == returnedArray[i],
                message: "Mismatch at index ".concat(i.toString())
                    .concat(": expected ").concat(value.toString())
                    .concat(", got ").concat(returnedArray[i].toString())
            )
        }

        log("Array encoding test passed!")
    }
}
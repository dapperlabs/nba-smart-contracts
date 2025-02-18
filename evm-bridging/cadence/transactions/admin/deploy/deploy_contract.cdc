import "FungibleToken"
import "FlowToken"
import "EVM"

/// Deploys a compiled solidity contract from bytecode to the EVM, with the signer's COA as the deployer
///
transaction(bytecode: String, constructorArgs: [String], gasLimit: UInt64) {
    let coa: auth(EVM.Deploy) &EVM.CadenceOwnedAccount

    prepare(signer: auth(BorrowValue) &Account) {
        self.coa = signer.storage.borrow<auth(EVM.Deploy) &EVM.CadenceOwnedAccount>(from: /storage/evm)
            ?? panic("Could not borrow reference to the signer's bridged account")
    }

    execute {

        let code = bytecode.concat(String.encodeHex(constructorArgs)).decodeHex()

        let result = self.coa.deploy(
            code: code,
            gasLimit: gasLimit,
            value: EVM.Balance(attoflow: 0)
        )
        assert(result.status == EVM.Status.successful && result.deployedContract != nil,
            message: "EVM deployment failed with error code: ".concat(result.errorCode.toString())
        )
    }
}

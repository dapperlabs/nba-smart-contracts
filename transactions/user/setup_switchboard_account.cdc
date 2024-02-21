import FungibleTokenSwitchboard from 0xFUNGIBLETOKENSWITCHBOARDADDRESS
import FungibleToken from 0xFUNGIBLETOKENADDRESS

// This transaction is a template for a transaction that could be used by
// anyone to to add a Switchboard resource to their account so that they can
// receive multiple fungible tokens using a single {FungibleToken.Receiver}
transaction {

    prepare(signer: auth(BorrowValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {

        // Check if the account already has a Switchboard resource, return early if so
        if signer.storage.borrow<&FungibleTokenSwitchboard.Switchboard>(from: FungibleTokenSwitchboard.StoragePath) != nil {
            return
        }

        // Create a new Switchboard resource and put it into storage
        signer.storage.save(
            <- FungibleTokenSwitchboard.createSwitchboard(),
            to: FungibleTokenSwitchboard.StoragePath
        )

        // Clear existing Capabilities at canonical paths
        signer.capabilities.unpublish(FungibleTokenSwitchboard.ReceiverPublicPath)
        signer.capabilities.unpublish(FungibleTokenSwitchboard.PublicPath)

        // Create a public capability to the Switchboard exposing the deposit
        // function through the {FungibleToken.Receiver} interface
        let receiverCap = signer.capabilities.storage.issue<&{FungibleToken.Receiver}>(
                FungibleTokenSwitchboard.StoragePath
            )
        signer.capabilities.publish(receiverCap, at: FungibleTokenSwitchboard.ReceiverPublicPath)

        // Create a public capability to the Switchboard exposing both the
        // {FungibleTokenSwitchboard.SwitchboardPublic} and the
        // {FungibleToken.Receiver} interfaces
        let switchboardPublicCap = signer.capabilities.storage.issue<&{FungibleTokenSwitchboard.SwitchboardPublic, FungibleToken.Receiver}>(
                FungibleTokenSwitchboard.StoragePath
            )
        signer.capabilities.publish(switchboardPublicCap, at: FungibleTokenSwitchboard.PublicPath)

    }

}
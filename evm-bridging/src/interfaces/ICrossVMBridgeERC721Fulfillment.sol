pragma solidity 0.8.24;

import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title ICrossVMBridgeERC721Fulfillment
 * @dev Related to https://github.com/onflow/flips/issues/318[FLIP-318] Cross VM NFT implementations
 * on Flow in the context of Cadence-native NFTs. The following interface must be implemented to
 * integrate with the Flow VM bridge connecting Cadence & EVM implementations so that the canonical
 * VM bridge may move the Cadence NFT into EVM in a mint/escrow pattern.
 */
interface ICrossVMBridgeERC721Fulfillment is IERC165 {

    // Encountered when attempting to fulfill a token that has been previously minted and is not
    // escrowed in EVM under the VM bridge
    error FulfillmentFailedTokenNotEscrowed(uint256 id, address escrowAddress);

    // Emitted when an NFT is moved from Cadence into EVM
    event FulfilledToEVM(address indexed recipient, uint256 indexed tokenId);

    /**
     * @dev Returns whether the token is currently escrowed under custody of the designated VM bridge
     * 
     * @param _id the ID of the token in question
     */
    function isEscrowed(uint256 _id) external view returns (bool);

    function exists(uint256 _id) external view returns (bool);

    /**
     * @dev Fulfills the bridge request, minting (if non-existent) or transferring (if escrowed) the
     * token with the given ID to the provided address. For dynamic metadata handling between
     * Cadence & EVM, implementations should override and assign metadata as encoded from Cadence
     * side. If overriding, be sure to preserve the mint/escrow pattern as shown in the default
     * implementation. See `_beforeFulfillment` and `_afterFulfillment` hooks to enable pre-and/or
     * post-processing without the need to override this function.
     * 
     * @param _to address of the token recipient
     * @param _id the id of the token being moved into EVM from Cadence
     * @param _data any encoded metadata passed by the corresponding Cadence NFT at the time of
     *      bridging into EVM
     */
    function fulfillToEVM(address _to, uint256 _id, bytes memory _data) external;
}

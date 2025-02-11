// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.24;

import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {ICrossVMBridgeERC721Fulfillment} from "../interfaces/ICrossVMBridgeERC721Fulfillment.sol";
import {ICrossVMBridgeCallable} from "../interfaces/ICrossVMBridgeCallable.sol";
import {CrossVMBridgeCallableUpgradeable} from "./CrossVMBridgeCallableUpgradeable.sol";

/**
 * @title CrossVMBridgeERC721Fulfillment
 * @dev Related to https://github.com/onflow/flips/issues/318[FLIP-318] Cross VM NFT implementations
 * on Flow in the context of Cadence-native NFTs. The following base contract should be implemented to
 * integrate with the Flow VM bridge connecting Cadence & EVM implementations so that the canonical
 * VM bridge may move the Cadence NFT into EVM in a mint/escrow pattern.
 */
abstract contract CrossVMBridgeERC721FulfillmentUpgradeable is ICrossVMBridgeERC721Fulfillment, CrossVMBridgeCallableUpgradeable, ERC721Upgradeable {

    /**
     * Initializes the bridge EVM address such that only the bridge COA can call privileged methods
     */
    function __CrossVMBridgeERC721Fulfillment_init(address vmBridgeAddress_) internal onlyInitializing {
        __CrossVMBridgeERC721Fulfillment_init_unchained(vmBridgeAddress_);
    }

    function __CrossVMBridgeERC721Fulfillment_init_unchained(address vmBridgeAddress_) internal onlyInitializing {
        _init_vm_bridge_address(vmBridgeAddress_);
    }

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
    function fulfillToEVM(address _to, uint256 _id, bytes memory _data) external onlyVMBridge {
        _beforeFulfillment(_to, _id, _data); // hook allowing implementation to perform pre-fulfillment validation
        if (_ownerOf(_id) == address(0)) {
            _mint(_to, _id); // Doesn't exist, mint the token
        } else {
            // Should be escrowed under vm bridge - transfer from escrow to recipient
            _requireEscrowed(_id);
            safeTransferFrom(vmBridgeAddress(), _to, _id);
        }
        _afterFulfillment(_to, _id, _data); // hook allowing implementation to perform post-fulfillment processing
        emit FulfilledToEVM(_to, _id);
    }

    /**
     * @dev Returns whether the token is currently escrowed under custody of the designated VM bridge
     * 
     * @param _id the ID of the token in question
     */
    function isEscrowed(uint256 _id) public view returns (bool) {
        return _ownerOf(_id) == vmBridgeAddress();
    }

    /**
     * @dev Allows a caller to determine the contract conforms to implemented interfaces
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override(CrossVMBridgeCallableUpgradeable, ERC721Upgradeable, IERC165) returns (bool) {
        return interfaceId == type(ICrossVMBridgeERC721Fulfillment).interfaceId
            || interfaceId == type(ICrossVMBridgeCallable).interfaceId
            || super.supportsInterface(interfaceId);
    }

    /**
     * @dev Internal method that reverts with FulfillmentFailedTokenNotEscrowed if the provided
     * token is not escrowed with the assigned vm bridge address as owner.
     * 
     * @param _id the token id that must be escrowed
     */
    function _requireEscrowed(uint256 _id) internal view {
        if (!isEscrowed(_id)) {
            revert FulfillmentFailedTokenNotEscrowed(_id, vmBridgeAddress());
        }
    }

    /**
     * @dev This internal method is included as a step implementations can override and have
     * executed in the default fullfillToEVM call.
     * 
     * @param _to address of the pending token recipient
     * @param _id the id of the token to be moved into EVM from Cadence
     * @param _data any encoded metadata passed by the corresponding Cadence NFT at the time of
     *      bridging into EVM
     */
    function _beforeFulfillment(address _to, uint256 _id, bytes memory _data) internal virtual {
        // No-op by default, meant to be overridden by implementations
    }

    /**
     * @dev This internal method is included as a step implementations can override and have
     * executed in the default fullfillToEVM call.
     * 
     * @param _to address of the pending token recipient
     * @param _id the id of the token to be moved into EVM from Cadence
     * @param _data any encoded metadata passed by the corresponding Cadence NFT at the time of
     *      bridging into EVM
     */
    function _afterFulfillment(address _to, uint256 _id, bytes memory _data) internal virtual {
        // No-op by default, meant to be overridden by implementations for things like processing
        // and setting metadata
    }
}

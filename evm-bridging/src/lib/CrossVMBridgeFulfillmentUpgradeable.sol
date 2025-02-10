// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.24;

import {CrossVMBridgeCallableUpgradeable} from "./CrossVMBridgeCallableUpgradeable.sol";
import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {ERC165Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/introspection/ERC165Upgradeable.sol";
import {ICrossVMBridgeFulfillment} from "../interfaces/ICrossVMBridgeFulfillment.sol";

abstract contract CrossVMBridgeFulfillmentUpgradeable is CrossVMBridgeCallableUpgradeable, ERC165Upgradeable, ERC721Upgradeable {

    error FulfillmentFailedTokenNotEscrowed(uint256 id, address escrowAddress);

    function __CrossVMBridgeFulfillment_init(address vmBridgeAddress_) internal onlyInitializing {
        __CrossVMBridgeFulfillment_init_unchained(vmBridgeAddress_);
    }

    function __CrossVMBridgeFulfillment_init_unchained(address vmBridgeAddress_) internal onlyInitializing {
        _init_vm_bridge_address(vmBridgeAddress_);
    }

    /**
     * @dev Fulfills the bridge request, minting (if non-existent) or transferring (if escrowed) the
     * token with the given ID to the provided address. For dynamic metadata handling between
     * Cadence & EVM, implementations should override and assign metadata as encoded from Cadence
     * side. If overriding, be sure to preserve the mint/escrow pattern as shown in the default
     * implementation.
     * 
     * @param _to address of the token recipient
     * @param _id the id of the token being moved into EVM from Cadence
     */
    function fulfillToEVM(address _to, uint256 _id, bytes memory /*_data*/) external onlyVMBridge {
        if (_ownerOf(_id) == address(0)) {
            _validateMint(_to, _id);
            _mint(_to, _id); // Doesn't exist, mint the token
        } else {
            // Should be escrowed under vm bridge - transfer from escrow to recipient
            _requireEscrowed(_id);
            safeTransferFrom(vmBridgeAddress(), _to, _id);
        }
    }

    function _validateMint(address _to, uint256 _id) internal view {
        // no-op, override in implementation if needed
    }

    /**
     * @dev Allows a caller to determine the contract conforms to the `ICrossVMFulfillment` interface
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC165Upgradeable, ERC721Upgradeable) returns (bool) {
        return interfaceId == type(ICrossVMBridgeFulfillment).interfaceId || super.supportsInterface(interfaceId);
    }

    /**
     * @dev Internal method that reverts with FulfillmentFailedTokenNotEscrowed if the provided
     * token is not escrowed with the assigned vm bridge address as owner.
     * 
     * @param _id the token id that must be escrowed
     */
    function _requireEscrowed(uint256 _id) internal view {
        address owner = _ownerOf(_id);
        address vmBridgeAddress_ = vmBridgeAddress();
        if (owner != vmBridgeAddress_) {
            revert FulfillmentFailedTokenNotEscrowed(_id, vmBridgeAddress_);
        }
    }
}

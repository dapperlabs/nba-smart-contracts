// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.24;

import {ICrossVMBridgeCallable} from "../interfaces/ICrossVMBridgeCallable.sol";
import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import {ERC165Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/introspection/ERC165Upgradeable.sol";

/**
 * @title CrossVMBridgeCallable
 * @dev A base contract intended for use in implementations on Flow, allowing a contract to define
 * access to the Cadence X EVM bridge on certain methods.
 */
abstract contract CrossVMBridgeCallableUpgradeable is ICrossVMBridgeCallable, ContextUpgradeable, ERC165Upgradeable {

    address private _vmBridgeAddress;

    /**
     * @dev Sets the bridge EVM address such that only the bridge COA can call the privileged methods
     */
    function _init_vm_bridge_address(address vmBridgeAddress_) internal {
        if (vmBridgeAddress_ == address(0)) {
            revert CrossVMBridgeCallableZeroInitialization();
        }
        _vmBridgeAddress = vmBridgeAddress_;
    }

    /**
     * @dev Modifier restricting access to the designated VM bridge EVM address 
     */
    modifier onlyVMBridge() {
        _checkVMBridgeAddress();
        _;
    }

    /**
     * @dev Returns the designated VM bridge’s EVM address
     */
    function vmBridgeAddress() public view virtual returns (address) {
        return _vmBridgeAddress;
    }

    /**
     * @dev Checks that msg.sender is the designated VM bridge address
     */
    function _checkVMBridgeAddress() internal view virtual {
        if (_vmBridgeAddress != _msgSender()) {
            revert CrossVMBridgeCallableUnauthorizedAccount(_msgSender());
        }
    }

    /**
     * @dev Allows a caller to determine the contract conforms to the `ICrossVMFulfillment` interface
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC165Upgradeable) returns (bool) {
        return interfaceId == type(ICrossVMBridgeCallable).interfaceId || super.supportsInterface(interfaceId);
    }
}

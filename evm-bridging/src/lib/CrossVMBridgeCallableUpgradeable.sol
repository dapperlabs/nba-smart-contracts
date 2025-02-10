// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.24;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import {ICrossVMBridgeFulfillment} from "../interfaces/ICrossVMBridgeFulfillment.sol";

/**
 * @title CrossVMBridgeCallable
 * @dev A base contract intended for use in implementations on Flow, allowing a contract to define
 *      access to the Cadence X EVM bridge on certain methods.
 */
abstract contract CrossVMBridgeCallableUpgradeable is ContextUpgradeable {

    address private _vmBridgeAddress;

    error CrossVMBridgeCallableZeroInitialization();
    error CrossVMBridgeCallableUnauthorizedAccount(address account);

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
     * @dev Checks that msg.sender is the designated vm bridge address
     */
    function _checkVMBridgeAddress() internal view virtual {
        if (vmBridgeAddress() != _msgSender()) {
            revert CrossVMBridgeCallableUnauthorizedAccount(_msgSender());
        }
    }
}

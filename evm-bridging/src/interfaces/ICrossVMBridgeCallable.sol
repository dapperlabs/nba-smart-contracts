// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.24;

/**
 * @title ICrossVMBridgeCallable
 * @dev An interface intended for use by implementations on Flow EVM, allowing a contract to define
 * access to the Cadence X EVM bridge on certain methods.
 */
interface ICrossVMBridgeCallable {

    /// @dev Should encounter when the vmBridgeAddress is initialized to 0x0
    error CrossVMBridgeCallableZeroInitialization();
    /// @dev Should encounter when a VM bridge privileged method is triggered by unauthorized caller
    error CrossVMBridgeCallableUnauthorizedAccount(address account);

    /**
     * @dev Returns the designated VM bridge’s EVM address
     */
    function vmBridgeAddress() external view returns (address);
}
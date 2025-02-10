// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {ERC165Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/introspection/ERC165Upgradeable.sol";
import {IBridgePermissions} from "../interfaces/IBridgePermissions.sol";

/**
 * @dev Contract for which implementation is checked by the Flow VM bridge as an opt-out mechanism
 * for non-standard asset contracts that wish to opt-out of bridging between Cadence & EVM. By
 * default, the VM bridge operates on a permissionless basis, meaning anyone can request an asset
 * be onboarded. However, some small subset of non-standard projects may wish to opt-out of this
 * and this contract provides a way to do so while also enabling future opt-in.
 *
 * Note: The Flow VM bridge checks for permissions at asset onboarding. If your asset has already
 * been onboarded, setting `permissions` to `false` will not affect movement between VMs.
 */
abstract contract BridgePermissionsUpgradeable is Initializable, ERC165Upgradeable, IBridgePermissions {
    // The permissions for the contract to allow or disallow bridging of its assets.
    bool private _permissions;

    function __BridgePermissions_init() internal onlyInitializing {
        _permissions = false;
    }

    /**
     * @dev See {IERC165-supportsInterface}.
     */
    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC165Upgradeable, IERC165) returns (bool) {
        return interfaceId == type(IBridgePermissions).interfaceId || super.supportsInterface(interfaceId);
    }

    /**
     * @dev Returns true if the contract allows bridging of its assets. Checked by the Flow VM
     *     bridge at asset onboarding to enable non-standard asset contracts to opt-out of bridging
     *     between Cadence & EVM. Implementing this contract opts out by default but can be
     *     overridden to opt-in or used in conjunction with a switch to enable opting in.
     */
    function allowsBridging() external view virtual returns (bool) {
        return _permissions;
    }

    /**
     * @dev Set the permissions for the contract to allow or disallow bridging of its assets.
     *
     * Emits a {PermissionsUpdated} event.
     */
    function _setPermissions(bool permissions) internal {
        _permissions = permissions;
        emit PermissionsUpdated(permissions);
    }
}

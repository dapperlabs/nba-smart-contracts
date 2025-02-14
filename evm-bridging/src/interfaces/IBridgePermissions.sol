// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

interface IBridgePermissions is IERC165 {
    /**
     * @dev Emitted when the permissions for the contract are updated.
     */
    event PermissionsUpdated(bool newPermissions);

    /**
     * @dev Returns true if the contract allows bridging of its assets.
     */
    function allowsBridging() external view returns (bool);
}

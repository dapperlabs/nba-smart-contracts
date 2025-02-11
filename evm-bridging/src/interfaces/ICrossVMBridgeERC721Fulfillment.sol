pragma solidity 0.8.24;

import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

interface ICrossVMBridgeERC721Fulfillment is IERC165 {
    function fulfillToEVM(address to, uint256 id, bytes memory data) external;
}

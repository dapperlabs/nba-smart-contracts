pragma solidity 0.8.24;

interface ICrossVMBridgeFulfillment {
    function fulfillToEVM(address to, uint256 id, bytes memory data) external;
    function vmBridgeAddress() external view returns (address);
}

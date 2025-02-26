pragma solidity 0.8.24;

interface ICrossVM {
    function getCadenceAddress() external view returns (string memory);
    function getCadenceIdentifier() external view returns (string memory);
}

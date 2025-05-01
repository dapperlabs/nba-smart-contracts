// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

contract TestContract {
    function testArrayEncoding(uint256[] calldata values) external pure
    returns (uint256[] memory)
    {
        return values;
    }
}
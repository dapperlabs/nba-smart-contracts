// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Script} from "forge-std/Script.sol";
import "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/src/Upgrades.sol";
import {TestNFTContract} from "../src/test-contracts/TestNFTContract.sol";

contract InitialTestingDeployScript is Script {
    function setUp() public {}

    function run() external returns (address, address) {
        // Start broadcast with deployer private key
        vm.startBroadcast(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        console.log("Deployer address:", msg.sender);

        // Set testnet contract initialization parameters
        address owner = msg.sender;
        string memory name = "Test NFT";
        string memory symbol = "TEST";
        string memory baseTokenURI = "https://api.cryptokitties.co/tokenuri/";
        string memory cadenceNFTAddress = "abcdef1234567890";
        string memory cadenceNFTIdentifier = "A.abcdef1234567890.TestNFT.NFT";
        string memory contractURI = 'data:application/json;utf8,{"name": "Name of NFT","description":"Description of NFT"}';
        address underlyingNftContractAddress = address(0x12345);
        address vmBridgeAddress = address(0x67890);

        // Deploy NFT contract using UUPS proxy for upgradeability
        address proxyAddr = Upgrades.deployUUPSProxy(
            "TestNFTContract.sol",
            abi.encodeCall(
                TestNFTContract.initialize,
                (
                    owner,
                    underlyingNftContractAddress,
                    vmBridgeAddress,
                    name,
                    symbol,
                    baseTokenURI,
                    cadenceNFTAddress,
                    cadenceNFTIdentifier,
                    contractURI
                )
            )
        );
        console.log("Proxy contract deployed at address:", proxyAddr);

        // Get implementation contract address
        address implementationAddr = Upgrades.getImplementationAddress(
            proxyAddr
        );
        console.log("Implementation contract deployed at address:", implementationAddr);

        // Stop broadcast and return implementation and proxy addresses
        vm.stopBroadcast();
        return (implementationAddr, proxyAddr);
    }
}
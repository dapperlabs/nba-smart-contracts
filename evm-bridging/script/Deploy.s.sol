// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Script} from "forge-std/Script.sol";
import "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/src/Upgrades.sol";
import {BridgedTopShotMoments} from "../src/BridgedTopShotMoments.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() external returns (address, address) {
        // Start broadcast with deployer private key
        vm.startBroadcast(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        console.log("Deployer address:", msg.sender);

        // Set testnet contract initialization parameters
        address owner = msg.sender;
        string memory name = "NBA Top Shot";
        string memory symbol = "TOPSHOT";
        string memory baseTokenURI = "https://api.cryptokitties.co/tokenuri/";
        string memory cadenceNFTAddress = "877931736ee77cff";
        string memory cadenceNFTIdentifier = "A.877931736ee77cff.NFT";
        string memory contractURI = "add-contract-URI-here";
        address underlyingNftContractAddress = address(0x12345);
        address vmBridgeAddress = address(0x67890);
        // Deploy NFT contract using UUPS proxy for upgradeability
        address proxyAddr = Upgrades.deployUUPSProxy(
            "BridgedTopShotMoments.sol",
            abi.encodeCall(
                BridgedTopShotMoments.initialize,
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
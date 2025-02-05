// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/src/Upgrades.sol";
import {BridgedTopShotMoments} from "../src/BridgedTopShotMoments.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() external returns (address, address) {
        // Start broadcast with deployer private key
        vm.startBroadcast(vm.envUint("DEPLOYER_PRIVATE_KEY"));
        console.log("Deployer address:", msg.sender);

        // Set contract initialization parameters
        address owner = msg.sender;
        string memory name = "Bidged NBA TopShot Moments";
        string memory symbol = "TOPSHOT";
        string memory cadenceNFTAddress = "cadenceNFTAddress";
        string memory cadenceNFTIdentifier = "cadenceNFTIdentifier";
        string memory contractURI = "contractURI";

        // Deploy NFT contract using UUPS proxy for upgradeability
        address proxyAddr = Upgrades.deployUUPSProxy(
            "BridgedTopShotMoments.sol",
            abi.encodeCall(
                BridgedTopShotMoments.initialize,
                (
                    owner,
                    name,
                    symbol,
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
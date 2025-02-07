// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/src/Upgrades.sol";
import {BridgedTopShotMoments} from "../src/BridgedTopShotMoments.sol";

contract BridgedTopShotMomentsTest is Test {
    address owner;
    string name;
    string symbol;
    string cadenceNFTAddress;
    string cadenceNFTIdentifier;
    string contractURI;
    BridgedTopShotMoments private nftContract;

    // Runs before each test
    function setUp() public {
        // Set initialization parameters
        owner = msg.sender;
        name = "name";
        symbol = "symbol";
        cadenceNFTAddress = "cadenceNFTAddress";
        cadenceNFTIdentifier = "cadenceNFTIdentifier";
        contractURI = "contractURI";

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

        // Set contract instance
        nftContract = BridgedTopShotMoments(proxyAddr);
    }

    function test_GetContractInfo() public view {
        assertEq(nftContract.owner(), owner);
        assertEq(nftContract.name(), name);
        assertEq(nftContract.symbol(), symbol);
        assertEq(nftContract.getCadenceAddress(), cadenceNFTAddress);
        assertEq(nftContract.getCadenceIdentifier(), cadenceNFTIdentifier);
        assertEq(nftContract.contractURI(), contractURI);
    }

    function test_MintNFT() public {
        address recipient = address(27);
        uint256 tokenId = 101;
        vm.startPrank(owner);
        nftContract.safeMint(recipient, tokenId, "MOCK_URI");
        vm.stopPrank();
        assertTrue(nftContract.exists(tokenId));
        assertEq(nftContract.ownerOf(tokenId), recipient);
        assertEq(nftContract.balanceOf(recipient), 1);
    }

    function test_UpdateTokenURI() public {
        uint256 tokenId = 100;
        vm.startPrank(owner);
        nftContract.safeMint(owner, tokenId, "MOCK_URI");

        string memory newURI = "NEW_URI";
        nftContract.updateTokenURI(tokenId, newURI);
        vm.stopPrank();
        assertEq(nftContract.tokenURI(tokenId), newURI);
    }

    function test_UpdateERC721Symbol() public {
        string memory _symbol = nftContract.symbol();
        assertEq(_symbol, symbol);

        string memory newSymbol = "NEW_SYMBOL";
        vm.startPrank(owner);
        nftContract.setSymbol(newSymbol);
        vm.stopPrank();

        _symbol = nftContract.symbol();
        assertEq(_symbol, newSymbol);
    }

    function test_RevertMintToZeroAddress() public {
        vm.startPrank(owner);
        vm.expectRevert();
        nftContract.safeMint(address(0), 1, "MOCK_URI");
        vm.stopPrank();
    }
}
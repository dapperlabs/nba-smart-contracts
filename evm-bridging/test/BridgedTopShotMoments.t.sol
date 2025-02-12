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

        // Mint NFT to recipient
        vm.startPrank(owner);
        nftContract.safeMint(recipient, tokenId, "MOCK_URI");
        vm.stopPrank();

        // Check NFT exists, owner, balance, and total supply
        assertTrue(nftContract.exists(tokenId));
        assertEq(nftContract.ownerOf(tokenId), recipient);
        assertEq(nftContract.balanceOf(recipient), 1);
        assertEq(nftContract.totalSupply(), 1);
    }

    function test_RevertMintToZeroAddress() public {
        // Expect revert when minting to zero address
        vm.startPrank(owner);
        vm.expectRevert();
        nftContract.safeMint(address(0), 1, "MOCK_URI");
        vm.stopPrank();
    }

    function test_TransferNFT() public {
        address account1 = address(27);
        address account2 = address(28);
        uint256 tokenId = 101;

        // Mint NFT to account1 and check balance
        vm.startPrank(owner);
        nftContract.safeMint(account1, tokenId, "MOCK_URI");
        vm.stopPrank();
        assertEq(nftContract.balanceOf(account1), 1);

        // Transfer NFT from account1 to account2 and check balances
        vm.startPrank(account1);
        nftContract.safeTransferFrom(account1, account2, tokenId);
        vm.stopPrank();
        assertEq(nftContract.balanceOf(account1), 0);
        assertEq(nftContract.balanceOf(account2), 1);
        assertEq(nftContract.ownerOf(tokenId), account2);
    }

    function test_ApproveNFT() public {
        address recipient = address(27);
        address operator = address(28);
        uint256 tokenId = 101;

        // Mint NFT to recipient and check balance and approval
        vm.startPrank(owner);
        nftContract.safeMint(recipient, tokenId, "MOCK_URI");
        vm.stopPrank();
        assertEq(nftContract.balanceOf(recipient), 1);
        assertEq(nftContract.getApproved(tokenId), address(0));

        // Approve operator for NFT and check approval
        vm.startPrank(recipient);
        nftContract.approve(operator, tokenId);
        vm.stopPrank();
        assertEq(nftContract.getApproved(tokenId), operator);

        // Transfer NFT from recipient to operator and check owner
        vm.startPrank(operator);
        nftContract.safeTransferFrom(recipient, address(29), tokenId);
        vm.stopPrank();
        assertEq(nftContract.ownerOf(tokenId), address(29));
    }

    function test_ApproveForAllNFTs() public {
        address recipient = address(27);
        address operator = address(28);
        uint256 tokenId1 = 101;
        uint256 tokenId2 = 102;

        // Mint NFTs to recipient and check balance and total supply
        vm.startPrank(owner);
        nftContract.safeMint(recipient, tokenId1, "MOCK_URI");
        nftContract.safeMint(recipient, tokenId2, "MOCK_URI");
        vm.stopPrank();
        assertEq(nftContract.balanceOf(recipient), 2);
        assertEq(nftContract.totalSupply(), 2);

        // Approve operator for all NFTs and check approval
        vm.startPrank(recipient);
        nftContract.setApprovalForAll(operator, true);
        vm.stopPrank();
        assertEq(nftContract.isApprovedForAll(recipient, operator), true);

        // Transfer NFTs from recipient to operator and check balances
        vm.startPrank(operator);
        nftContract.safeTransferFrom(recipient, address(30), tokenId1);
        nftContract.safeTransferFrom(recipient, address(30), tokenId2);
        vm.stopPrank();
        assertEq(nftContract.balanceOf(recipient), 0);
        assertEq(nftContract.balanceOf(address(30)), 2);
    }

    function test_BurnNFT() public {
        address recipient = address(27);
        uint256 tokenId = 101;

        // Mint NFT to recipient and check balance
        vm.startPrank(owner);
        nftContract.safeMint(recipient, tokenId, "MOCK_URI");
        vm.stopPrank();
        assertEq(nftContract.balanceOf(recipient), 1);

        // Burn NFT and check balance
        vm.startPrank(recipient);
        nftContract.burn(tokenId);
        vm.stopPrank();
        assertFalse(nftContract.exists(tokenId));
        assertEq(nftContract.balanceOf(recipient), 0);
    }

    function test_UpdateTokenURI() public {
        uint256 tokenId = 100;

        // Mint NFT to owner and check tokenURI
        vm.startPrank(owner);
        nftContract.safeMint(owner, tokenId, "MOCK_URI");
        vm.stopPrank();
        assertEq(nftContract.tokenURI(tokenId), "MOCK_URI");

        // Update tokenURI and check newURI
        string memory newURI = "NEW_URI";
        vm.startPrank(owner);
        nftContract.updateTokenURI(tokenId, newURI);
        vm.stopPrank();
        assertEq(nftContract.tokenURI(tokenId), newURI);
    }

    function test_UpdateERC721Symbol() public {
        // Check initial symbol
        string memory initialSymbol = nftContract.symbol();
        assertEq(initialSymbol, symbol);

        // Update symbol and check new symbol
        string memory newSymbol = "NEW_SYMBOL";
        vm.startPrank(owner);
        nftContract.setSymbol(newSymbol);
        vm.stopPrank();
        assertEq(nftContract.symbol(), newSymbol);
    }

    function test_TestTransferContractOwnership() public {
        address newOwner = address(27);
        vm.startPrank(owner);
        nftContract.transferOwnership(newOwner);
        vm.stopPrank();
        assertEq(nftContract.owner(), newOwner);
    }

    function test_RevertTransferContractOwnershipToZeroAddress() public {
        address newOwner = address(0);
        vm.startPrank(owner);
        vm.expectRevert();
        nftContract.transferOwnership(newOwner);
        vm.stopPrank();
    }

    function test_SetTransferValidator() public {
        // Check initial transfer validator
        assertEq(nftContract.getTransferValidator(), address(0));

        // Set transfer validator and check new validator
        address transferValidator = address(27);
        vm.startPrank(owner);
        nftContract.setTransferValidator(transferValidator);
        vm.stopPrank();
        assertEq(nftContract.getTransferValidator(), transferValidator);
    }

    function test_SetRoyaltyInfo() public {
        // Check initial royalty info
        assertEq(nftContract.royaltyAddress(), address(0));

        // Set royalty info and check new info
        uint96 royaltyBps = 500;
        vm.startPrank(owner);
        nftContract.setRoyaltyInfo(
            BridgedTopShotMoments.RoyaltyInfo({
                royaltyAddress: owner,
                royaltyBps: royaltyBps
            })
        );
        vm.stopPrank();
        assertEq(nftContract.royaltyAddress(), owner);
        assertEq(nftContract.royaltyBasisPoints(), royaltyBps);
    }
}

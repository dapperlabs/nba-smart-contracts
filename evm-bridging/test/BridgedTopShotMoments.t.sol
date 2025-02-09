// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/src/Upgrades.sol";
import {BridgedTopShotMoments} from "../src/BridgedTopShotMoments.sol";
import {ERC721} from "openzeppelin-contracts/contracts/token/ERC721/ERC721.sol";
import {Ownable} from "openzeppelin-contracts/contracts/access/Ownable.sol";
import {Strings} from "openzeppelin-contracts/contracts/utils/Strings.sol";

// Add this minimal ERC721 implementation for testing
contract UnderlyingERC721 is ERC721, Ownable {
    constructor(string memory name, string memory symbol) ERC721(name, symbol) Ownable(msg.sender) {}

    function safeMint(address to, uint256 tokenId) public onlyOwner {
        _safeMint(to, tokenId);
    }
}

contract BridgedTopShotMomentsTest is Test {
    address owner;
    string name;
    string symbol;
    string baseTokenURI;
    string cadenceNFTAddress;
    string cadenceNFTIdentifier;
    string contractURI;
    BridgedTopShotMoments private nftContract;
    UnderlyingERC721 private underlyingNftContract;
    address underlyingNftContractAddress;
    address underlyingNftContractOwner;
    uint256[] nftIDs;
    // Runs before each test
    function setUp() public {
        // Set owner
        owner = msg.sender;

        // Deploy underlying NFT contract and mint underlying NFTs to owner
        underlyingNftContractOwner = address(0x1111);
        nftIDs = [101, 102, 103];
        vm.startPrank(underlyingNftContractOwner);
        underlyingNftContract = new UnderlyingERC721("Underlying NFT", "UNFT");
        for (uint256 i = 0; i < nftIDs.length; i++) {
            underlyingNftContract.safeMint(owner, nftIDs[i]);
        }
        vm.stopPrank();
        assertEq(underlyingNftContract.balanceOf(owner), nftIDs.length);

        // Set NFT contract initialization parameters
        underlyingNftContractAddress = address(underlyingNftContract);
        name = "name";
        symbol = "symbol";
        baseTokenURI = "https://example.com/";
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
                    underlyingNftContractAddress,
                    name,
                    symbol,
                    baseTokenURI,
                    cadenceNFTAddress,
                    cadenceNFTIdentifier,
                    contractURI
                )
            )
        );

        // Set contract instance
        nftContract = BridgedTopShotMoments(proxyAddr);
    }

    /* Test contract initialization */

    function test_GetContractInfo() public view {
        assertEq(nftContract.owner(), owner);
        assertEq(nftContract.name(), name);
        assertEq(nftContract.symbol(), symbol);
        assertEq(nftContract.getCadenceAddress(), cadenceNFTAddress);
        assertEq(nftContract.getCadenceIdentifier(), cadenceNFTIdentifier);
        assertEq(nftContract.contractURI(), contractURI);
        assertEq(address(nftContract.underlying()), underlyingNftContractAddress);
    }


    /* Test ERC721Wrapper operations */

    function test_WrapNFTs() public {
        // Approve and wrap NFT
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();
        assertEq(nftContract.balanceOf(owner), nftIDs.length);
        assertEq(underlyingNftContract.balanceOf(owner), 0);
        for (uint256 i = 0; i < nftIDs.length; i++) {
            assertEq(nftContract.ownerOf(nftIDs[i]), owner);
        }
    }

    function test_RevertWrapNFTsNotApproved() public {
        vm.startPrank(owner);
        vm.expectRevert();
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();
    }

    function test_RevertWrapNFTsZeroAddress() public {
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        vm.expectRevert();
        nftContract.depositFor(address(0), nftIDs);
        vm.stopPrank();
    }

    function test_UnwrapNFTs() public {
        // Approve and wrap NFT
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();

        // Unwrap NFT
        vm.startPrank(owner);
        nftContract.withdrawTo(owner, nftIDs);
        vm.stopPrank();
        assertEq(underlyingNftContract.balanceOf(owner), nftIDs.length);
        assertEq(underlyingNftContract.balanceOf(address(nftContract)), 0);
        for (uint256 i = 0; i < nftIDs.length; i++) {
            assertEq(underlyingNftContract.ownerOf(nftIDs[i]), owner);
        }
    }

    /* Test core ERC721 operations */

    function test_TransferNFT() public {
        address recipient = address(27);

        // Approve and wrap NFT
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();

        // Transfer NFT from account1 to account2 and check balances
        vm.startPrank(owner);
        nftContract.safeTransferFrom(owner, recipient, nftIDs[0]);
        vm.stopPrank();
        assertEq(nftContract.balanceOf(owner), nftIDs.length - 1);
        assertEq(nftContract.balanceOf(recipient), 1);
        assertEq(nftContract.ownerOf(nftIDs[0]), recipient);
    }

    function test_ApproveNFT() public {
        address operator = address(28);

        // Approve and wrap NFT
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();

        // Approve operator for NFT and check approval
        vm.startPrank(owner);
        nftContract.approve(operator, nftIDs[0]);
        vm.stopPrank();
        assertEq(nftContract.getApproved(nftIDs[0]), operator);
    }

    function test_BurnNFT() public {
        // Approve and wrap NFT
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();

        // Burn NFT and check balance
        vm.startPrank(owner);
        nftContract.burn(nftIDs[0]);
        vm.stopPrank();
        assertFalse(nftContract.exists(nftIDs[0]));
        assertEq(nftContract.balanceOf(owner), nftIDs.length - 1);
    }

    function test_UpdateBaseTokenURI() public {
        string memory newBaseTokenURI = "NEW_BASE_URI";

        // Approve and wrap NFT
        vm.startPrank(owner);
        underlyingNftContract.setApprovalForAll(address(nftContract), true);
        nftContract.depositFor(owner, nftIDs);
        vm.stopPrank();
        assertEq(nftContract.tokenURI(nftIDs[0]), string(abi.encodePacked(baseTokenURI, Strings.toString(nftIDs[0]))));

        // Update tokenURI and check newURI
        vm.startPrank(owner);
        nftContract.setBaseTokenURI(newBaseTokenURI);
        vm.stopPrank();
        assertEq(nftContract.tokenURI(nftIDs[0]), string(abi.encodePacked(newBaseTokenURI, Strings.toString(nftIDs[0]))));
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

    /* Test Creator Token operations */

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

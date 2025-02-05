// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {IERC721Metadata} from "@openzeppelin/contracts/token/ERC721/extensions/IERC721Metadata.sol";
import {IERC721Enumerable} from "@openzeppelin/contracts/token/ERC721/extensions/IERC721Enumerable.sol";

import {ERC721Upgradeable} from "openzeppelin-contracts-upgradeable/contracts/token/ERC721/ERC721Upgradeable.sol";
import {ERC721URIStorageUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import {ERC721EnumerableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/token/ERC721/extensions/ERC721EnumerableUpgradeable.sol";
import {ERC721BurnableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/token/ERC721/extensions/ERC721BurnableUpgradeable.sol";
import {OwnableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {Initializable} from "openzeppelin-contracts-upgradeable/contracts/proxy/utils/Initializable.sol";

import {ICrossVM} from "./interfaces/ICrossVM.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {ERC165} from "@openzeppelin/contracts/utils/introspection/ERC165.sol";

contract BridgedTopShotMoments is
    Initializable,
    ERC721Upgradeable,
    ERC721URIStorageUpgradeable,
    ERC721BurnableUpgradeable,
    ERC721EnumerableUpgradeable,
    OwnableUpgradeable,
    ICrossVM
{
    string public cadenceNFTAddress;
    string public cadenceNFTIdentifier;
    string public contractMetadata;
    string private _customSymbol;

    function initialize(
        address owner,
        string memory name_,
        string memory symbol_,
        string memory _cadenceNFTAddress,
        string memory _cadenceNFTIdentifier,
        string memory _contractMetadata) public initializer
    {
        __ERC721_init(name_, symbol_);
        __Ownable_init(owner);
        _customSymbol = symbol_;
        cadenceNFTAddress = _cadenceNFTAddress;
        cadenceNFTIdentifier = _cadenceNFTIdentifier;
        contractMetadata = _contractMetadata;
    }

    function getCadenceAddress() external view returns (string memory) {
        return cadenceNFTAddress;
    }

    function getCadenceIdentifier() external view returns (string memory) {
        return cadenceNFTIdentifier;
    }

    function symbol() public view override returns (string memory) {
        return _customSymbol;
    }

    function safeMint(address to, uint256 tokenId, string memory uri) public onlyOwner {
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
    }

    function updateTokenURI(uint256 tokenId, string memory uri) public onlyOwner {
        _setTokenURI(tokenId, uri);
    }

    function setSymbol(string memory newSymbol) public onlyOwner {
        _setSymbol(newSymbol);
    }

    function contractURI() public view returns (string memory) {
        return contractMetadata;
    }

    function tokenURI(uint256 tokenId) public view override(ERC721Upgradeable, ERC721URIStorageUpgradeable) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721Upgradeable, ERC721EnumerableUpgradeable, ERC721URIStorageUpgradeable)
        returns (bool)
    {
        return interfaceId == type(IERC165).interfaceId || interfaceId == type(IERC721Metadata).interfaceId
            || interfaceId == type(IERC721Enumerable).interfaceId || interfaceId == type(ERC721BurnableUpgradeable).interfaceId
            || interfaceId == type(OwnableUpgradeable).interfaceId
            || super.supportsInterface(interfaceId);
    }

    function exists(uint256 tokenId) public view returns (bool) {
        return _ownerOf(tokenId) != address(0);
    }

    function _setSymbol(string memory newSymbol) internal {
        _customSymbol = newSymbol;
    }

    function _update(address to, uint256 tokenId, address auth)
        internal
        override(ERC721Upgradeable, ERC721EnumerableUpgradeable)
        returns (address)
    {
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(address account, uint128 value) internal override(ERC721Upgradeable, ERC721EnumerableUpgradeable) {
        super._increaseBalance(account, value);
    }
}

// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {ERC721URIStorageUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import {ERC721EnumerableUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721EnumerableUpgradeable.sol";
import {ERC721BurnableUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721BurnableUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import {IERC721Metadata} from "@openzeppelin/contracts/token/ERC721/extensions/IERC721Metadata.sol";
import {IERC721Enumerable} from "@openzeppelin/contracts/token/ERC721/extensions/IERC721Enumerable.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {IERC2981} from "@openzeppelin/contracts/interfaces/IERC2981.sol";
import {ERC165} from "@openzeppelin/contracts/utils/introspection/ERC165.sol";

import {ICreatorToken, ILegacyCreatorToken} from "./interfaces/ICreatorToken.sol";
import {ITransferValidator721} from "./interfaces/ITransferValidator.sol";
import {ERC721TransferValidator} from "./lib/ERC721TransferValidator.sol";

import {ICrossVM} from "./interfaces/ICrossVM.sol";

// Initial draft version of the BridgedTopShotMoments contract
contract BridgedTopShotMoments is
    Initializable,
    ERC721Upgradeable,
    ERC721URIStorageUpgradeable,
    ERC721BurnableUpgradeable,
    ERC721EnumerableUpgradeable,
    OwnableUpgradeable,
    ERC721TransferValidator,
    ICrossVM
{
    string public cadenceNFTAddress;
    string public cadenceNFTIdentifier;
    string public contractMetadata;
    string private _customSymbol;
    RoyaltyInfo private _royaltyInfo;

    error InvalidRoyaltyBasisPoints(uint256 basisPoints);
    error RoyaltyAddressCannotBeZeroAddress();
    event RoyaltyInfoUpdated(address receiver, uint256 bps);
    struct RoyaltyInfo {
        address royaltyAddress;
        uint96 royaltyBps;
    }

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
            || interfaceId == type(OwnableUpgradeable).interfaceId || interfaceId == type(ICrossVM).interfaceId
            || interfaceId == type(ICreatorToken).interfaceId || interfaceId == type(ILegacyCreatorToken).interfaceId
            || interfaceId == type(IERC2981).interfaceId || super.supportsInterface(interfaceId);
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
        // Add the beforeTokenTransfer hook
        _beforeTokenTransfer(_ownerOf(tokenId), to, tokenId);

        // Call parent implementation
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(address account, uint128 value) internal override(ERC721Upgradeable, ERC721EnumerableUpgradeable) {
        super._increaseBalance(account, value);
    }

    function setRoyaltyInfo(RoyaltyInfo calldata newInfo) external onlyOwner {
        // Revert if the new royalty address is the zero address.
        if (newInfo.royaltyAddress == address(0)) {
            revert RoyaltyAddressCannotBeZeroAddress();
        }

        // Revert if the new basis points is greater than 10_000.
        if (newInfo.royaltyBps > 10_000) {
            revert InvalidRoyaltyBasisPoints(newInfo.royaltyBps);
        }

        // Set the new royalty info.
        _royaltyInfo = newInfo;

        // Emit an event with the updated params.
        emit RoyaltyInfoUpdated(newInfo.royaltyAddress, newInfo.royaltyBps);
    }

    function royaltyAddress() external view returns (address) {
        return _royaltyInfo.royaltyAddress;
    }

    function royaltyBasisPoints() external view returns (uint256) {
        return _royaltyInfo.royaltyBps;
    }

    /**
     * @dev Implements the IERC2981 interface.
     */
    function royaltyInfo(
        uint256 /* _tokenId */,
        uint256 _salePrice
    ) external view returns (address receiver, uint256 royaltyAmount) {
        // Put the royalty info on the stack for more efficient access.
        RoyaltyInfo storage info = _royaltyInfo;

        // Set the royalty amount to the sale price times the royalty basis
        // points divided by 10_000.
        royaltyAmount = (_salePrice * info.royaltyBps) / 10_000;

        // Set the receiver of the royalty.
        receiver = info.royaltyAddress;
    }

    function getTransferValidationFunction()
        external
        pure
        returns (bytes4 functionSignature, bool isViewFunction)
    {
        functionSignature = ITransferValidator721.validateTransfer.selector;
        isViewFunction = false;
    }

    function setTransferValidator(address newValidator) external onlyOwner {
        _setTransferValidator(newValidator);
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 startTokenId
    ) internal virtual {
        if (from != address(0) && to != address(0)) {
            // Call the transfer validator if one is set.
            address transferValidator = _transferValidator;
            if (transferValidator != address(0)) {
                ITransferValidator721(transferValidator).validateTransfer(
                    msg.sender,
                    from,
                    to,
                    startTokenId
                );
            }
        }
    }
}

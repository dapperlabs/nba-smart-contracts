// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {ERC721EnumerableUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721EnumerableUpgradeable.sol";
import {ERC721BurnableUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721BurnableUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {ERC721WrapperUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721WrapperUpgradeable.sol";

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";

import {IERC721Metadata} from "@openzeppelin/contracts/token/ERC721/extensions/IERC721Metadata.sol";
import {IERC721Enumerable} from "@openzeppelin/contracts/token/ERC721/extensions/IERC721Enumerable.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {IERC2981} from "@openzeppelin/contracts/interfaces/IERC2981.sol";
import {ERC165} from "@openzeppelin/contracts/utils/introspection/ERC165.sol";

import {ICreatorToken, ILegacyCreatorToken} from "./interfaces/ICreatorToken.sol";
import {ITransferValidator721} from "./interfaces/ITransferValidator.sol";
import {ERC721TransferValidator} from "./lib/ERC721TransferValidator.sol";

import {ICrossVM} from "./interfaces/ICrossVM.sol";
import {BridgePermissionsUpgradeable} from "./lib/BridgePermissionsUpgradeable.sol";
import {CrossVMBridgeERC721FulfillmentUpgradeable} from "./lib/CrossVMBridgeERC721FulfillmentUpgradeable.sol";

/**
 * @title ERC-721 BridgedTopShotMoments
 * @notice An upgradeable ERC721 contract for bridged NBA Top Shot Moments
 * @dev This contract implements multiple features:
 * - ERC721 standard functionality with enumeration and burning capabilities
 * - Wrapper functionality to handle NFTs from bridged-deployed contract
 * - Fulfillment functionality for Flow -> EVM bridging, once bridge onboarding allowed
 * - Cross-VM compatibility for Flow <-> EVM bridging
 * - Royalty management for secondary sales
 */
contract BridgedTopShotMoments is
    Initializable,
    ERC721Upgradeable,
    ERC721BurnableUpgradeable,
    ERC721EnumerableUpgradeable,
    OwnableUpgradeable,
    ERC721WrapperUpgradeable,
    ERC721TransferValidator,
    CrossVMBridgeERC721FulfillmentUpgradeable,
    BridgePermissionsUpgradeable,
    ICrossVM,
    IERC2981
{
    // Cadence-specific identifiers for cross-chain bridging
    string public cadenceNFTAddress;
    string public cadenceNFTIdentifier;

    // Metadata-related fields
    string public contractMetadata;
    string private _customSymbol;
    string private _baseTokenURI;

    // Royalty configuration for secondary sales
    RoyaltyInfo private _royaltyInfo;

    // Error declarations
    error InvalidRoyaltyBasisPoints(uint256 basisPoints);
    error RoyaltyAddressCannotBeZeroAddress();
    error InvalidUnderlyingTokenAddress();

    // Event declarations
    event RoyaltyInfoUpdated(address receiver, uint256 bps);
    event ContractURIUpdated();
    event MetadataUpdate(uint256 tokenId);

    /**
     * @notice Stores royalty configuration for secondary sales
     * @dev royaltyBps is in basis points (1/100th of a percent)
     * e.g., 500 = 5%, max value is 10000 = 100%
     */
    struct RoyaltyInfo {
        address royaltyAddress;
        uint96 royaltyBps;
    }

    /**
     * @dev Initializes the contract.
     */
    function initialize(
        address owner,
        address underlyingNftContractAddress,
        address vmBridgeAddress,
        string memory name_,
        string memory symbol_,
        string memory baseTokenURI_,
        string memory _cadenceNFTAddress,
        string memory _cadenceNFTIdentifier,
        string memory _contractMetadata
    ) public initializer {
        if (underlyingNftContractAddress == address(0)) {
            revert InvalidUnderlyingTokenAddress();
        }
        __ERC721_init(name_, symbol_);
        __Ownable_init(owner);
        __ERC721Wrapper_init(IERC721(underlyingNftContractAddress));
        __CrossVMBridgeERC721Fulfillment_init(vmBridgeAddress);
        __BridgePermissions_init();
        _customSymbol = symbol_;
        _baseTokenURI = baseTokenURI_;
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

    function contractURI() public view returns (string memory) {
        return contractMetadata;
    }

    function setSymbol(string memory newSymbol) public onlyOwner {
        _setSymbol(newSymbol);
    }

    /**
     * @notice Sets the contract URI, whether an offchain metadata URL or a JSON object
     * (i.e. `data:application/json;utf8,{"name":"...","description":"..."}`).
     */
    function setContractURI(string memory newMetadata) external onlyOwner {
        contractMetadata = newMetadata;

        // Indicate that the metadata has been updated (https://docs.opensea.io/docs/contract-level-metadata)
        emit ContractURIUpdated();
    }

    function setBaseTokenURI(string memory newBaseTokenURI) public onlyOwner {
        _baseTokenURI = newBaseTokenURI;

        // Indicate that the metadata has been updated (https://docs.opensea.io/docs/metadata-standards#metadata-updates)
        emit MetadataUpdate(type(uint256).max);
    }

    /**
     * @notice Returns the token URI for a given token ID.
     */
    function tokenURI(uint256 tokenId) public view override(ERC721Upgradeable) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721Upgradeable, ERC721EnumerableUpgradeable, BridgePermissionsUpgradeable, CrossVMBridgeERC721FulfillmentUpgradeable, IERC165)
        returns (bool)
    {
        return interfaceId == type(IERC165).interfaceId || interfaceId == type(IERC721Metadata).interfaceId
            || interfaceId == type(IERC721Enumerable).interfaceId || interfaceId == type(ERC721BurnableUpgradeable).interfaceId
            || interfaceId == type(OwnableUpgradeable).interfaceId || interfaceId == type(ICrossVM).interfaceId
            || interfaceId == type(ICreatorToken).interfaceId || interfaceId == type(ILegacyCreatorToken).interfaceId
            || interfaceId == type(IERC2981).interfaceId || super.supportsInterface(interfaceId);
    }

    function _setSymbol(string memory newSymbol) internal {
        _customSymbol = newSymbol;
    }

    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
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

    function setBridgePermissions(bool permissions) external onlyOwner {
        _setPermissions(permissions);
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

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {ERC165} from "@openzeppelin/contracts/utils/introspection/ERC165.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {IERC721} from "@openzeppelin/contracts/interfaces/IERC721.sol";

import {IERC7857, SealedKeyEntry} from "./interfaces/IERC7857.sol";
import {IERC7857Authorize} from "./interfaces/IERC7857Authorize.sol";
import {IERC7857Cloneable} from "./interfaces/IERC7857Cloneable.sol";
import {IntelligentData} from "./interfaces/IERC7857Metadata.sol";
import {IERC7857DataVerifier, TransferValidityProof} from "./interfaces/IERC7857DataVerifier.sol";

/// @title PitDeskID
/// @notice Aristotle (16661) desk identity. Production ERC-7857 interface IDs.
///         Mint / own / encrypted intelligence hash / authorizeUsage / revoke are live.
///         iTransferFrom / iCloneFrom revert: Foundation TEE attestor is Galileo-only.
///         ERC-721 transferFrom is disabled (ERC7857UseITransferFrom).
///         Self-appointed TEE oracles are forbidden — that is not a Foundation transfer.
contract PitDeskID is ERC721, IERC7857, IERC7857Authorize, IERC7857Cloneable {
    error AttestorNotOnAristotle();
    error NotDeskOwner();
    error ZeroAddress();
    error EmptyURI();
    error ZeroDataHash();
    error MaxAuthorized();

    uint256 public constant MAX_AUTHORIZED = 16;
    uint256 public nextId = 1;

    mapping(uint256 => IntelligentData[]) private _idata;
    mapping(uint256 => string) private _uri;
    mapping(uint256 => address[]) private _authz;
    mapping(uint256 => mapping(address => uint256)) private _authzIndex; // 1-based

    constructor() ERC721("PIT Desk ID", "PITDESK") {}

    /// @notice Mint a desk NFT to `to` with encrypted-intelligence pointer.
    /// @dev `dataHash` is keccak256 of the ciphertext (or Storage root commitment).
    ///      `uri` may be ipfs:// or 0g://<root>. Never store plaintext book on-chain.
    function mint(address to, string calldata uri, bytes32 dataHash, string calldata dataDescription)
        external
        returns (uint256 tokenId)
    {
        if (to == address(0)) revert ZeroAddress();
        if (bytes(uri).length == 0) revert EmptyURI();
        if (dataHash == bytes32(0)) revert ZeroDataHash();
        tokenId = nextId++;
        _safeMint(to, tokenId);
        _uri[tokenId] = uri;
        _idata[tokenId].push(IntelligentData({dataDescription: dataDescription, dataHash: dataHash}));
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId);
        return _uri[tokenId];
    }

    function intelligentDatasOf(uint256 tokenId) external view returns (IntelligentData[] memory) {
        _requireOwned(tokenId);
        return _idata[tokenId];
    }

    function verifier() external pure returns (IERC7857DataVerifier) {
        return IERC7857DataVerifier(address(0));
    }

    function authorizeUsage(uint256 tokenId, address user) public {
        _onlyOwner(tokenId);
        if (user == address(0)) revert IERC7857Authorize.ERC7857InvalidAuthorizedUser(user);
        if (_authzIndex[tokenId][user] != 0) revert IERC7857Authorize.ERC7857AlreadyAuthorized();
        if (_authz[tokenId].length >= MAX_AUTHORIZED) revert IERC7857Authorize.ERC7857TooManyAuthorizedUsers();
        _authz[tokenId].push(user);
        _authzIndex[tokenId][user] = _authz[tokenId].length;
        emit AuthorizationGranted(msg.sender, user, tokenId);
    }

    function batchAuthorizeUsage(uint256 tokenId, address[] calldata users) external {
        for (uint256 i; i < users.length; ++i) {
            authorizeUsage(tokenId, users[i]);
        }
    }

    function revokeAuthorization(uint256 tokenId, address user) public {
        _onlyOwner(tokenId);
        uint256 idx = _authzIndex[tokenId][user];
        if (idx == 0) revert IERC7857Authorize.ERC7857NotAuthorized();
        address[] storage list = _authz[tokenId];
        uint256 last = list.length;
        if (idx != last) {
            address moved = list[last - 1];
            list[idx - 1] = moved;
            _authzIndex[tokenId][moved] = idx;
        }
        list.pop();
        delete _authzIndex[tokenId][user];
        emit AuthorizationRevoked(msg.sender, user, tokenId);
    }

    function clearAuthorizedUsers(uint256 tokenId) external {
        _onlyOwner(tokenId);
        address[] storage list = _authz[tokenId];
        for (uint256 i; i < list.length; ++i) {
            delete _authzIndex[tokenId][list[i]];
        }
        delete _authz[tokenId];
        emit AuthorizationCleared(msg.sender, tokenId);
    }

    function authorizedUsersOf(uint256 tokenId) external view returns (address[] memory) {
        _requireOwned(tokenId);
        return _authz[tokenId];
    }

    /// @notice Load-bearing off-chain gate: owner OR authorized user may use the desk.
    function isAuthorized(uint256 tokenId, address user) external view returns (bool) {
        if (_ownerOf(tokenId) == user) return true;
        return _authzIndex[tokenId][user] != 0;
    }

    function iTransferFrom(address, address, uint256, TransferValidityProof[] calldata)
        external
        pure
        returns (SealedKeyEntry[] memory)
    {
        revert AttestorNotOnAristotle();
    }

    function iCloneFrom(address, address, uint256, TransferValidityProof[] calldata)
        external
        pure
        returns (uint256)
    {
        revert AttestorNotOnAristotle();
    }

    function transferFrom(address, address, uint256) public pure override(ERC721, IERC721) {
        revert IERC7857.ERC7857UseITransferFrom();
    }

    function safeTransferFrom(address, address, uint256, bytes memory) public pure override(ERC721, IERC721) {
        revert IERC7857.ERC7857UseITransferFrom();
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, IERC165)
        returns (bool)
    {
        return interfaceId == type(IERC7857).interfaceId || interfaceId == type(IERC7857Authorize).interfaceId
            || interfaceId == type(IERC7857Cloneable).interfaceId || super.supportsInterface(interfaceId);
    }

    function _onlyOwner(uint256 tokenId) internal view {
        if (_ownerOf(tokenId) != msg.sender) revert NotDeskOwner();
    }
}

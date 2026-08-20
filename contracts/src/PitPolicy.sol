// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title PitPolicy
/// @notice Pins a workspace policy hash on the selected 0G chain. The host engine
///         is still authoritative; this contract is the public commitment.
contract PitPolicy {
    error NotWorkspace();
    error EmptyHash();

    event PolicyPinned(bytes32 indexed workspaceId, address indexed wallet, bytes32 policyHash, uint64 version);

    mapping(bytes32 => bytes32) public policyHashOf;
    mapping(bytes32 => address) public walletOf;
    mapping(bytes32 => uint64) public versionOf;

    function pin(bytes32 workspaceId, bytes32 policyHash) external {
        if (workspaceId == bytes32(0) || policyHash == bytes32(0)) revert EmptyHash();
        address existing = walletOf[workspaceId];
        if (existing != address(0) && existing != msg.sender) revert NotWorkspace();
        walletOf[workspaceId] = msg.sender;
        policyHashOf[workspaceId] = policyHash;
        uint64 v = versionOf[workspaceId] + 1;
        versionOf[workspaceId] = v;
        emit PolicyPinned(workspaceId, msg.sender, policyHash, v);
    }
}

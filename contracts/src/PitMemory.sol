// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract PitMemory {
    error Empty();
    error AlreadyPinned();

    event MemoryPinned(bytes32 indexed objectId, address indexed wallet, bytes32 root);

    mapping(bytes32 => bytes32) public rootOf;
    mapping(bytes32 => address) public ownerOf;

    function pin(bytes32 objectId, bytes32 storageRoot) external {
        if (objectId == bytes32(0) || storageRoot == bytes32(0)) revert Empty();
        if (rootOf[objectId] != bytes32(0)) revert AlreadyPinned();
        rootOf[objectId] = storageRoot;
        ownerOf[objectId] = msg.sender;
        emit MemoryPinned(objectId, msg.sender, storageRoot);
    }
}

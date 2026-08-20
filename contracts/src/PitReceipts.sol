// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract PitReceipts {
    error AlreadyFiled();
    error Empty();

    event ReceiptFiled(bytes32 indexed previewHash, address indexed wallet, bytes32 storageRoot);

    mapping(bytes32 => bool) public filed;
    mapping(bytes32 => address) public filer;
    mapping(bytes32 => bytes32) public rootOf;

    function file(bytes32 previewHash, bytes32 storageRoot) external {
        if (previewHash == bytes32(0) || storageRoot == bytes32(0)) revert Empty();
        if (filed[previewHash]) revert AlreadyFiled();
        filed[previewHash] = true;
        filer[previewHash] = msg.sender;
        rootOf[previewHash] = storageRoot;
        emit ReceiptFiled(previewHash, msg.sender, storageRoot);
    }
}

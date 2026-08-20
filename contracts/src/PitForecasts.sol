// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract PitForecasts {
    error NotCommitter();
    error Empty();
    error AlreadyResolved();

    event Committed(bytes32 indexed forecastId, address indexed wallet, bytes32 commitHash);
    event Resolved(bytes32 indexed forecastId, bytes32 outcomeHash);

    mapping(bytes32 => address) public committer;
    mapping(bytes32 => bytes32) public commitHashOf;
    mapping(bytes32 => bytes32) public outcomeOf;

    function commit(bytes32 forecastId, bytes32 commitHash) external {
        if (forecastId == bytes32(0) || commitHash == bytes32(0)) revert Empty();
        address existing = committer[forecastId];
        if (existing != address(0) && existing != msg.sender) revert NotCommitter();
        committer[forecastId] = msg.sender;
        commitHashOf[forecastId] = commitHash;
        emit Committed(forecastId, msg.sender, commitHash);
    }

    function resolve(bytes32 forecastId, bytes32 outcomeHash) external {
        if (committer[forecastId] != msg.sender) revert NotCommitter();
        if (outcomeOf[forecastId] != bytes32(0)) revert AlreadyResolved();
        outcomeOf[forecastId] = outcomeHash;
        emit Resolved(forecastId, outcomeHash);
    }
}

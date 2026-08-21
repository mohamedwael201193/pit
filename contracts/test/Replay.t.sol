// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitReceipts} from "../src/PitReceipts.sol";
import {PitMemory} from "../src/PitMemory.sol";

contract ReplayTest is Test {
    PitReceipts internal rec;
    PitMemory internal mem;
    address internal alice = address(0xA11CE);

    function setUp() public {
        rec = new PitReceipts();
        mem = new PitMemory();
    }

    function test_receiptReplayReverts() public {
        bytes32 h = keccak256("preview");
        bytes32 root = keccak256("root");
        vm.prank(alice);
        rec.file(h, root);
        vm.prank(alice);
        vm.expectRevert(PitReceipts.AlreadyFiled.selector);
        rec.file(h, root);
    }

    function test_memoryReplayReverts() public {
        bytes32 id = keccak256("obj");
        bytes32 root = keccak256("root");
        vm.prank(alice);
        mem.pin(id, root);
        vm.prank(alice);
        vm.expectRevert(PitMemory.AlreadyPinned.selector);
        mem.pin(id, root);
    }
}

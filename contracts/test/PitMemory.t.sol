// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitMemory} from "../src/PitMemory.sol";

contract PitMemoryTest is Test {
    PitMemory internal mem;
    address internal alice = address(0xA11CE);

    function setUp() public {
        mem = new PitMemory();
    }

    function test_pinOnce() public {
        bytes32 id = keccak256("obj");
        bytes32 root = keccak256("root");
        vm.prank(alice);
        mem.pin(id, root);
        assertEq(mem.ownerOf(id), alice);
        vm.prank(alice);
        vm.expectRevert(PitMemory.AlreadyPinned.selector);
        mem.pin(id, root);
    }
}

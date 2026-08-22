// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitReceipts} from "../src/PitReceipts.sol";
import {PitForecasts} from "../src/PitForecasts.sol";
import {PitMemory} from "../src/PitMemory.sol";

contract IsolationTest is Test {
    PitReceipts internal rec;
    PitForecasts internal fc;
    PitMemory internal mem;
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    function setUp() public {
        rec = new PitReceipts();
        fc = new PitForecasts();
        mem = new PitMemory();
    }

    function test_bobCannotReplayAliceReceipt() public {
        bytes32 h = keccak256("preview-a");
        vm.prank(alice);
        rec.file(h, keccak256("root-a"));
        vm.prank(bob);
        vm.expectRevert(PitReceipts.AlreadyFiled.selector);
        rec.file(h, keccak256("root-b"));
        assertEq(rec.filer(h), alice);
    }

    function test_bobCannotResolveAliceForecast() public {
        bytes32 id = keccak256("f-a");
        vm.prank(alice);
        fc.commit(id, keccak256("c"));
        vm.prank(bob);
        vm.expectRevert(PitForecasts.NotCommitter.selector);
        fc.resolve(id, keccak256("o"));
    }

    function test_bobCannotOverwriteAliceMemory() public {
        bytes32 id = keccak256("obj-a");
        vm.prank(alice);
        mem.pin(id, keccak256("root-a"));
        vm.prank(bob);
        vm.expectRevert(PitMemory.AlreadyPinned.selector);
        mem.pin(id, keccak256("root-b"));
        assertEq(mem.ownerOf(id), alice);
    }
}

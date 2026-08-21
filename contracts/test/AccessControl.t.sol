// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitPolicy} from "../src/PitPolicy.sol";
import {PitReceipts} from "../src/PitReceipts.sol";
import {PitForecasts} from "../src/PitForecasts.sol";
import {PitMemory} from "../src/PitMemory.sol";

contract AccessControlTest is Test {
    PitPolicy internal pol;
    PitReceipts internal rec;
    PitForecasts internal fc;
    PitMemory internal mem;
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    function setUp() public {
        pol = new PitPolicy();
        rec = new PitReceipts();
        fc = new PitForecasts();
        mem = new PitMemory();
    }

    function test_policyIsolation() public {
        bytes32 ws = keccak256("ws-a");
        vm.prank(alice);
        pol.pin(ws, keccak256("p1"));
        vm.prank(bob);
        vm.expectRevert(PitPolicy.NotWorkspace.selector);
        pol.pin(ws, keccak256("p2"));
    }

    function test_forecastNotCommitter() public {
        bytes32 id = keccak256("f");
        vm.prank(alice);
        fc.commit(id, keccak256("c"));
        vm.prank(bob);
        vm.expectRevert(PitForecasts.NotCommitter.selector);
        fc.resolve(id, keccak256("o"));
    }
}

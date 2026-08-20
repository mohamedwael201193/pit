// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitPolicy} from "../src/PitPolicy.sol";

contract PitPolicyTest is Test {
    PitPolicy internal pol;
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);
    bytes32 internal ws = keccak256("workspace-a");

    function setUp() public {
        pol = new PitPolicy();
    }

    function test_pinAndIsolate() public {
        vm.prank(alice);
        pol.pin(ws, keccak256("policy-v1"));
        assertEq(pol.walletOf(ws), alice);
        assertEq(pol.versionOf(ws), 1);

        vm.prank(bob);
        vm.expectRevert(PitPolicy.NotWorkspace.selector);
        pol.pin(ws, keccak256("hostile"));

        vm.prank(alice);
        pol.pin(ws, keccak256("policy-v2"));
        assertEq(pol.versionOf(ws), 2);
        assertEq(pol.policyHashOf(ws), keccak256("policy-v2"));
    }

    function test_rejectEmpty() public {
        vm.prank(alice);
        vm.expectRevert(PitPolicy.EmptyHash.selector);
        pol.pin(bytes32(0), keccak256("x"));
    }

    function testFuzz_onlyOwnerRepins(address other) public {
        vm.assume(other != address(0) && other != alice);
        vm.prank(alice);
        pol.pin(ws, keccak256("p"));
        vm.prank(other);
        vm.expectRevert(PitPolicy.NotWorkspace.selector);
        pol.pin(ws, keccak256("q"));
    }
}

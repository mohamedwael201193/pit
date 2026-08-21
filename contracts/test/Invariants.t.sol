// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitPolicy} from "../src/PitPolicy.sol";
import {PitReceipts} from "../src/PitReceipts.sol";
import {PitForecasts} from "../src/PitForecasts.sol";

contract InvariantsTest is Test {
    PitPolicy internal pol;
    PitReceipts internal rec;
    PitForecasts internal fc;

    function setUp() public {
        pol = new PitPolicy();
        rec = new PitReceipts();
        fc = new PitForecasts();
    }

    function testFuzz_receiptCannotRefile(bytes32 preview, bytes32 root, bytes32 other) public {
        vm.assume(preview != bytes32(0) && root != bytes32(0) && other != bytes32(0));
        rec.file(preview, root);
        vm.expectRevert(PitReceipts.AlreadyFiled.selector);
        rec.file(preview, other);
    }

    function testFuzz_policyOwnerIsolation(address a, address b, bytes32 ws, bytes32 h1, bytes32 h2) public {
        vm.assume(a != address(0) && b != address(0) && a != b);
        vm.assume(ws != bytes32(0) && h1 != bytes32(0) && h2 != bytes32(0));
        vm.prank(a);
        pol.pin(ws, h1);
        vm.prank(b);
        vm.expectRevert(PitPolicy.NotWorkspace.selector);
        pol.pin(ws, h2);
    }

    function testFuzz_forecastResolverIsCommitter(address a, address b, bytes32 id, bytes32 commitHash, bytes32 outcome) public {
        vm.assume(a != address(0) && b != address(0) && a != b);
        vm.assume(id != bytes32(0) && commitHash != bytes32(0) && outcome != bytes32(0));
        vm.prank(a);
        fc.commit(id, commitHash);
        vm.prank(b);
        vm.expectRevert(PitForecasts.NotCommitter.selector);
        fc.resolve(id, outcome);
    }
}

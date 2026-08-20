// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitReceipts} from "../src/PitReceipts.sol";
import {PitForecasts} from "../src/PitForecasts.sol";

contract PitReceiptsForecastsTest is Test {
    PitReceipts internal rec;
    PitForecasts internal fc;
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    function setUp() public {
        rec = new PitReceipts();
        fc = new PitForecasts();
    }

    function test_fileOnce() public {
        bytes32 h = keccak256("preview");
        bytes32 root = keccak256("root");
        vm.prank(alice);
        rec.file(h, root);
        vm.prank(alice);
        vm.expectRevert(PitReceipts.AlreadyFiled.selector);
        rec.file(h, root);
        assertEq(rec.filer(h), alice);
    }

    function test_forecastIsolation() public {
        bytes32 id = keccak256("f1");
        vm.prank(alice);
        fc.commit(id, keccak256("c"));
        vm.prank(bob);
        vm.expectRevert(PitForecasts.NotCommitter.selector);
        fc.commit(id, keccak256("x"));
        vm.prank(alice);
        fc.resolve(id, keccak256("o"));
        vm.prank(alice);
        vm.expectRevert(PitForecasts.AlreadyResolved.selector);
        fc.resolve(id, keccak256("o2"));
    }
}

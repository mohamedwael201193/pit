// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {PitPolicy} from "../src/PitPolicy.sol";
import {PitReceipts} from "../src/PitReceipts.sol";
import {IERC7857} from "../src/interfaces/IERC7857.sol";
import {PitDeskID} from "../src/PitDeskID.sol";

contract SecurityTest is Test {
    PitPolicy internal pol;
    PitReceipts internal rec;
    PitDeskID internal desk;

    function setUp() public {
        pol = new PitPolicy();
        rec = new PitReceipts();
        desk = new PitDeskID();
    }

    function test_emptyInputsRevert() public {
        vm.expectRevert(PitPolicy.EmptyHash.selector);
        pol.pin(bytes32(0), keccak256("x"));
        vm.expectRevert(PitReceipts.Empty.selector);
        rec.file(bytes32(0), keccak256("r"));
    }

    function test_transferFromDisabled() public {
        vm.expectRevert(IERC7857.ERC7857UseITransferFrom.selector);
        desk.transferFrom(address(this), address(0xB0B), 1);
    }
}

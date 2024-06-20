
// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./node_modules/@iota/iscmagic/ISC.sol";

contract GetAllowance {
    event AllowanceFrom(ISCAssets assets);
    event AllowanceTo(ISCAssets assets);
    event Allowance(ISCAssets assets);

    function getAllowanceFrom(address _address) public {
        ISCAssets memory assets = ISC.sandbox.getAllowanceFrom(_address);
        emit AllowanceFrom(assets);
    }
}
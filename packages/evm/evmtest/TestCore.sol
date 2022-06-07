// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import "@isccontract/ISC.sol";

contract TestCore {

    function fibonacci(uint32 n) public view returns(uint32) {
        if (n == 0 || n == 1) {
            return n;
        }

        return fibonacci(n-1) + fibonacci(n-2);
    }
}

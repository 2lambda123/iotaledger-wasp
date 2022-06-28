// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import "@isccontract/ISC.sol";

contract TestCore {

    function fibonacci(uint64 n) public pure returns(uint64) {
        if (n == 0 || n == 1) {
            return n;
        }

        return fibonacci(n - 1) + fibonacci(n - 2);
    }

    function fibonacciIndirect(uint64 n) public view returns(uint64) {
        if (n == 0 || n == 1) {
            return n;
        }
        uint64 n1;
        uint64 n2;
        n1 = TestCore(this).fibonacciIndirect(n - 1);
        n2 = TestCore(this).fibonacciIndirect(n - 2);
        return n1 + n2;
    }

    event FibonacciResultEvent(uint64 n);
    function fibonacciLoop(uint64 n) public returns(uint64) {
        if (n == 0) {
            return 0;
        }
        uint64 a = 1;
        uint64 b = 1;
        for (uint64 i = 2;i < n;i++) {
            uint64 c = a + b;
            a = b;
            b = c;
        }
        emit FibonacciResultEvent(b);
        return b;
    }

    function loop(uint64 n) public pure {
        for (uint64 i = 0;i < n;i++) {
            // do nothing just burn gas
        }
    }
}

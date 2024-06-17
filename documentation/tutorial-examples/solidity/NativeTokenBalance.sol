// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./node_modules/@iota/iscmagic/ISC.sol";

contract NativeTokenBalance {
    event NativeTokenBalance(uint balance);

    function getNativeTokenBalance(bytes memory nativeTokenID) public {
        ISCAgentID memory agentID = ISC.sandbox.getSenderAccount();

        ISCDict memory params = ISCDict(new ISCDictItem[](2));
        params.items[0] = ISCDictItem("a", agentID.data);
        params.items[1] = ISCDictItem("N", nativeTokenID);

        ISCDict memory result = ISC.sandbox.callView(
            ISC.util.hn("accounts"),
            ISC.util.hn("balanceNativeToken"),
            params
        );

        emit NativeTokenBalance(bytesToUint(result.items[0].value));
    }

    function bytesToUint(bytes memory b) internal pure virtual returns (uint256) {
        require(b.length <= 32, "Bytes length exceeds 32.");
        return abi.decode(abi.encodePacked(new bytes(32 - b.length), b), (uint256));
    }


    function simpleFunction() public pure returns (uint) {
        return 42;
    }
}
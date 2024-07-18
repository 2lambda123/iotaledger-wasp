// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iscmagic/ISC.sol";

contract MyToken {
    event MintedToken(uint32 foundrySN);

    constructor(
        string memory _tokenName,
        string memory _tokenSymbol,
        uint8 _tokenDecimals,
        uint256 _maximumSupply,
        uint64 _storageDeposit
    ) payable {
        require(
            msg.value == _storageDeposit * (10 ** 12),
            "Please send exact funds to pay for storage deposit"
        );
        ISCAssets memory allowance;
        allowance.baseTokens = _storageDeposit;

        NativeTokenScheme memory nativeTokenScheme = NativeTokenScheme({
            mintedTokens: 0,
            meltedTokens: 0,
            maximumSupply: _maximumSupply
        });

        uint32 foundrySN = ISC.accounts.createNativeTokenFoundry(
            _tokenName,
            _tokenSymbol,
            _tokenDecimals,
            nativeTokenScheme,
            allowance
        );
        emit MintedToken(foundrySN);
    }
}

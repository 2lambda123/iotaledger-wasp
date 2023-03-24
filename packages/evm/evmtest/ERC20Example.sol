// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.5;

import "@iscmagic/ISC.sol";
import "@iscmagic/ERC20NativeTokens.sol";

contract ERC20Example {
    uint32 private foundrySN;

    function mint(uint256 amount, uint64 storageDeposit) public {
      ISC.accounts.mintNativeTokens(foundrySN, amount, getAllowance(storageDeposit));
    }

    function createFoundry(uint256 maxSupply, uint64 storageDeposit) public {
      NativeTokenScheme memory tokenScheme;
      tokenScheme.maximumSupply = maxSupply;
      foundrySN = ISC.accounts.foundryCreateNew(tokenScheme, getAllowance(storageDeposit));
    }

    function registerToken(string memory name, string memory symbol, uint8 decimals, uint64 storageDeposit) public {
      ISC.sandbox.registerERC20NativeToken(foundrySN, name, symbol, decimals, getAllowance(storageDeposit));
    }

    function getAllowance(uint64 amount) private pure returns (ISCAssets memory) {
      ISCAssets memory assets;
      assets.baseTokens = amount;
      return assets;
    }
}

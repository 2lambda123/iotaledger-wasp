// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iscmagic/ISC.sol";

contract L1Assets {
  function withdraw(L1Address memory to) public {
    ISCAssets memory allowance = ISC.sandbox.getAllowanceFrom(msg.sender);
    ISC.sandbox.takeAllowedFunds(msg.sender, allowance);

    ISCSendMetadata memory metadata;
    ISCSendOptions memory options;
    ISC.sandbox.send(to, allowance, false, metadata, options);
  }
}

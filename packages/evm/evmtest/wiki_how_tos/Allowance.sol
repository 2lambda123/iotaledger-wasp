// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iscmagic/ISC.sol";

contract allowance {
    function takeAllowedFunds(
        address _address,
        bytes32 _allowanceNFTID
    ) public {
        NFTID[] memory nftIDs = new NFTID[](1);
        nftIDs[0] = NFTID.wrap(_allowanceNFTID);
        ISCAssets memory assets;
        assets.nfts = nftIDs;
        ISC.sandbox.takeAllowedFunds(_address, assets);
    }
}

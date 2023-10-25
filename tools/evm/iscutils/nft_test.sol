// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.5;

import "@iota/iscmagic/ISCTypes.sol";
import "./nft.sol";

contract NFTTest {
    using nft for *;

    // Testing the MintNFT function
    function MintNFT(
        bytes memory immutableData,
        ISCAgentID memory agentID,
        NFTID nftID,
        bool withdrawal
    ) public {
        nft.MintNFT(immutableData, agentID, nftID, withdrawal);
    }
}
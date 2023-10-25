// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: MIT
pragma solidity >=0.8.5;

import "@iota/iscmagic/ISC.sol";

library nft {
    function MintNFT(
        bytes memory immutableData,
        ISCAgentID memory agentID,
        NFTID nftID,
        bool withdrawal
    ) public {
        ISCDict memory params = ISCDict(new ISCDictItem[](4));
        // ([]byte): ImmutableData for the NFT.
        params.items[0] = ISCDictItem("I", immutableData);
        // (AgentID): AgentID that will be the owner of the NFT.
        params.items[1] = ISCDictItem("a", agentID.data);
        // (optional NFTID - default empty): collectionID (NFTID) for the NFT.
        params.items[2] = ISCDictItem("C", abi.encodePacked(nftID));
        // (optional bool - default false): whether to withdrawal the NFT in the minting step
        // (can only be true when the targetAgentID is a L1 address).
        params.items[3] = ISCDictItem("W", boolToBytes(withdrawal));

        ISCAssets memory allowance;
        ISC.sandbox.call(
            ISC.util.hn("accounts"),
            ISC.util.hn("mintNFT"),
            params,
            allowance
        );
    }

    function boolToBytes(bool b) internal pure returns (bytes memory) {
        return abi.encodePacked(b ? uint8(1) : uint8(0));
    }
}
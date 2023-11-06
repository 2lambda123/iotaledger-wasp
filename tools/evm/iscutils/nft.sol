// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: MIT
pragma solidity >=0.8.5;

import "@iota/iscmagic/ISC.sol";

library NFT {
    /// @notice Mints an NFT with ImmutableData that will be owned by the AgentID.
    /// It's possible to mint as part of a collection (the caller must be the owner
    /// of the collection NFT to mint new NFTs as part of said collection).
    /// The mint can be done directly to any L1 address (it is not necessary for
    /// the target to have an account on the chain).
    /// @param immutableData Immutable data for the NFT
    /// @param agentID The AgentID that will be the owner of the NFT
    /// @param nftID CollectionID (NFTID) for the NFT
    /// @param withdrawal Whether to withdrawal the NFT in the minting step
    /// (can only be true when the targetAgentID is a L1 address)
    /// @return MintID the internal ID of the NFT at the time of minting that can be
    /// used by users/contracts to obtain the resulting NFTID on the next block
    function MintNFT(
        bytes memory immutableData,
        ISCAgentID memory agentID,
        NFTID nftID,
        bool withdrawal,
        uint64 allowance
    ) public returns (bytes6 MintID) {
        ISCDict memory params = ISCDict(new ISCDictItem[](4));
        params.items[0] = ISCDictItem("I", immutableData);
        params.items[1] = ISCDictItem("a", agentID.data);
        params.items[2] = ISCDictItem("C", abi.encodePacked(nftID));
        params.items[3] = ISCDictItem("W", boolToBytes(withdrawal));

        ISCAssets memory assets= makeAllowanceBaseTokens(allowance);

        ISCDict memory returnedDict = ISC.sandbox.call(
            ISC.util.hn("accounts"),
            ISC.util.hn("mintNFT"),
            params,
            assets
        );

        ISCDictItem memory item = returnedDict.items[0];
        return bytes6(item.value);
    }
    function makeAllowanceBaseTokens(uint64 amount) internal pure returns (ISCAssets memory) {
        ISCAssets memory assets;
        assets.baseTokens = amount;
        return assets;
    }
    function boolToBytes(bool b) internal pure returns (bytes memory) {
        return abi.encodePacked(b ? uint8(1) : uint8(0));
    }
}
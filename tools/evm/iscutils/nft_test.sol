// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.5;

import "@iota/iscmagic/ISC.sol";
import "./nft.sol";

contract NFTTest {
    using NFT for *;
    event Minted(bytes6 mintID);

    bytes public constant immutableData = "0xabc123";
    ISCAgentID public agentID = ISCAgentID({data: abi.encodePacked(msg.sender)});
    NFTID public constant nftID = NFTID.wrap(bytes32(0));
    bool public constant withdrawal = false;

    // Testing the MintNFT function
    function MintTestNFT() public {
        require(true == false,"kaboom");
        bytes6 mintID = NFT.MintNFT(immutableData, agentID, nftID, withdrawal);

        require(bytes32(mintID) != bytes32(0),"empty MintID returned");
       // mintID = NFT.MintNFT(immutableData, agentID, nftID, withdrawal);
       // require(mintID.length == 0,"empty MintID returned");
        emit Minted(mintID);
    }

    function MintNFTWithReturnedDict() public returns (ISCDict memory) {

        ISCDict memory params = ISCDict(new ISCDictItem[](4));
        params.items[0] = ISCDictItem("I", immutableData);
        params.items[1] = ISCDictItem("a", agentID.data);
        params.items[2] = ISCDictItem("C", abi.encodePacked(nftID));
        params.items[3] = ISCDictItem("W", boolToBytes(withdrawal));

        ISCAssets memory allowance = makeAllowanceBaseTokens(1000000);

       return ISC.sandbox.call(
            ISC.util.hn("accounts"),
            ISC.util.hn("mintNFT"),
            params,
            allowance
        );


    }
    function makeAllowanceBaseTokens(uint64 amount) private pure returns (ISCAssets memory) {
        ISCAssets memory assets;
        assets.baseTokens = amount;
        return assets;
    }
    function boolToBytes(bool b) internal pure returns (bytes memory) {
        return abi.encodePacked(b ? uint8(1) : uint8(0));
    }
    /*
    function MintCustomNFT(
        bytes memory immutableData,
        ISCAgentID memory agentID,
        NFTID nftID,
        bool withdrawal
    ) public {
        bytes memory mintID = NFT.MintNFT(immutableData, agentID, nftID, withdrawal);
        emit Minted(mintID);
    }
    */
}
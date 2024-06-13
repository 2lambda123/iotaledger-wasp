# Improvement of ERC721NFTs contract

## Add name and symbol

### Reason

For ERC721 NFT collection, the [name](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC721/ERC721.sol#L74) and [symbol](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC721/ERC721.sol#L81) can be set for better display on the explorer.
Currently, due to no name/symbol set, the block explorer displays a default value `Unnamed token`.

### Modified functions

```
    function name() external view virtual returns (string memory) {
        return "CollectionL1";
    }

    function symbol() external view virtual returns (string memory) {
        return "CollectionL1";
    }
```

## Handle system error for non-existent NFT token

### Reason

The call of `getNFTData()` (e.g. in the function `ownerOf()` or `tokenURI()`) will throw the error "invalid memory address or nil pointer dereference" if the NFT tokenId does not exist (not yet minted or already burnt).

Thus the new `function _requireOwned()` (according to OpenZeppelin [ERC721.sol](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC721/ERC721.sol#L449)) is added to revert with user-friendly error message in such case. This function is called by other functions like `ownerOf()` or `tokenURI()`.

### Added function

```
    function _requireOwned(uint256 tokenId) internal view returns (address)
```

## Remove the dummy function \_isManagedByThisContract()

### Reason

The function `_isManagedByThisContract()` is just dummy and thus it should be removed.

### Removed function

```
function _isManagedByThisContract(ISCNFT memory)
```

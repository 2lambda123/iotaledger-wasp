## ISCP accounts. Controlling token balances

ISCP provides secure and trustless transfers of digitized assets:
- between smart contracts on the same or on different chains
- between smart contracts and addresses on the UTXO Ledger.

On the UTXO Ledger, just like in any DLT, we have **trustless** and
**atomic** transfers of assets between addresses on the ledger. The tokens
contained in the address can be moved to another address by providing a valid
signature with the private key which controls the source address.

In ISCP, the smart contracts which reside on chains are also owners of their
tokens. Each smart contract can receive tokens transferred to it and can move
tokens controlled by it to any other owner, be it a smart contract or an
ordinary address on the UTXO Ledger.

So, there are 2 types of entities which control tokens:

* Addresses on the UTXO Ledger
* Smart contracts on ISCP chains

There are 3 different types of trustless token transfers possible between those 
entities. Each type involves different mechanism of transfer:

* between address and smart contract.
* between smart contracts on the same chain
* between smart contracts on different chains

To make the system homogenous, we introduce the following two concepts:

* _Agent ID_, which represents the ID of the token owner abstracted away from 
  the type of the owning entity
* _On-chain account_ which represents the unit of ownership on the chain

### Smart contract ID

Unlike in blockchain systems like Ethereum, we cannot simply represent the smart
contract by a blockchain address: ISCP can have more than one blockchain. Each
chain in ISCP is identified by its _chain address_. A chain can contain many
smart contracts on it. So, in ISCP each contract is identified by concatenation
of a chain identifier, the ChainID, and the _hname_ of the smart contract:
`chainID || hname`. In human-readable form the smart _contract ID_ looks
like this:

```
A/RJNmyghMeM4Yr3UtBnric8mmBBwWdt9yVifetdpCQj7J::cebf5908
```

The part before `::` is the chain ID (the chain address), the part after `::` 
is the _hname_ of the smart contract, which is the contract identifier on the 
chain interpreted as a hexadecimal number. The `A/` prefix indicates that 
this is an agent ID.

### Agent ID

The agent ID is an identifier which generalizes and represents one of two 
possible identifiers: either an address on the UTXO Ledger, or a smart
_contract ID_.

It is easily possible to determine which one is represented by the particular
agent ID: when it is an address the _hname_ part will be all zeroes.

Address is a data
type [defined by Goshimmer](https://github.com/iotaledger/goshimmer/blob/master/packages/ledgerstate/address.go#L43).

The `AgentID` type
is [defined by Wasp](https://github.com/iotaledger/wasp/blob/master/packages/coretypes/agentid.go#L25):
The _agent ID_ value contains information which one of two types it represents:
address or contract ID.

### On-chain accounts

Each chain can contain any number of accounts. Each account contains colored
tokens: a collection of `color: balance` pairs.

Each account on the chain is controlled by some `agent ID`. That means that 
tokens contained in the account can only be moved by the entity that is 
associated with the agent ID:

* If the _agent ID_ represents an address on the UTXO Ledger, the tokens can
  only be moved by a request, sent (and signed) by that address.
* If the _agent ID_ represents a smart contract, the tokens on the account can
  only be moved by that smart contract: independent of whether the smart 
  contract resides on the same chain or on another chain.

![](accounts.png)

The picture illustrates an example situation. There are two chains deployed,
with respective IDs
`Pmc7iH8b..` and `Klm314noP8..`. The pink chain `Pmc7iH8b..` has two smart
contracts on it (`3037` and `2225`) and the blue chain `Klm314noP8..` has one
smart contract (`7003`).

The UTXO Ledger has 1 address `P6ZxYXA2..`. The address `P6ZxYXA2..` controls
1337 iotas and 42 red tokens on the UTXO Ledger. The same address also controls
42 iotas on the pink chain and 8 green tokens on the blue chain. So, the owner
of the private key behind the address controls 3 different accounts: 1 on the L1
UTXO Ledger, and 2 accounts on 2 different chains on L2.

At the same time, smart contract `7003` on the blue chain has 5 iotas on its 
native chain and controls 11 iotas on the pink chain.

Note that “control over account” means that the entity which has the private key
can move funds from it. For an ordinary address it means its private key. For a
smart contract it means the private keys of the committee which runs the chain
(the smart contract program can only be executed with those private keys).

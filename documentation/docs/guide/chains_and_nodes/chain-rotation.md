# Chain Rotation

## Get public key of wallet address

The “Chain Owner” is the only one who can perform administrative tasks. Any iota/SMR L1 pub key can be used, and whoever owns that key, basically owns the chain. In this example we'll use the public key of the wallet address. In a consortium it would probably be wise to use a multisig setup.

Get wallet public key with the command below

```bash
wasp-cli address
```

Copy the public key and pass that to the `--gov-contraller` flag on the `wasp-cli chain deploy` command.

## Start chain

Before starting a chain, all nodes need to trust each other. See [Trust setup](./setting-up-a-chain.md#trust-setup) for instructions on chain peering. Deploy a new chain with the command below

```bash
wasp-cli chain deploy --description "Test Chain" --gov-controller ${wallet_public_key} --committee 0,1 --evm-chainid 1076
```

## Add access nodes

Depending on your intentions for rotating the chain, you may need to add some access nodes. If all you want to do is remove existing nodes from the chain then skip this step. To learn more about adding access nodes see [chain management](./chain-management.md#changing-access-nodes). Once you have configuted at least one access node, add their public keys with the command below.

```bash
wasp-cli chain change-access-nodes accept ${node_pub_key}
```

Run that command for each public key.

## Run DKG

Update your `wasp-cli` config and run the dkg command with the indices for the new committee members

```bash
wasp-cli chain rundkg --committee=0,2 --quorum=2
```

Copy the new public key

## Rotate chain

```bash
wasp-cli chain rotate ${new_chain_pub_key}
```

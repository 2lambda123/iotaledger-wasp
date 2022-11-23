# Chain Rotation

## Get public key of wallet address

```bash
wasp-cli chain address
```

Copy the public key and pass that to the `--gov-contraller` flag on the `wasp-cli chain deploy` command.

## Start chain

Before starting a chain, all nodes need to trust each other. See [Trust setup](./setting-up-a-chain.md#trust-setup) for instructions on chain peering. Deploy a new chain with the command below

```bash
wasp-cli chain deploy --description "Test Chain" --gov-controller ${wallet_public_key} --committee 0,1 --evm-chainid 1076
```

## Add access nodes

Before rotating the chain, you need to some access nodes. To learn more about access nodes see [chain management](./chain-management.md#changing-access-nodes).

```bash
wasp-cli chain change-access-nodes accept ${node_pub_key}
```

## Run DKG

Run dkg commannd with new committee indices.

```bash
wasp-cli chain rundkg --committee=0,2 --quorum=2
```

Copy the new public key

## Rotate chain

```bash
wasp-cli chain rotate ${new_chain_pub_key}
```

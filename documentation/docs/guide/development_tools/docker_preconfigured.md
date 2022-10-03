---
description: How to run the preconfigured Docker setup.
image: /img/logo/WASP_logo_dark.png
keywords:

- smart contracts
- validator node
- Docker
- image
- build
- configure
- arguments
- Hornet
- how to

---

# Preconfigured Development Docker Setup

This page describes how you can use the preconfigured developer Docker setup.

:::note
This setup is intended for **local** development only (you will have your own private network/tangle).
:::

private tangle, ready to run out of the box.

## Requirements

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)
* [Git](https://git-scm.com/)

## Quick Start Guide


1. Checkout the project:

    ```shell
    git clone https://github.com/iotaledger/wasp.git
    ```

2. Mode into the project directory and check out the `develop` branch:

    ```shell
   cd wasp
   git checkout develop
    ```

3. Move into the `tools/devnet` folder:

    ```shell
     cd tools/devnet
    ```

4. Run the following command to start the setup.

    ```shell
    docker-compose up
    ```

This command will initialize Hornet and create a fresh image of the checked-out Wasp code.

If you modify the branch, docker-compose will include your modifications in the Wasp image.

:::note Default Ports

All Wasp ports will bind to 127.0.0.1 by default.
If you want to expose the ports to the outside world, run `HOST=0.0.0.0 docker-compose up`.

:::


## Wasp-CLI Configuration

All of the `wasp-cli` ports are locally available. You should use this `wasp-cli.json` configuration:

```json
{
  "l1": {
    "apiaddress": "http://localhost:14265",
    "faucetaddress": "http://localhost:8091"
  },
  "wasp": {
    "0": {
      "api": "127.0.0.1:9090",
      "nanomsg": "127.0.0.1:5550",
      "peering": "127.0.0.1:4000"
    }
  }
}
```

Run `wasp-cli init` to generate a seed, and you are ready to go.

See [Configuring wasp-cli](../chains_and_nodes/wasp-cli) for further information.

## Usage

Wasp is configured to allow any connection coming from `wasp-cli`. This is fine for development purposes, but please
ensure you don’t run it on a publicly available server or create matching firewall filter rules if you do so.

Other than that, everything should simply work as expected. Faucet requests will be handled. You will be able
to deploy and run smart contracts. All useful [ports](#reachable-ports) such as are available to the local machine.

### Start

To start the setup, run:

```bash
docker-compose up
```

During the startup you might see a few failed restarts of Wasp with the message:
`panic: error getting node event client: mqtt plugin not available on the current node`

This is normal, as Wasp starts faster than Hornet. Wasp retries the connection until it succeeds.

### Stop

You can shut down the setup by running the following command in a new terminal:

```shell
docker-compose down
```

You can also press `Ctrl-C` to shut down the setup, but **don't press it twice to force it**. Otherwise, you can corrupt
the Hornet database.

### Reset

To shut down the nodes and to remove all databases run:

```shell
docker-compose down --volumes
```

### Recreation

If you made changes to the Wasp code and want to use it inside the setup, you need to recreate the Wasp image by
running:

```bash
docker-compose build
```

## Reachable Ports

The nodes will then be reachable under these ports:

### Wasp

#### API

[http://localhost:9090]( http://localhost:9090)

#### Dashboard

[http://localhost:7000](http://localhost:7000) (**Username**: wasp  **Password**: wasp)

#### Nanomsg

[tcp://localhost:5550](tcp://localhost:5550)

### Hornet

#### API

[http://localhost:14265](http://localhost:14265)

#### Faucet

[http://localhost:8091](http://localhost:8091)

#### Dashboard

[http://localhost:8081](http://localhost:8081) (**Username**: admin **Password**: admin)

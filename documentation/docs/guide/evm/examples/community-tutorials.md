---
title: Community Tutorials for the IOTA EVM Testnet
keywords:
- ISCP
- IOTA
- Smart Contracts
- EVM
- Metamask
- Solidity
- Tutorial
- Step by Step Guide
description: Easy to follow Step by Step Tutorials created by community member ZAIRED. No coding skills required!
image: /img/logo/WASP_logo_dark.png
---

# Community Tutorials for the IOTA EVM Testnet
## Setting up Metamask for IOTA smart contracts
## 1. Install metamask

Available at [https://metamask.io](https://metamask.io/) and on the chrome webstore

[https://metamask.io](https://metamask.io)

## 2. Create an account

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/5ef05765-a772-4ad9-9222-7d290d1018e8/Untitled.png)

Create a wallet unless you already have one.

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/fe61f801-bceb-4a7b-add5-7ce7164e9ce1/Untitled.png)

Create a strong password

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/d1dc6808-1cff-476e-a7a9-8d8b20cb9de8/Untitled.png)

### Save your keyphrase somewhere secure and never share it with anyone

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/efca6dd5-fb41-4542-9bdd-53f7e5d8cf7d/Untitled.png)

## 3. Add the IOTA custom testnet

Time to setup the custom RPC network

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/b2610275-3dfd-4f6c-9085-60018377081d/Untitled.png)

Next, add the following values

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/b643d406-acfd-450f-80a3-90c777ed3ae8/Untitled.png)

[Values](https://www.notion.so/de5dd24f7d5645b2a5f8d24b56d0097c)

## Now hit Save!

Congratulations, you are done, you can now use IOTA smart contracts.

If you haven't already, go setup metamask [using this guide](https://www.notion.so/Setting-up-Metamask-for-IOTA-smart-contracts-fa52b6d49f3446e5947f8f37606c82cc)!

This guide is going to explain how to create an ERC20 token, not ERC721 (aka NFT)

### Note: Metamask is required for this guide. See installation guide above.

## 1. Create the smart contract

Go to [https://wizard.openzeppelin.com](https://wizard.openzeppelin.com/)

**Create a contract, with (1.) a custom name and symbol, (2.) a premint amout (aka the amount you start with) and finally (3.) select if you want it to be mintable if you want to add more later and if you want it burnable (if you want to be able to destroy them)**

![ISCP1.png](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/f03c8652-e11e-40ee-946b-9bb4bfd292ff/ISCP1.png)

## 2. Export to Remix IDE

![ISCP2.png](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/34905c2a-89b4-4d83-907d-a5e5d5536042/ISCP2.png)

## 3. Compile the code

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/c3115292-f40d-4214-89d9-ea5a90314315/Untitled.png)

Simply press the big blue button

## 4. Export your smart contract

Switch to the deploy tab

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/5fd3decd-0362-420e-b86e-1530e55bbd38/Untitled.png)

Select Injected Web3 environment

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/bf00e41b-99ca-48d0-ab37-1df6a131212c/Untitled.png)

A window should popup to connect metamask

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/bd90362f-99ff-4cf2-8bef-fe64dbaab5a4/Untitled.png)

Press **Next** and **Connect**

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/54f1d0fa-758e-42e8-a3c3-759e62cd5114/Untitled.png)

### Select your smart contract

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/e2dc938c-9bff-4164-bb1d-10e9540573c5/Untitled.png)

### Deploy and confirm with metamask

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/efb5c511-0a3d-4cef-b431-c535c9b21074/Untitled.png)

## 5. Import the token to metamask

Find and copy the token address:  

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/17fe5cae-97d6-46a1-a3f5-8ceae610fff9/Untitled.png)

Import the token:

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/f74d6c3d-c03d-4e44-bcb8-1c2cbd144c68/Untitled.png)

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/960f31ec-f346-4ce4-a1b7-82377e5c8a4a/Untitled.png)

### Congratulation, you've successfully created your own token using IOTA smart contracts!

## 6. (optionnal) Mint more tokens

Expand the view of the contract

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/191e6b62-9cbb-42a6-bbc9-155a25842f80/Untitled.png)

Do this for mint too

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/8b3dafb6-d325-4e17-b7a7-0d88b4c6b316/Untitled.png)

Now enter the address you want to send it to

(example, my address: 0x3AC5c7b1e9A2F6ecE3e9F879D22DD45797A96D51) 

For the amount,  you want to add 18 zeros to whatever you are minting:
Example you want to min 10 tokens, you would write 10000000000000000000 

If you want to copy it here is 18 zeros: 000000000000000000

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/749ecf25-473a-4124-8479-4efe8219316f/Untitled.png)

Then simply press the transact button, and confirm it with metamask.
It should appear in a couple seconds in your wallet

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/16044e34-0cbc-47c6-a598-a18605dbec97/Untitled.png)

## This wraps up this guide for now. More may be added later.

Last update: 23/10/2021




# Gas Management

This solution is in charge of distributing gas to the different LACChain Besu writer nodes, it is composed of backend components such as smart contracts. Gas distribution is automatic, whose logic is written in smart contracts. 

## Package overview

1. **audit** contains ways to log.
2. **blockchain** contains connections to Ethereum.
3. **controller** controller layer that receives all external requests and redirects requests to the service layer
4. **service** contains main logic
5. **model** contains data models of requests and responses of APIs
6. **errors** contains different errors types
7. **relayhub** contains all smart contract 
8. **rpc** contains models and ways to interact with RPC request and response
9. **docs** contains documentation about architecture and developer interaction with this 
solution

## Prerequisites

* Being a validator node in LACChain network
* Go 1.13+ installation or later
* **GOPATH** environment variable is set correctly

## Install

```
$ git clone https://github.com/lacchain/gas-management

$ cd gas-management
$ go build
```

## Run

Execute the executable file generated previously in a Validator node

```
$ ./gas-relay-signer
```

## Know More

* [In depth overview of the GAS distribution mechanism](https://github.com/LACNetNetworks/gas-management/blob/master/docs/OVERVIEW.md)
* [How to adapt you solution to the GAS distribution mechanism](https://github.com/LACNetNetworks/gas-management/blob/master/docs/How_adapt_your_Dapp.md)
* [Deploy your first ERC20 and time-stamping (notarization) smart contracts](https://github.com/LACNetNetworks/gas-management/blob/master/docs/tutorial/Deploy_SmartContract.md)
* [Deploy and interact with the LACChain ID verifiable credential registry smart contract](https://github.com/LACNetNetworks/gas-management/blob/master/docs/tutorial/VC_en.md)
* [Stress testing and performance of the network with the GAS distribution mechanism](https://github.com/LACNetNetworks/gas-management/blob/master/docs/STRESS_TESTING.md)
* [Comparison with Ethereum](https://github.com/LACNetNetworks/gas-management/blob/master/docs/COMPARISON_WITH_ETHEREUM.md)
* [FAQ](https://github.com/LACNet-Networks/gas-management/blob/master/docs/FAQ.md)


## Concurrency nonce tests

There are two specific unit tests for handle with nonce concurrecy into relaySignerService_test.go:

TestNonceRaceSendMetatransactionDifferentSenders(t *testing.T)
TestNonceCollisionSendMetatransactionSameSender(t *testing.T)

Run this command for check the tests:

```
go test ./service -v -run ^TestNonceRace
go test ./service -race -v -run ^TestNonceRace

go test ./service -v -run ^TestNonceCollision
go test ./service -race -v -run ^TestNonceCollision

```

## Copyright 2022 LACNet

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

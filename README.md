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
$ make build      # compila con la versión inyectada desde el tag git (ver "Versión")
```

> `go build` a secas también funciona, pero deja la versión en `dev`. Usa `make build`
> (o los `-ldflags`) para que el binario reporte la versión real.

## Run

Execute the executable file generated previously in a Validator node

```
$ ./gas-relay-signer
```

## Versión

El binario reporta su versión:

```
$ ./gas-relay-signer --version
gas-relay-signer v1.1.0 (commit 5b7a3a7, built 2026-07-01T22:48:36Z, go1.23.0)
```

La versión es el **tag git** (`git describe --tags`), inyectado en compilación vía
`-ldflags "-X main.version=... -X main.commit=... -X main.date=..."` (lo hace `make build`).
Compilado en el tag `v1.1.0` reporta `v1.1.0`; en `develop` sin tag, algo como
`v1.0.1-9-g5b7a3a7`. Sin `ldflags` reporta `dev`.

### Publicar un release (manual)

1. Mergear `develop` → `master` (PR) y situarse en `master` actualizado.
2. Crear el tag anotado y empujarlo:
   ```
   git tag -a v1.1.0 -m "gas-relay-signer v1.1.0"
   git push origin v1.1.0
   ```
3. Compilar el artefacto con la versión inyectada y publicar el release:
   ```
   make build VERSION=v1.1.0
   gh release create v1.1.0 gas-relay-signer --title "v1.1.0" --notes "..."
   ```

## Know More

* [In depth overview of the GAS distribution mechanism](https://github.com/LACNetNetworks/gas-management/blob/master/docs/OVERVIEW.md)
* [How to adapt you solution to the GAS distribution mechanism](https://github.com/LACNetNetworks/gas-management/blob/master/docs/How_adapt_your_Dapp.md)
* [Deploy your first ERC20 and time-stamping (notarization) smart contracts](https://github.com/LACNetNetworks/gas-management/blob/master/docs/tutorial/Deploy_SmartContract.md)
* [Deploy and interact with the LACChain ID verifiable credential registry smart contract](https://github.com/LACNetNetworks/gas-management/blob/master/docs/tutorial/VC_en.md)
* [Stress testing and performance of the network with the GAS distribution mechanism](https://github.com/LACNetNetworks/gas-management/blob/master/docs/STRESS_TESTING.md)
* [Comparison with Ethereum](https://github.com/LACNetNetworks/gas-management/blob/master/docs/COMPARISON_WITH_ETHEREUM.md)
* [FAQ](https://github.com/LACNet-Networks/gas-management/blob/master/docs/FAQ.md)
* [Reporte del fallo de la llamada interna (status=1 → fallida)](docs/RECEIPT-FALLO-INTERNO.md) — cómo el RelaySigner reescribe el receipt a `status=0`+`revertReason` y expone `relay_getMetaTxResult` (rama `develop`).

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

# MVM Token List

MVM token list fetched from api.mvm.dev.

| file | info |
-------------------------------------------|------------|
| [mvm-tokenlist.json](mvm-tokenlist.json) | all tokens |
| [pure-tokenlist.json](pure-tokenlist.json) | list without 4swap LP tokens, Pando Rings rTokens, and each asset include Mixin AssetID and Mixin ChainID.|
| [mvm-chainlist.json](mvm-chainlist.json)| all of the chain assets |
| [evm-asset-chain-map.json](evm-asset-chain-map.json)| map of asset id and chain id of evm-compatible chains |
| [evm-chain-asset-map.json](evm-chain-asset-map.json)| map of chain id and asset id of evm-compatible chains |
| [asset-symbol-key.json](asset-symbol-key.json)| map of asset symbol to asset key|


### Generate the latest token list

1. clone this repo

2. $ go mod tidy

3. $ go run main.go

# MVM Token List

MVM token list fetched from api.mvm.dev.

| file | info |
-------------------------------------------|------------|
| [mvm-tokenlist.json](mvm-tokenlist.json) | all tokens |
| [pure-tokenlist.json](pure-tokenlist.json) | list without 4swap LP tokens, Pando Rings rTokens, and each asset include Mixin AssetID and Mixin ChainID.|
| [mvm-chainlist.json](mvm-chainlist.json)| all of the chain assets |





### Generate the latest token list

1. clone this repo

2. $ go mod tidy

3. $ go run main.go

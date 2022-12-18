# MVM Token List

MVM token list in [defillama format](https://defillama-datasets.s3.eu-central-1.amazonaws.com/tokenlist/ethereum.json) fetched from api.mvm.dev.

See [mvm-tokenlist.json](mvm-tokenlist.json) for the list of all tokens.

See [MVG-tokenlist.json](MVG-tokenlist.json) for the list without 4swap LP tokens, Pando Rings rTokens, and each asset include Mixin AssetID and Mixin ChainID.

[Uniswap token list format](https://github.com/Uniswap/token-lists) is WIP.


### Generate the latest token list

1. clone this repo

2. $ go mod tidy

3. $ go run main.go

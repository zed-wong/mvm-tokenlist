package main

import (
	"log"
	"fmt"
	"strings"
	"context"
	"io/ioutil"

	"github.com/Jeffail/gabs"
	"github.com/go-resty/resty/v2"
	"github.com/fox-one/mixin-sdk-go"
)

const (
	ENDPOINT = "https://api.mvm.dev/asset_contract?asset="
	NULL_ADDR = "0x0000000000000000000000000000000000000000"
)

var (
	NAMES = []string{"pure-tokenlist.json", "mvm-tokenlist.json", "mvm-chainlist.json"}
	STABLE_LIST = []string{"USDT", "USDC", "pUSD", "DAI"}
	LP_LIST = []string{"LP Token"}
	RINGS_LIST = []string{"Pando Rings"}
	EVM_LIST = []string{
		"43d61dcd-e413-450d-80b8-101d5e903357", // ETH
		"b7938396-3f94-4e0a-9179-d3440718156f", // Polygon
		"1949e683-6a08-49e2-b087-d6b72398588f", // BSC
		"", // Arbitrum
		"", // Optimism
		"", // Avalance
		"", // Fantom
		"", // Gnosis
		"", // Celo
	}
)

type Result struct {
	AssetContract string `json:"asset_contract"`
}

func writeFile(filename, data string) {
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getContract(rest *resty.Client, assetID string) *resty.Response{
	resp, err := rest.R().SetResult(&Result{}).Get(ENDPOINT+assetID)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func isStable(name string) bool {
	for _, stb := range STABLE_LIST {
		if strings.Contains(name, stb) {
			return true
		}
	}
	return false
}

func isLpToken(name string) bool {
	for _, lpl := range LP_LIST {
		if strings.Contains(name, lpl) {
			return true
		}
	}
	return false
}

func isRings(name string) bool {
	for _, rsl := range RINGS_LIST {
		if strings.Contains(name, rsl) {
			return true
		}
	}
	return false
}

func isChainAsset(assetID, chainID string) bool {
	if (assetID == chainID) {
		return true
	}
	return false
}

func isEVMChain(assetID string) bool {
	for _, n := range EVM_LIST {
		if (assetID == n) {
			return true
		}
	}
	return false
}

func llamaTokenlist(name string) {
	rest := resty.New()
	ctx := context.Background()

	fmt.Println("Started collecting token info...")
	topAssets, err := mixin.ReadTopNetworkAssets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	o := gabs.New()
	for _, asset := range topAssets {
		obj := gabs.New()
		res := getContract(rest, asset.AssetID).Result().(*Result)
		if (res.AssetContract == NULL_ADDR) { continue }

		obj.Set(asset.Name, res.AssetContract, "name")
		obj.Set(asset.Symbol, res.AssetContract, "symbol")
		obj.Set(asset.IconURL, res.AssetContract, "logoURI")
		obj.Set(73927, res.AssetContract, "chainId")
		obj.Set(8, res.AssetContract, "decimals")
		fmt.Println(obj.StringIndent("", " "))
		o.Merge(obj)
	}
	writeFile(name, o.StringIndent("", " "))
	fmt.Println("Token info saved in", name)
}

func PureTokenlist(name string) {
	rest := resty.New()
	ctx := context.Background()

	fmt.Println("Started collecting token info...")
	topAssets, err := mixin.ReadTopNetworkAssets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	o := gabs.New()
	for _, asset := range topAssets {
		obj := gabs.New()
		res := getContract(rest, asset.AssetID).Result().(*Result)
		if (res.AssetContract == NULL_ADDR) { continue }

		if (isLpToken(asset.Name)) { continue }
		if (isRings(asset.Name)) { continue }
		obj.Set(res.AssetContract, res.AssetContract, "contract")
		obj.Set(isStable(asset.Symbol), res.AssetContract, "stable")
		obj.Set(asset.AssetID, res.AssetContract, "mixinAssetId")
		obj.Set(asset.ChainID, res.AssetContract, "mixinChainId")
		obj.Set(asset.Name, res.AssetContract, "name")
		obj.Set(asset.Symbol, res.AssetContract, "symbol")
		obj.Set(asset.IconURL, res.AssetContract, "logoURI")
		obj.Set(73927, res.AssetContract, "chainId")
		obj.Set(8, res.AssetContract, "decimals")
		fmt.Println(obj.StringIndent("", " "))
		o.Merge(obj)
	}
	writeFile(name, o.StringIndent("", " "))
	fmt.Println("Token info saved in", name)
}

func SymbolAddressList(name string) {
	// WIP
}

func MVMChainList(name string) {
	// Only Chain Asset
	rest := resty.New()
	ctx := context.Background()

	fmt.Println("Started collecting token info...")
	topAssets, err := mixin.ReadTopNetworkAssets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	o := gabs.New()
	for _, asset := range topAssets {
		obj := gabs.New()
		res := getContract(rest, asset.AssetID).Result().(*Result)
		if (res.AssetContract == NULL_ADDR) { continue }
		if (isLpToken(asset.Name)) { continue }
		if (isRings(asset.Name)) { continue }
		if (!isChainAsset(asset.AssetID, asset.ChainID)) { continue }

		obj.Set(res.AssetContract, res.AssetContract, "contract")
		obj.Set(isStable(asset.Symbol), res.AssetContract, "stable")
		obj.Set(asset.AssetID, res.AssetContract, "mixinAssetId")
		obj.Set(asset.ChainID, res.AssetContract, "mixinChainId")
		obj.Set(asset.Name, res.AssetContract, "name")
		obj.Set(asset.Symbol, res.AssetContract, "symbol")
		obj.Set(asset.IconURL, res.AssetContract, "logoURI")
		obj.Set(73927, res.AssetContract, "chainId")
		obj.Set(8, res.AssetContract, "decimals")
		obj.Set(false, res.AssetContract, "evm")
		if (isEVMChain(asset.AssetID)) {
			obj.Set(true, res.AssetContract, "evm")
		}
		fmt.Println(obj.StringIndent("", " "))
		o.Merge(obj)
	}
	writeFile(name, o.StringIndent("", " "))
	fmt.Println("Token info saved in", name)
}

func main() {
	PureTokenlist(NAMES[0])
	llamaTokenlist(NAMES[1])
	MVMChainList(NAMES[2])
}

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
	NAMES = []string{"MVG-tokenlist.json", "mvm-tokenlist.json"}
	STABLE_LIST = []string{"USDT", "USDC", "pUSD", "DAI"}
	LP_LIST = []string{"LP Token"}
	RINGS_LIST = []string{"Pando Rings"}
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

func MVGTokenlist(name string) {
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
		obj.Set(isStable(asset.Name), res.AssetContract, "stable")
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

func main() {
	MVGTokenlist(NAMES[0])
	llamaTokenlist(NAMES[1])
}

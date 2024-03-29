package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

const (
	ENDPOINT  = "https://api.mvm.dev/asset_contract?asset="
	NULL_ADDR = "0x0000000000000000000000000000000000000000"
	NaNa_ADDR = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

	// ETH used in the registry contract, but deprecated after the native currency change.
	DEPRECATED_ETH = "0x181251D3A501961d4Af2AF46E33C71A5D808c25B"
	// WETH used in the registry contract, but useless when it comes to MVM because the lack of methods.
	DEPRECATED_WETH = "0x6D759901Aa3104BAAE6c15EA19eaE06A84d4cC3d"
	// WETH9 deployed on MVM Mainnet to replace WETH created by the registry
	WETH9_ADDRESS = "0xBac65f64cd7Ac8a2e71800C504b1E61D8c405015"
)

var (
	NAMES       = []string{"pure-tokenlist.json", "mvm-tokenlist.json", "mvm-chainlist.json", "asset-symbol-key.json"}
	STABLE_LIST = []string{"USDT", "USDC", "pUSD", "DAI"}
	LP_LIST     = []string{"LP Token", "4swap"}
	RINGS_LIST  = []string{"Pando Rings"}
	EVM_LIST    = []string{
		"43d61dcd-e413-450d-80b8-101d5e903357", // ETH
		"b7938396-3f94-4e0a-9179-d3440718156f", // Polygon
		"1949e683-6a08-49e2-b087-d6b72398588f", // BSC
		"",                                     // Arbitrum
		"",                                     // Optimism
		"",                                     // Avalance
		"",                                     // Fantom
		"",                                     // Gnosis
		"",                                     // Celo
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

func readFile(filename string) (string, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
		return "", err
	}
	return string(body), nil
}

func getContract(rest *resty.Client, assetID string) *resty.Response {
	resp, err := rest.R().SetResult(&Result{}).Get(ENDPOINT + assetID)
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
	if assetID == chainID {
		return true
	}
	return false
}

func isEVMChain(assetID string) bool {
	for _, n := range EVM_LIST {
		if assetID == n {
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
		if res.AssetContract == NULL_ADDR {
			continue
		}
		if res.AssetContract == DEPRECATED_ETH {
			res.AssetContract = NaNa_ADDR
		}
		if res.AssetContract == DEPRECATED_WETH {
			res.AssetContract = WETH9_ADDRESS
		}

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
		if res.AssetContract == NULL_ADDR {
			continue
		}
		if res.AssetContract == DEPRECATED_ETH {
			res.AssetContract = NaNa_ADDR
		}
		if res.AssetContract == DEPRECATED_WETH {
			res.AssetContract = WETH9_ADDRESS
		}

		if isLpToken(asset.Name) {
			continue
		}
		if isRings(asset.Name) {
			continue
		}
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
		if res.AssetContract == NULL_ADDR {
			continue
		}
		if isLpToken(asset.Name) {
			continue
		}
		if isRings(asset.Name) {
			continue
		}
		if !isChainAsset(asset.AssetID, asset.ChainID) {
			continue
		}
		if res.AssetContract == DEPRECATED_ETH {
			res.AssetContract = NaNa_ADDR
		}

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
		if isEVMChain(asset.AssetID) {
			obj.Set(true, res.AssetContract, "evm")
		}
		fmt.Println(obj.StringIndent("", " "))
		o.Merge(obj)
	}
	writeFile(name, o.StringIndent("", " "))
	fmt.Println("Token info saved in", name)
}

// Only includes necessary parts for mixin usage
//
//	{
//		"asset_id":"",
//		"asset_name": "",
//		"asset_symbol": "",
//		"asset_icon": "",
//		"asset_key": "",
//		"chain_id":"",
//		"chain_name":"",
//		"chain_icon":"",
//	}
func MiniumTokenlist() {
	top, _ := readFile("mixin-top-assets.json")
	assets := gjson.Get(top, "assets").Array()
	for _, asset := range assets {
		obj := gabs.New()
		obj.Set(asset.Get("asset_id").String(), "asset_id")
		obj.Set(asset.Get("name").String(), "name")
		obj.Set(asset.Get("symbol").String(), "symbol")
		obj.Set(asset.Get("icon_url").String(), "icon")
		obj.Set(asset.Get("chain_id").String(), "chain_id")
		obj.Set(asset.Get("asset_key").String(), "asset_key")

		chainAsset := gjson.Get(top, fmt.Sprintf(`assets.#(asset_id==%s)`, asset.Get("chain_id").String()))
		obj.Set(chainAsset.Get("name").String(), "chain_name")
		obj.Set(chainAsset.Get("icon_url").String(), "chain_icon")
		fmt.Printf("%s,", obj)
	}
}

// Only chain tokens
func MiniumChainlist() {
	top, _ := readFile("mixin-top-assets.json")
	assets := gjson.Get(top, "assets").Array()
	for _, asset := range assets {
		obj := gabs.New()
		obj.Set(asset.Get("asset_id").String(), "asset_id")
		obj.Set(asset.Get("name").String(), "name")
		obj.Set(asset.Get("symbol").String(), "symbol")
		obj.Set(asset.Get("icon_url").String(), "icon")
		obj.Set(asset.Get("chain_id").String(), "chain_id")
		obj.Set(asset.Get("asset_key").String(), "asset_key")
		if asset.Get("asset_id").String() == asset.Get("chain_id").String() {
			fmt.Printf("%s,", obj)
		}
	}
}

func AssetIDChainNameList() {

}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func AssetKeyList(name string) {
	// Asset symbol to Asset Key (Only ERC20)
	Assets, err := readFile("mixin-top-assets.json")
	if err != nil {
		log.Println(err)
	}
	assets := gjson.Get(Assets, "assets").Array()

	var keys []string
	o := gabs.New()
	// set btc as wbtc
	o.Set("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599", "btc")

	for _, asset := range assets {
		SYMBOL := asset.Get("symbol").String()
		NAME := asset.Get("name").String()
		KEY := asset.Get("asset_key").String()

		symbol := strings.ToLower(SYMBOL)

		if isLpToken(NAME) {
			continue
		}
		if isRings(NAME) {
			continue
		}

		// Skip non-ERC20
		if len(KEY) != 42 {
			continue
		}
		// Skip duplicated
		if contains(keys, symbol) {
			continue
		}
		keys = append(keys, symbol)

		obj := gabs.New()
		obj.Set(KEY, symbol)
		print(symbol, ":", KEY, "\n")
		o.Merge(obj)
	}
	writeFile(name, o.StringIndent("", " "))
	fmt.Println("Asset key list saved in", name)
}

func main() {
	/*
		PureTokenlist(NAMES[0])
		llamaTokenlist(NAMES[1])
		MVMChainList(NAMES[2])
		AssetKeyList(NAMES[3])
	*/
	MiniumTokenlist()
}

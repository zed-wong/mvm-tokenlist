package main

import (
	"log"
	"fmt"
	"context"
	"io/ioutil"

	"github.com/Jeffail/gabs"
	"github.com/go-resty/resty/v2"
	"github.com/fox-one/mixin-sdk-go"
)

const (
	NAME = "mvm-tokenlist.json"
	ENDPOINT = "https://api.mvm.dev/asset_contract?asset="
	NULL_ADDR = "0x0000000000000000000000000000000000000000"
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

func main() {
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
		fmt.Println(obj.StringIndent("", " "))
		o.Merge(obj)
	}
	writeFile(NAME, o.StringIndent("", " "))
	fmt.Println("Token info saved in", NAME)
}


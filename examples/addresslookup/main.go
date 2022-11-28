package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/go-ethutils/addresslookup"
	"github.com/metachris/go-ethutils/utils"
)

func main() {
	ethNodeUri := os.Getenv("ETH_NODE")
	if ethNodeUri == "" {
		fmt.Println("Please set ETH_NODE environment variable")
		os.Exit(1)
	}

	client, err := ethclient.Dial(ethNodeUri)
	utils.Perror(err)
	addressLookup := addresslookup.NewAddressLookupService(client)

	err = addressLookup.AddAllAddresses()
	utils.Perror(err)

	a, f := addressLookup.GetAddressDetail("0x3ecef08d0e2dad803847e052249bb4f8bff2d5bb") // MiningPoolHub
	fmt.Println(f, a)

	a, f = addressLookup.GetAddressDetail("0xdac17f958d2ee523a2206206994597c13d831ec7") // ERC20: Tether
	fmt.Println(f, a)

	a, f = addressLookup.GetAddressDetail("0x21a31Ee1afC51d94C2eFcCAa2092aD1028285549") // Binance
	fmt.Println(f, a)

	err = addressLookup.AddAddressesFromJsonUrl(addresslookup.JsonUrlEthplorerExchangeAddresses)
	utils.Perror(err)

	a, f = addressLookup.GetAddressDetail("0x2b5634c42055806a59e9107ed44d43c426e58258") // KuCoin
	fmt.Println(f, a)
}

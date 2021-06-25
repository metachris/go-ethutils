package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/go-ethutils/addresslookup"
	"github.com/metachris/go-ethutils/utils"
)

func main() {
	client, err := ethclient.Dial(os.Getenv("ETH_NODE"))
	utils.Perror(err)

	// block, err := client.BlockByNumber(context.Background(), big.NewInt(12691459))
	// utils.Perror(err)

	addressLookup := addresslookup.NewAddressLookupService(client)

	err = addressLookup.AddAddressFromDefaultJson()
	utils.Perror(err)

	a, f := addressLookup.GetAddressDetail("0x3ecef08d0e2dad803847e052249bb4f8bff2d5bb") // MiningPoolHub
	fmt.Println(f, a)

	a, f = addressLookup.GetAddressDetail("0xdac17f958d2ee523a2206206994597c13d831ec7") // ERC20: Tether
	fmt.Println(f, a)

	a, f = addressLookup.GetAddressDetail("0x21a31Ee1afC51d94C2eFcCAa2092aD1028285549") // Binance
	fmt.Println(f, a)
}

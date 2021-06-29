package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/metachris/go-ethutils/addresslookup"
	"github.com/metachris/go-ethutils/utils"
)

func main() {
	log.SetOutput(os.Stdout)

	addressPtr := flag.String("addr", "", "Address to look up")
	flag.Parse()

	if *addressPtr == "" {
		return
	}

	fmt.Println("Address:", *addressPtr)
	res, err := addresslookup.EthplorerServiceAddressLookup(*addressPtr)
	utils.Perror(err)
	fmt.Println("- IsContract:", res.IsContract)
	fmt.Println("- Tags:", strings.Join(res.PublicTags, ", "))
}

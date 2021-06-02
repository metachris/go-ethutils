package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/go-eth-utils/blockswithtx"
	"github.com/metachris/go-eth-utils/utils"
)

func perror(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	client, err := ethclient.Dial(os.Getenv("ETH_NODE"))
	perror(err)

	// blockswithtx.GetBlocksWithReceipts(client, processBlockWithReceipts, 12381372, 12381372, 5)
	t, err := utils.DateToTime("2021-05-20", 0, 0, 0)
	perror(err)

	h, err := utils.GetFirstBlockHeaderAtOrAfterTime(client, t)
	perror(err)
	fmt.Println(h.Number)
}

func processBlockWithReceipts(b *blockswithtx.BlockWithTxReceipts) {
	fmt.Println(*b)
}

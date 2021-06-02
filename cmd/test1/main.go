package main

import (
	"fmt"
	"os"
	"time"

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

	TestGetBlocksWithReceipts(client)
}

func TestGetBlocksWithReceipts(client *ethclient.Client) {
	processBlockWithReceipts := func(b *blockswithtx.BlockWithTxReceipts) {
		fmt.Println("block:", b.Block.Number(), "receipts:", len(b.TxReceipts))
	}

	var startBlock int64 = 12381372
	var endBlock int64 = startBlock + 20
	var numThreads int = 5

	t1 := time.Now()
	blockswithtx.GetBlocksWithTxReceipts(client, processBlockWithReceipts, startBlock, endBlock, numThreads)
	t2 := time.Since(t1)
	fmt.Println("needed", t2.Seconds(), "sec")
}

func TestGetBlockAtTime(client *ethclient.Client) {
	t, err := utils.DateToTime("2021-05-20", 0, 0, 0)
	perror(err)

	h, err := utils.GetFirstBlockHeaderAtOrAfterTime(client, t)
	perror(err)
	fmt.Println(h.Number)
}

package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/go-ethutils/blockswithtx"
	"github.com/metachris/go-ethutils/utils"
)

func main() {
	var startBlock int64 = 12600000
	var numBlocks int64 = 100
	var concurrency int = 10
	var ethNode string = os.Getenv("ETH_NODE")

	fmt.Println(ethNode, concurrency, startBlock, numBlocks)

	// Connect the geth client
	client, err := ethclient.Dial(ethNode)
	utils.Perror(err)

	// Create the channel to receive BlockWithTxReceipt
	blockChan := make(chan *blockswithtx.BlockWithTxReceipts, 100)

	// Create worker thread to process received items
	var lock sync.Mutex
	var numTx int64
	go func() {
		lock.Lock()
		defer lock.Unlock()
		for b := range blockChan {
			numTx += int64(len(b.Block.Transactions()))
			fmt.Println(b.Block.Number())
		}
	}()

	// Time retrieving the data
	t1 := time.Now()
	blockswithtx.GetBlocksWithTxReceipts(client, blockChan, startBlock, startBlock+numBlocks, concurrency)
	close(blockChan)
	lock.Lock() // wait until all blocks have been processed
	t2 := time.Since(t1)
	fmt.Printf("Processed %d transactions in %.3f seconds (%.2f tx/sec)\n", numTx, t2.Seconds(), float64(numTx)/t2.Seconds())
}

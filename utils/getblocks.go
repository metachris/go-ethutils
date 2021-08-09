package utils

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlocks is a fast block query pipeline. It queries blocks concurrently and pushes it into a channel for processing.
func GetBlocks(blockChan chan<- *types.Block, client *ethclient.Client, startBlock int64, endBlock int64, concurrency int) {
	var blockWorkerWg sync.WaitGroup         // for waiting for all workers to finish
	blockHeightChan := make(chan int64, 100) // channel for workers to know which heights to download

	// Start eth client thread pool
	for w := 1; w <= concurrency; w++ {
		blockWorkerWg.Add(1)

		// Worker gets a block height from blockHeightChan, downloads it, and puts it in the blockChan
		go func() {
			defer blockWorkerWg.Done()
			for blockHeight := range blockHeightChan {
				// fmt.Println(blockHeight)
				block, err := client.BlockByNumber(context.Background(), big.NewInt(blockHeight))
				if err != nil {
					log.Println("Error getting block:", blockHeight, err)
					continue
				}
				blockChan <- block
			}
		}()
	}

	// Push blockheights into channel, for workers to pick up
	for currentBlockNumber := startBlock; currentBlockNumber <= endBlock; currentBlockNumber++ {
		blockHeightChan <- currentBlockNumber
	}

	// Close worker channel and wait for workers to finish
	close(blockHeightChan)
	blockWorkerWg.Wait()
}

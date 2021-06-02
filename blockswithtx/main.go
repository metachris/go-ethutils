package blockswithtx

import (
	"context"
	"errors"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockWithTxReceipts struct {
	Block      *types.Block
	TxReceipts map[common.Hash]*types.Receipt
}

// GetBlockWithTxReceipts returns a single block with receipts for all transactions
func GetBlockWithTxReceipts(client *ethclient.Client, height int64) (res BlockWithTxReceipts, err error) {
	res.TxReceipts = make(map[common.Hash]*types.Receipt)

	// Get the block
	res.Block, err = client.BlockByNumber(context.Background(), big.NewInt(height))
	if err != nil {
		return res, err
	}

	// Get receipts for all transactions
	for _, tx := range res.Block.Transactions() {
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				// can apparently happen if 0 tx: https://etherscan.io/block/10102170
				continue
			}
			return res, err

		}
		res.TxReceipts[tx.Hash()] = receipt
	}

	return res, nil
}

// GetBlocksWithTxReceipts downloads a range of blocks with tx receipts and sends them to a user-defined function for processing.
// Uses numThreads concurrent geth connections to speed things up. 5 is usually a good number for a direct IPC connection.
func GetBlocksWithTxReceipts(client *ethclient.Client, callback func(b *BlockWithTxReceipts), startBlock int64, endBlock int64, numThreads int) {
	var blockWorkerWg sync.WaitGroup
	blockHeightChan := make(chan int64, 100)          // blockHeight to fetch with receipts
	blockChan := make(chan *BlockWithTxReceipts, 100) // channel for resulting BlockWithTxReceipt

	// Start eth client thread pool
	for w := 1; w <= numThreads; w++ {
		blockWorkerWg.Add(1)
		go func() {
			defer blockWorkerWg.Done()
			for blockHeight := range blockHeightChan {
				res, err := GetBlockWithTxReceipts(client, blockHeight)
				if err != nil {
					log.Println("Error getting block with tx receipts:", err)
					continue
				}
				blockChan <- &res
			}
		}()
	}

	// Start thread to pass blocks back to caller
	var processLock sync.Mutex
	processLock.Lock()
	go func() {
		defer processLock.Unlock() // we unlock when done
		for block := range blockChan {
			callback(block)
		}
	}()

	// Push blocks into channel, for workers to pick up
	for currentBlockNumber := startBlock; currentBlockNumber <= endBlock; currentBlockNumber++ {
		blockHeightChan <- currentBlockNumber
	}

	// Close worker channel and wait for workers to finish
	close(blockHeightChan)
	blockWorkerWg.Wait()

	// Close callback channel and wait for processing to finish
	close(blockChan)
	processLock.Lock()
}

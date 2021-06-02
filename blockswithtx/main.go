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
func GetBlockWithTxReceipts(client *ethclient.Client, height int64) (res *BlockWithTxReceipts, err error) {
	res = &BlockWithTxReceipts{}
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

// GetBlocksWithTxReceipts downloads a range of blocks with tx receipts and sends them to a user-defined function for processing
// Uses concurrency parallel connections to get data from the eth node fast. 5 is usually a good number for a direct IPC connection.
func GetBlocksWithTxReceipts(client *ethclient.Client, blockChan chan<- *BlockWithTxReceipts, startBlock int64, endBlock int64, concurrency int) {
	var blockWorkerWg sync.WaitGroup
	blockHeightChan := make(chan int64, 100) // blockHeight to fetch with receipts

	// Start eth client thread pool
	for w := 1; w <= concurrency; w++ {
		blockWorkerWg.Add(1)
		go func() {
			defer blockWorkerWg.Done()
			for blockHeight := range blockHeightChan {
				res, err := GetBlockWithTxReceipts(client, blockHeight)
				if err != nil {
					log.Println("Error getting block with tx receipts:", err)
					continue
				}
				blockChan <- res
			}
		}()
	}

	// Push blocks into channel, for workers to pick up
	for currentBlockNumber := startBlock; currentBlockNumber <= endBlock; currentBlockNumber++ {
		blockHeightChan <- currentBlockNumber
	}

	// Close worker channel and wait for workers to finish
	close(blockHeightChan)
	blockWorkerWg.Wait()

	// Close blockChan
	close(blockChan)
}

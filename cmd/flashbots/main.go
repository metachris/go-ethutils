package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
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
	blockHeightPtr := flag.Int("block", 0, "specific block to check")
	datePtr := flag.String("date", "", "date (yyyy-mm-dd or -1d)")
	hourPtr := flag.Int("hour", 0, "hour (UTC)")
	minPtr := flag.Int("min", 0, "hour (UTC)")
	lenPtr := flag.String("len", "", "num blocks or timespan (4s, 5m, 1h, ...)")
	watchPtr := flag.Bool("watch", false, "watch and process new blocks")
	flag.Parse()

	client, err := ethclient.Dial(os.Getenv("ETH_NODE"))
	perror(err)

	if len(*datePtr) != 0 || *blockHeightPtr != 0 {
		// A start for historical analysis was given
		// log.Fatal("Missing start (date or block). Add with -date <yyyy-mm-dd> or -block <blockNum>")
		processHistoricBlocks(*blockHeightPtr, *datePtr, *hourPtr, *minPtr, *lenPtr)
	}

	if *watchPtr {
		watch(client)
	}

}

func processHistoricBlocks(blockHeight int, date string, hour int, min int, len string) {
	// // Parse -len argument (can be either 1s, 5m, 2h or 4d, or without suffix a number of blocks)
	// numBlocks := 1
	// timespanSec := 0
	// switch {
	// case strings.HasSuffix(*lenPtr, "s"):
	// 	timespanSec, _ = strconv.Atoi(strings.TrimSuffix(*lenPtr, "s"))
	// case strings.HasSuffix(*lenPtr, "m"):
	// 	timespanSec, _ = strconv.Atoi(strings.TrimSuffix(*lenPtr, "m"))
	// 	timespanSec *= 60
	// case strings.HasSuffix(*lenPtr, "h"):
	// 	timespanSec, _ = strconv.Atoi(strings.TrimSuffix(*lenPtr, "h"))
	// 	timespanSec *= 60 * 60
	// case strings.HasSuffix(*lenPtr, "d"):
	// 	timespanSec, _ = strconv.Atoi(strings.TrimSuffix(*lenPtr, "d"))
	// 	timespanSec *= 60 * 60 * 24
	// case len(*lenPtr) == 0:
	// 	numBlocks = 1
	// default:
	// 	// No suffix: number of blocks
	// 	numBlocks, _ = strconv.Atoi(*lenPtr)
	// }

	// startBlockHeight := int64(*blockHeightPtr)
	// var startTimestamp int64
	// var startTime time.Time

	// // if *blockHeightPtr > 0 { // start at timestamp
	// // 	startBlockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(*blockHeightPtr)))
	// // 	core.Perror(err)
	// // 	startTimestamp = int64(startBlockHeader.Time)
	// // 	startTime = time.Unix(startTimestamp, 0)

	// // startHeight := int64(*blockHeightPtr)
	// // endHeight := startHeight + int64(numBlocks) - 1
	// // checkBlocks(startHeight, endHeight)
}

func checkBlocks(startHeight int64, endHeight int64) {
	client, err := ethclient.Dial(os.Getenv("ETH_NODE"))
	perror(err)

	processBlockWithReceipts := func(b *blockswithtx.BlockWithTxReceipts) {
		// fmt.Println("block", b.Block.Number(), "\t receipts:", len(b.TxReceipts))
		for _, tx := range b.Block.Transactions() {
			receipt := b.TxReceipts[tx.Hash()]
			if receipt == nil {
				continue
			}
			if receipt.Status == 1 {
				continue
			}

			// Failed transaction
			if len(tx.Data()) == 0 {
				sender, _ := utils.GetTxSender(tx)
				fmt.Printf("Failed Flashbots tx in block %v: %s from %v\n", b.Block.Number(), tx.Hash(), sender)
			}
		}
	}

	t1 := time.Now()
	blockswithtx.GetBlocksWithTxReceipts(client, processBlockWithReceipts, startHeight, endHeight, 5)
	t2 := time.Since(t1)
	fmt.Println("Processed", endHeight-startHeight, "blocks in", t2.Seconds(), "sec")
}

func watch(client *ethclient.Client) {
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			b, err := blockswithtx.GetBlockWithTxReceipts(client, header.Number.Int64())
			perror(err)
			checkBlockWithReceipts(b)
		}
	}
}

func checkBlockWithReceipts(b blockswithtx.BlockWithTxReceipts) {
	fmt.Println("block", b.Block.Number(), "\t receipts:", len(b.TxReceipts))
	for _, tx := range b.Block.Transactions() {
		receipt := b.TxReceipts[tx.Hash()]
		if receipt == nil {
			continue
		}

		if tx.GasPrice().Uint64() == 0 && len(tx.Data()) > 0 {
			sender, _ := utils.GetTxSender(tx)
			if receipt.Status == 1 { // successful tx
				// fmt.Printf("Flashbots tx in block %v: %s from %v\n", b.Block.Number(), tx.Hash(), sender)
			} else { // failed tx
				fmt.Printf("Failed Flashbots tx in block %v: %s from %v\n", b.Block.Number(), tx.Hash(), sender)
			}
		}
	}
}

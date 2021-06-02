package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
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

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func colorPrintf(color string, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	fmt.Printf(string(color), str)
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

	if *datePtr != "" || *blockHeightPtr != 0 {
		// A start for historical analysis was given
		// log.Fatal("Missing start (date or block). Add with -date <yyyy-mm-dd> or -block <blockNum>")
		startBlock, endBlock := getBlockRangeFromArguments(client, *blockHeightPtr, *datePtr, *hourPtr, *minPtr, *lenPtr)
		// fmt.Println(startBlock, endBlock)
		checkBlocks(client, startBlock, endBlock)
	}

	if *watchPtr {
		watch(client)
	}

}

func getBlockRangeFromArguments(client *ethclient.Client, blockHeight int, date string, hour int, min int, length string) (startBlock int64, endBlock int64) {
	if date != "" && blockHeight != 0 {
		panic("cannot use both -block and -date arguments")
	}

	if date == "" && blockHeight == 0 {
		panic("need to use either -block or -date arguments")
	}

	// Parse -len argument. Can be either a number of blocks or a time duration (eg. 1s, 5m, 2h or 4d)
	numBlocks := 0
	timespanSec := 0
	var err error
	switch {
	case strings.HasSuffix(length, "s"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "s"))
		perror(err)
	case strings.HasSuffix(length, "m"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "m"))
		perror(err)
		timespanSec *= 60
	case strings.HasSuffix(length, "h"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "h"))
		perror(err)
		timespanSec *= 60 * 60
	case strings.HasSuffix(length, "d"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "d"))
		perror(err)
		timespanSec *= 60 * 60 * 24
	case length == "": // default 1 block
		numBlocks = 1
	default: // No suffix: number of blocks
		numBlocks, err = strconv.Atoi(length)
		perror(err)
	}

	// startTime is set from date argument, or block timestamp if -block argument was used
	var startTime time.Time

	// Get start block
	if blockHeight > 0 { // start at block height
		startBlockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blockHeight)))
		perror(err)
		startBlock = startBlockHeader.Number.Int64()
		startTime = time.Unix(int64(startBlockHeader.Time), 0)
	} else {
		// Negative date prefix (-1d, -2m, -1y)
		if strings.HasPrefix(date, "-") {
			if strings.HasSuffix(date, "d") {
				t := time.Now().AddDate(0, 0, -1)
				startTime = t.Truncate(24 * time.Hour)
			} else if strings.HasSuffix(date, "m") {
				t := time.Now().AddDate(0, -1, 0)
				startTime = t.Truncate(24 * time.Hour)
			} else if strings.HasSuffix(date, "y") {
				t := time.Now().AddDate(-1, 0, 0)
				startTime = t.Truncate(24 * time.Hour)
			} else {
				panic(fmt.Sprintf("Not a valid date offset: %s", date))
			}
		} else {
			startTime, err = utils.DateToTime(date, hour, min, 0)
			perror(err)
		}

		startBlockHeader, err := utils.GetFirstBlockHeaderAtOrAfterTime(client, startTime)
		perror(err)
		startBlock = startBlockHeader.Number.Int64()
	}

	if numBlocks > 0 {
		endBlock = startBlock + int64(numBlocks-1)
	} else if timespanSec > 0 {
		endTime := startTime.Add(time.Duration(timespanSec) * time.Second)
		// fmt.Printf("endTime: %v\n", endTime.UTC())
		endBlockHeader, _ := utils.GetFirstBlockHeaderAtOrAfterTime(client, endTime)
		endBlock = endBlockHeader.Number.Int64() - 1
	} else {
		panic("No valid block range")
	}

	if endBlock < startBlock {
		endBlock = startBlock
	}

	return startBlock, endBlock
}

func checkBlocks(client *ethclient.Client, startHeight int64, endHeight int64) {
	fmt.Printf("Checking blocks %d to %d...\n", startHeight, endHeight)
	t1 := time.Now()
	blockChan := make(chan *blockswithtx.BlockWithTxReceipts, 100) // channel for resulting BlockWithTxReceipt

	// Start thread listening for blocks (with tx receipts) from geth worker pool
	var numTx uint64
	go func() {
		for b := range blockChan {
			checkBlockWithReceipts(b)
			numTx += uint64(len(b.Block.Transactions()))
		}
	}()

	blockswithtx.GetBlocksWithTxReceipts(client, blockChan, startHeight, endHeight, 5)
	t2 := time.Since(t1)
	fmt.Printf("Processed %d blocks (%d transactions) in %.3f seconds\n", endHeight-startHeight, numTx, t2.Seconds())
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

func checkBlockWithReceipts(b *blockswithtx.BlockWithTxReceipts) {
	// fmt.Printf("block %v: %d tx\n", b.Block.Number(), len(b.Block.Transactions()))
	utils.PrintBlock(b.Block)
	for _, tx := range b.Block.Transactions() {
		receipt := b.TxReceipts[tx.Hash()]
		if receipt == nil {
			continue
		}

		if utils.IsBigIntZero(tx.GasPrice()) && len(tx.Data()) > 0 {
			sender, _ := utils.GetTxSender(tx)
			if receipt.Status == 1 { // successful tx
				// fmt.Printf("Flashbots tx in block %v: %s from %v\n", b.Block.Number(), tx.Hash(), sender)
			} else { // failed tx
				// fmt.Printf("block %v: failed Flashbots tx %s from %v\n", WarningColor, b.Block.Number(), tx.Hash(), sender)
				colorPrintf(ErrorColor, "block %v: failed Flashbots tx %s from %v\n", b.Block.Number(), tx.Hash(), sender)
			}
		}
	}
}

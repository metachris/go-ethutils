package utils

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlockRangeFromArguments returns start and end blocks for a given block or date range. startBlock is first block at or after the given time, endBlock is the last before the given end time.
func FindBlockRange(client *ethclient.Client, blockHeight int, date string, hour int, min int, length string) (startBlock int64, endBlock int64, err error) {
	if date != "" && blockHeight != 0 {
		return startBlock, endBlock, errors.New("cannot use both block and date arguments")
	}

	if date == "" && blockHeight == 0 {
		return startBlock, endBlock, errors.New("need to use either block or date arguments")
	}

	// length can be a number of blocks (eg. 40) or a time duration (eg. 1s, 5m, 2h or 4d) or "." meaning the latest block
	numBlocks := 0
	timespanSec := 0
	switch {
	case length == ".":
		// is handled further down
	case strings.HasSuffix(length, "s"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "s"))
	case strings.HasSuffix(length, "m"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "m"))
		timespanSec *= 60
	case strings.HasSuffix(length, "h"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "h"))
		timespanSec *= 60 * 60
	case strings.HasSuffix(length, "d"):
		timespanSec, err = strconv.Atoi(strings.TrimSuffix(length, "d"))
		timespanSec *= 60 * 60 * 24
	case length == "": // default 1 block
		numBlocks = 1
	default: // No suffix: number of blocks
		numBlocks, err = strconv.Atoi(length)
	}

	if err != nil {
		return startBlock, endBlock, err
	}

	// startTime is set from date argument, or block timestamp if -block argument was used
	var startTime time.Time

	// Get start block
	if blockHeight > 0 { // start at block height
		startBlockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blockHeight)))
		if err != nil {
			return startBlock, endBlock, err
		}

		startBlock = startBlockHeader.Number.Int64()
		startTime = time.Unix(int64(startBlockHeader.Time), 0)
	} else {
		// Negative date prefix (-1d, -2m, -1y)
		if strings.HasPrefix(date, "-") {
			if strings.HasSuffix(date, "d") {
				days, _ := strconv.Atoi(date[1 : len(date)-1])
				t := time.Now().AddDate(0, 0, -days) // todo
				startTime = t.Truncate(24 * time.Hour)
			} else if strings.HasSuffix(date, "m") {
				months, _ := strconv.Atoi(date[1 : len(date)-1])
				t := time.Now().AddDate(0, -months, 0)
				startTime = t.Truncate(24 * time.Hour)
			} else if strings.HasSuffix(date, "y") {
				years, _ := strconv.Atoi(date[1 : len(date)-1])
				t := time.Now().AddDate(-years, 0, 0)
				startTime = t.Truncate(24 * time.Hour)
			} else {
				return startBlock, endBlock, fmt.Errorf("not a valid date offset: '%s'. Can be d, m, y", date)
			}
		} else {
			startTime, err = DateToTime(date, hour, min, 0)
			if err != nil {
				return startBlock, endBlock, err
			}
		}

		startBlockHeader, err := GetFirstBlockHeaderAtOrAfterTime(client, startTime)
		if err != nil {
			return startBlock, endBlock, err
		}
		startBlock = startBlockHeader.Number.Int64()
	}

	// Find end block
	if length == "." {
		latestBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return startBlock, endBlock, err
		}
		endBlock = latestBlockHeader.Number.Int64()
	} else if numBlocks > 0 {
		endBlock = startBlock + int64(numBlocks-1)
	} else if timespanSec > 0 {
		endTime := startTime.Add(time.Duration(timespanSec) * time.Second)
		endBlockHeader, err := GetFirstBlockHeaderAtOrAfterTime(client, endTime)
		if err != nil {
			return startBlock, endBlock, err
		}
		endBlock = endBlockHeader.Number.Int64() - 1
	} else {
		return startBlock, endBlock, errors.New("no valid block range")
	}

	if endBlock < startBlock {
		endBlock = startBlock
	}

	return startBlock, endBlock, nil
}

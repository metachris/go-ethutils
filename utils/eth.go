package utils

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetTxSender(tx *types.Transaction) (from common.Address, err error) {
	from, err = types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		from, err = types.Sender(types.HomesteadSigner{}, tx)
	}
	return from, err
}

// Roughly estimate a block number by target timestamp (might be off by a lot)
func EstimateTargetBlocknumber(utcTimestamp int64) int64 {
	// Calculate block-difference from reference block
	referenceBlockNumber := int64(12323940)
	referenceBlockTimestamp := int64(1619546404) // 2021-04-27 19:49:13 +0200 CEST

	secDiff := referenceBlockTimestamp - utcTimestamp
	// fmt.Println("secDiff", secDiff)
	blocksDiff := secDiff / 13
	targetBlock := referenceBlockNumber - blocksDiff
	return targetBlock
}

// GetBlockHeaderAtTimestamp returns the header of the first block at or after the timestamp. If timestamp is after
// latest block, then return latest block.
func GetFirstBlockHeaderAtOrAfterTime(client *ethclient.Client, targetTime time.Time) (header *types.Header, err error) {
	// Get latest header
	latestBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return header, err
	}

	// Ensure target timestamp is before latest block
	targetTimestamp := targetTime.Unix()
	if uint64(targetTimestamp) > latestBlockHeader.Time {
		return header, errors.New("target timestamp after latest block")
	}

	// Estimate a target block number
	currentBlockNumber := EstimateTargetBlocknumber(targetTimestamp)

	// If estimation later than latest block, then use latest block as estimation base
	if currentBlockNumber > latestBlockHeader.Number.Int64() {
		currentBlockNumber = latestBlockHeader.Number.Int64()
	}

	// approach the target block from below, to be sure it's the first one at/after the timestamp
	var isNarrowingDownFromBelow = false

	// Ringbuffer for the latest secDiffs, to avoid going in circles when narrowing down
	lastSecDiffs := make([]int64, 7)
	lastSecDiffsIncludes := func(a int64) bool {
		for _, b := range lastSecDiffs {
			if b == a {
				return true
			}
		}
		return false
	}

	// fmt.Printf("Finding start block:\n")
	var secDiff int64
	blockSecAvg := int64(13) // average block time. is adjusted when narrowing down

	for {
		// core.DebugPrintln("Checking block:", currentBlockNumber)
		blockNumber := big.NewInt(currentBlockNumber)
		header, err := client.HeaderByNumber(context.Background(), blockNumber)
		if err != nil {
			return header, err
		}

		secDiff = int64(header.Time) - targetTimestamp

		// fmt.Printf("%d \t blockTime: %d / %v \t secDiff: %5d\n", currentBlockNumber, header.Time, time.Unix(int64(header.Time), 0).UTC(), secDiff)

		// Check if this secDiff was already seen (avoid circular endless loop)
		if lastSecDiffsIncludes(secDiff) && blockSecAvg < 25 {
			blockSecAvg += 1
			// fmt.Println("- Increase blockSecAvg to", blockSecAvg)
		}

		// Pop & add secDiff to array of last values
		lastSecDiffs = lastSecDiffs[1:]
		lastSecDiffs = append(lastSecDiffs, secDiff)
		// core.DebugPrintln("lastSecDiffs:", lastSecDiffs)

		if Abs(secDiff) < 80 || isNarrowingDownFromBelow { // getting close
			if secDiff < 0 {
				// still before wanted startTime. Increase by 1 from here...
				isNarrowingDownFromBelow = true
				currentBlockNumber += 1
				continue
			}

			// Only return if coming block-by-block from below, making sure to take first block after target time
			if isNarrowingDownFromBelow {
				return header, nil
			} else {
				currentBlockNumber -= 1
				continue
			}
		}

		// Try for better block in big steps
		blockDiff := secDiff / blockSecAvg
		currentBlockNumber -= blockDiff
	}
}

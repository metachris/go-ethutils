## Go Ethereum Utilities

Helpers for working with Ethereum:

**Links**

* Docs: https://pkg.go.dev/github.com/metachris/go-ethutils
* Example usage: https://github.com/metachris/flashbots

**Contents**

* [blockswithtx](https://github.com/metachris/go-ethutils/blob/master/blockswithtx) - fast, concurrent block+receipts downloading pipeline (use a geth IPC connection)
* [smartcontracts](https://github.com/metachris/go-ethutils/blob/master/smartcontracts) - detect types of smart contracts
* [addressdetail](https://github.com/metachris/go-ethutils/blob/master/addressdetail) - helper type for smart contracts and addresses
* [utils/eth.go](https://github.com/metachris/go-ethutils/blob/master/utils/eth.go) - finding first block at or after a certain UTC timestamp
* [utils/blockrangefinder.go](https://github.com/metachris/go-ethutils/blob/master/utils/blockrangefinder.go) - find a block range based on date, timespans or blocks
* [utils/various.go](https://github.com/metachris/go-ethutils/blob/master/utils/various.go) - various utilities

**Feedback**

* Reach out to [twitter.com/metachris](https://twitter.com/metachris)

---

## BlockWithTxReceipts

`blockswithtx` is used for concurrently fetching blocks and tx-receipts from a geth node.

Benchmarks - running on the same machine as the geth node:

```
# IPC connection with concurrency 1, 5, 10, 15
ipc+1x:  Processed 17773 transactions in 29.671 seconds (599.01 tx/sec)
ipc+5x:  Processed 17773 transactions in 7.453 seconds (2384.59 tx/sec)
ipc+10x: Processed 17773 transactions in 5.403 seconds (3289.71 tx/sec)
ipc+15x: Processed 17773 transactions in 4.631 seconds (3837.83 tx/sec)

# WebSocket connection with concurrency 1, 5, 10, 15
ws+1x:  Processed 17773 transactions in 29.784 seconds (596.74 tx/sec)
ws+5x:  Processed 17773 transactions in 8.380 seconds (2120.87 tx/sec)
ws+10x: Processed 17773 transactions in 5.763 seconds (3084.00 tx/sec)
ws+15x: Processed 17773 transactions in 4.842 seconds (3670.32 tx/sec)

# HTTP connection with concurrency 1, 5, 10, 15
http+1x:  Processed 17773 transactions in 33.829 seconds (525.39 tx/sec)
http+5x:  Processed 17773 transactions in 8.828 seconds (2013.30 tx/sec)
http+10x: Processed 17773 transactions in 5.814 seconds (3056.90 tx/sec)
http+15x: Processed 17773 transactions in 5.510 seconds (3225.56 tx/sec)
```

Over the network I could only get ~200 tx/sec.

Example code: cmd/benchmark-blockswithtx/main.go

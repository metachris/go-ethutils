## Go helpers for dealing with Ethereum and geth

Various helpers that make my life easier.

There is an example for getting failed Flashbots transactions:

```bash
# Subscribe to new blocks and find failed Flashbots tx:
go run cmd/flashbots/main.go -watch

# Historic, using a starting block
go run cmd/flashbots/main.go -block 12539827           # 1 block
go run cmd/flashbots/main.go -block 12539827 -len 5    # 5 blocks
go run cmd/flashbots/main.go -block 12539827 -len 10m  # all blocks within 10 minutes of given block
go run cmd/flashbots/main.go -block 12539827 -len 1h   # all blocks within 1 hour of given block
go run cmd/flashbots/main.go -block 12539827 -len 1d   # all blocks within 1 day of given block

# Historic, using a starting date
go run cmd/flashbots/main.go -date -1d -len 1h         # all blocks within 1 hour of yesterday 00:00:00 (UTC)
go run cmd/flashbots/main.go -date 2021-05-25 -len 1h  # all blocks within 1 hour of given date 00:00:00 (UTC)
go run cmd/flashbots/main.go -date 2021-05-25 -hour 12 -min 5 -len 1h  # all blocks within 1 hour of given date 12:05:00 (UTC)
```

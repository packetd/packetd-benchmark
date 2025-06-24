module github.com/packetd/packetd-benchmark/redis/client

go 1.24

require github.com/redis/go-redis/v9 v9.7.1

require (
	github.com/olekukonko/tablewriter v1.0.7
	github.com/packetd/packetd-benchmark/common v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/olekukonko/errors v0.0.0-20250405072817-4e6d85265da6 // indirect
	github.com/olekukonko/ll v0.0.8 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

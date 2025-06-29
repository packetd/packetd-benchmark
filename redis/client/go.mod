module github.com/packetd/packetd-benchmark/redis/client

go 1.24

require github.com/redis/go-redis/v9 v9.7.1

require (
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/packetd/packetd-benchmark/common v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

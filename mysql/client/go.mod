module github.com/packetd/packetd-benchmark/mysql/client

go 1.24

require github.com/go-sql-driver/mysql v1.9.1

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

require (
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/packetd/packetd-benchmark/common v0.0.0
)

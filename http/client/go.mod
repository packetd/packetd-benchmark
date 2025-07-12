module github.com/packetd/packetd-benchmark/http/client

go 1.24

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

require (
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/packetd/packetd-benchmark/common v0.0.0
)

require (
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

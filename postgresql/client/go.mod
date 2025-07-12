module github.com/packetd/packetd-benchmark/postgresql/client

go 1.24

require github.com/jackc/pgx/v5 v5.7.4

require (
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/packetd/packetd-benchmark/common v0.0.0
)

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

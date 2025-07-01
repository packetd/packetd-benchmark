module github.com/packetd/packetd-benchmark/grpc/client

go 1.24

replace github.com/packetd/packetd-benchmark/grpc/pb => ./../pb

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

require (
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/packetd/packetd-benchmark/common v0.0.0
	github.com/packetd/packetd-benchmark/grpc/pb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.71.0
)

require (
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

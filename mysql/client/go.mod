module github.com/packetd/packetd-benchmark/mysql/client

go 1.23.4

require github.com/go-sql-driver/mysql v1.9.1

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/packetd/packetd-benchmark/common v0.0.0 => ./../../common

require github.com/packetd/packetd-benchmark/common v0.0.0

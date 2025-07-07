# gGRPC 压测

1）Running Server

```shell
$ go run server/main.go
```

2）Client Usage

```shell
$ ./client -h
Usage of ./client:
  -addr string
        grpc server address (default "localhost:8085")
  -body_size string
        request body size (default "1KB")
  -interval duration
        interval per request
  -total int
        requests total (default 1)
  -workers int
        concurrency workers (default 1)
```

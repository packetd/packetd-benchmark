# HTTP 压测

1）Running Server

```shell
$ go run server/main.go
```

2）Client Usage

```shell
$ ./client -h
Usage of ./client:
  -addr string
        http server address (default "localhost:8083")
  -body_size string
        request body size (default "1KB")
  -interval duration
        interval per request
  -status string
        http response status (default "200")
  -total int
        requests total (default 1)
  -workers int
        concurrency workers (default 1)
```

# Redis 压测

Client Usage

```shell
$ ./client -h                                                             
Usage of ./client:
  -addr string
        redis server address (default "localhost:6379")
  -body_size string
        request body size (default "1KB")
  -cmd string
        redis command, options: ping/set/get (default "ping")
  -interval duration
        interval per request
  -total int
        requests total (default 1)
  -workers int
        concurrency workers (default 1)
```

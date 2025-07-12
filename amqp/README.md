# AMQP 压测

## Producer Usage

```shell
$ ./producer -h
Usage of ./producer:
  -dsn string
        rabbitmq dsn connected to (default "amqp://guest:guest@localhost:5672/")
  -message_size int
        rabbitmq message size in bytes (default 10240)
  -queue string
        rabbitmq queue name (default "benchmark")
  -total int
        require produce messages count (default 100)
```

## Consumer Usage

```shell
$ ./consumer -h
Usage of ./consumer:
  -dsn string
        kafka bootstrap brokers, as a comma separated list (default "amqp://guest:guest@localhost:5672/")
  -queue string
        kafka topics to be consumed (default "benchmark")
  -total int
        require produce messages count (default 100)
```

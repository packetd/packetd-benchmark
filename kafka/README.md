# Kafka 压测

## Producer Usage

```shell
$ ./producer -h 
Usage of ./producer:
  -brokers string
        kafka bootstrap brokers, as a comma separated list (default "localhost:9092")
  -message_size int
        kafka message size in bytes (default 10240)
  -topic string
        kafka topics to be consumed (default "benchmark")
  -total int
        require produce messages count (default 100)
```

## Consumer Usage

```shell
$ ./consumer -h  
Usage of ./consumer:
  -brokers string
        kafka bootstrap brokers, as a comma separated list (default "localhost:9092")
  -topic string
        kafka topics to be consumed (default "benchmark")
  -total int
        require consume messages count (default 100)
```

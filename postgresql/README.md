# PostgreSQL 压测

## Prepare for data

1）建表语句

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE stress_test (
    id            BIGSERIAL PRIMARY KEY,
    uuid          UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id       INTEGER NOT NULL DEFAULT 0,
    amount        NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    create_time   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status        SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX idx_user_id ON stress_test (user_id);
CREATE INDEX idx_create_time ON stress_test (create_time);
CREATE INDEX idx_status ON stress_test (status);
```

2）创建记录

```sql
BEGIN;

INSERT INTO stress_test (uuid, user_id, amount, status, create_time)
SELECT
gen_random_uuid(),
floor(random() * 100000)::integer + 1,
round(random() * 9999.99, 2),
floor(random() * 4)::smallint,
current_timestamp - (random() * interval '365 days')
FROM generate_series(1, 10000); 

COMMIT;
```

## Usage

```shell
$ ./client -h
Usage of ./client:
  -dsn string
    	mysql server dsn
  -interval duration
    	interval per request
  -sql string
    	sql statement
  -total int
    	requests total (default 1)
  -workers int
    	concurrency workers (default 1)

# ./client -dsn 'postgres://postgres@localhost:5432/benchmark?sslmode=disable' -sql 'select * from stress_test limit 5000;' -workers 5 -total 10000
```

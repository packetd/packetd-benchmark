# MySQL 压测

## 前置准备

1）建表语句

```sql
CREATE TABLE `stress_test` (
	`id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
	`uuid` char(36) NOT NULL DEFAULT '' COMMENT 'UUID',
	`user_id` int UNSIGNED NOT NULL DEFAULT '0' COMMENT 'USER ID',
	`amount` decimal(10, 2) NOT NULL DEFAULT '0.00' COMMENT '金额',
	`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
	PRIMARY KEY (`id`),
	KEY `idx_user_id` (`user_id`),
	KEY `idx_create_time` (`create_time`),
	KEY `idx_status` (`status`)
) ENGINE = InnoDB CHARSET = utf8mb4 COMMENT '压测表';
```

2）创建记录

```python
import pymysql
from faker import Faker
import uuid
import random

fake = Faker()
conn = pymysql.connect(host='localhost', user='root', password='', db='benchmark')

with conn.cursor() as cursor:
    for i in range(10000):
        user_id = random.randint(1, 100000)
        amount = round(random.uniform(0.01, 9999.99), 2)
        status = random.randint(0, 3)
        
        sql = """
        INSERT INTO stress_test (uuid, user_id, amount, status)
        VALUES (%s, %s, %s, %s)
        """
        cursor.execute(sql, (str(uuid.uuid4()), user_id, amount, status))
        
        # 每 1000 条提交一次
        if i % 1000 == 0:
            conn.commit()

conn.commit()
conn.close()
```

## 工具用法

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
```

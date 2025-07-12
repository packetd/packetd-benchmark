# MongoDB 压测

## Prepare for data

1）建表语句

```sql
use benchmark;

var startTime = new Date();
for (var i = 1; i <= 90000; i++) {
var randomValue = Math.random() * 1000;
var isActive = i % 2 === 0;
var tags = ["tag" + (i % 5 + 1), "tag" + (i % 3 + 1)];
if (i % 7 === 0) tags.push("special");
    db.stress_test.insertOne({
        user_id: i,
        username: "user_" + i,
        email: "user_" + i + "@example.com",
        created_at: new Date(),
        active: isActive,
        score: Math.floor(randomValue),
        balance: (randomValue * 10).toFixed(2),
        tags: tags,
        last_login: isActive ? new Date() : null,
        metadata: {
            device: i % 3 === 0 ? "mobile" : "desktop",
            version: "1." + (i % 5)
        }
    });

    if (i % 1000 === 0) {
        print("Inserted " + i + " records");
    }
}
var endTime = new Date();
print("Inserted 10000 records in " + (endTime - startTime) + " ms");
```

## Usage

```shell
$ ./client -h  
Usage of ./client:
  -collection string
        collection name
  -database string
        database name
  -dsn string
        mysql server dsn
  -interval duration
        interval per request
  -limit int
        records count
  -total int
        requests total (default 1)
  -workers int
        concurrency workers (default 1)
        
# ./client -dsn 'mongodb://admin:admin@localhost:27017' -database benchmark -collection stress_test -workers 20 -limit 1000 -total 50000
```

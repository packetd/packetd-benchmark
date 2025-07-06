```python
import pymysql
from faker import Faker
import uuid
import random

fake = Faker()
conn = pymysql.connect(host='localhost', user='root', password='', db='benchmark')

with conn.cursor() as cursor:
    # 生成 100 万条数据
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
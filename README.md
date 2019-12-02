## change DB Connection
コンテナの中に入る
```
$ docker exec -it [NAMES] bash
mysql> mysql -u root -p
```

max connection の数を確認
```
mysql> show variables like "%max_connections%";
+-----------------+-------+
| Variable_name   | Value |
+-----------------+-------+
| max_connections | 150   |
+-----------------+-------+
1 row in set (0.00 sec)
```

max connection の数変更 ( 今回は10にした )
```
set global max_connections = 10;
```

処理中の接続の確認
```
show processlist;
```
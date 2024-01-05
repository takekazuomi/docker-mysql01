# MySQL+Docker+go

やってることは、下記のような感じ

- 開発用のubuntuとmysqlのコンテナを上げる
- 国土数値情報/医療機関データをダウンロードしてきて、mysqlに入れる
- goのコードから、近くの医療機関を検索する

## 動かし方

1. `make bootstrap` すると、国土数値情報の医療機関データをダウンロードして、mysqlに入れてくれる
2. `make query` とすると、四谷の事務所近辺の医療機関を検索してくれる

**例:** `make bootstrap`

```sh
$ make bootstrap
docker compose -f docker-compose.yml build
[+] Building 1.0s (9/9) FINISHED                                                                                          docker:default
 => [dev internal] load .dockerignore                                                                                               0.0s
 => => transferring context: 2B                                                                                                     0.0s
 => [dev internal] load build definition from Dockerfile                                                                            0.0s
 => => transferring dockerfile: 584B                                                                                                0.0s
 => [dev internal] load metadata for docker.io/library/golang:1.18.1-bullseye                                                       0.9s
 => [dev 1/5] FROM docker.io/library/golang:1.18.1-bullseye@sha256:3b1a72af045ad0fff9fe8e00736baae76d70ff51325ac5bb814fe4754044b97  0.0s
 => CACHED [dev 2/5] RUN groupadd -g 1000 gouser &&     useradd -m -s /bin/bash -u 1000 -g 1000 gouser                              0.0s
 => CACHED [dev 3/5] RUN apt-get update && apt-get install -y --no-install-recommends   mariadb-client   sudo   gettext   && apt-g  0.0s
 => CACHED [dev 4/5] RUN go install -v github.com/rubenv/sql-migrate/...@v1.1.1                                                     0.0s
 => CACHED [dev 5/5] WORKDIR /workspace                                                                                             0.0s
 => [dev] exporting to image                                                                                                        0.0s
 => => exporting layers                                                                                                             0.0s
 => => writing image sha256:17797fdb3910f5f54b95371937b401a47a5be1d60783cd38c8502d7447b6018d                                        0.0s
 => => naming to docker.io/library/docker-mysql01-dev                                                                               0.0s
docker compose -f docker-compose.yml up --force-recreate -d --wait
[+] Building 0.0s (0/0)                                                                                                   docker:default
[+] Running 3/3
 ✔ Network docker-mysql01_db-network  Created                                                                                       0.1s
 ✔ Container docker-mysql01-db-1      Healthy                                                                                       0.1s
 ✔ Container docker-mysql01-dev-1     Healthy                                                                                       0.1s
docker compose -f docker-compose.yml exec dev /bin/bash -c "sql-migrate up; sql-migrate status"
Applied 2 migrations
+---------------------------------------+-------------------------------+
|               MIGRATION               |            APPLIED            |
+---------------------------------------+-------------------------------+
| 20220319203906-create-table.sql       | 2024-01-05 06:41:10 +0000 UTC |
| 20220331150041-neighborhood-sproc.sql | 2024-01-05 06:41:10 +0000 UTC |
+---------------------------------------+-------------------------------+
docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make import"
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time DATA_FILENAME=../data/P04-20-4.data <sql/loaddata.sql envsubst | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} \
--local-infile --database=geo

real    1m22.860s
user    0m0.009s
sys     0m0.013s
```

**例:** `make query`

```sh
$ make query
docker compose -f docker-compose.yml exec dev /bin/bash -c "cd query && go run cmd/main.go"
go: downloading github.com/spf13/pflag v1.0.5
go: downloading github.com/go-sql-driver/mysql v1.6.0
+-------+----------------------------------------------------------------+--------------------------------+------------+
|  ID   |                              NAME                              |            LOCATION            |  DISTANCE  |
+-------+----------------------------------------------------------------+--------------------------------+------------+
| 58869 | 皿井医院                                                       | POINT(35.68844223022461        |   7.669178 |
|       |                                                                | 139.72506713867188)            |            |
| 59325 | 四ツ谷デンタルオフィス                                         | POINT(35.68891525268555        |  60.281013 |
|       |                                                                | 139.72479248046875)            |            |
| 48441 | 医療法人　平心会　ＴｏＣＲＯＭクリニック                       | POINT(35.68808364868164        |  72.979590 |
|       |                                                                | 139.7257080078125)             |            |
| 61032 | 新宿区四谷保健センター                                         | POINT(35.68876647949219        |  74.896163 |
|       |                                                                | 139.72430419921875)            |            |
| 67384 | 柳川歯科医院                                                   | POINT(35.68768310546875        |  80.201562 |
|       |                                                                | 139.72509765625)               |            |
| 64238 | 田中内科医院                                                   | POINT(35.687992095947266       |  94.791170 |
|       |                                                                | 139.72592163085938)            |            |
| 64962 | 藤原歯科医院                                                   | POINT(35.68879699707031        | 106.837551 |
|       |                                                                | 139.7239227294922)             |            |
| 62835 | 太田医院

... snip ...

```

## DBの中を見る

dbコンテナに入って、mysqlクライアントを使って見ます。こんな感じ。

```sh
$ make mysql-client
docker compose -f docker-compose.yml exec db /bin/bash -c "LANG=C.UTF-8 mysql -q -u root -p\${MYSQL_ROOT_PASSWORD} -D geo"
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 516
Server version: 8.0.28 MySQL Community Server - GPL

Copyright (c) 2000, 2022, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> select count(*) from geo.hospital;
+-----------------------+
| count(*)              |
+-----------------------+
|                181312 |
+-----------------------+
1 row in set (0.01 sec)
```

ホストのポート3306に繋ぐと、dbコンテナに繋がります。スクショは、[Azure Data Studio](https://learn.microsoft.com/en-us/azure-data-studio/quickstart-mysql)で繋いだものです。

![azure data studio](./docs/images/ads01.pngazure-data-studio.png)

## 追加のドキュメント

- [データ投入の詳細と入れ方の比較](./docs/import.md)
- [Geoクエリーを使った近傍検索](./docs/query.md)
- [go profile toolsの利用](./docs/profile.md.md)
- 
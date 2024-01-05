# データ投入 Import

ダウンロードした、GEO JsonをMySQLにロードする。

JSON読んで、SQLを生成するまで、とSQLをMySQLに流してinsertするまでを比べると、圧倒的に後半が時間がかかるので、まずは後半を試してみる。

## MySQLへの投入

geojson2sqlコマンドでSQLを生成し、mysqlに流し込んでデータを投入する。コードは、../import にある。コマンドライン引数で、SQLの生成方法を指定する。

```sh
import/bin$ ./geojson2sql --help
  -h, --help          show help message
  -j, --json string   source geo json file (default "_data/P04-20_11_GML/small.geojson")
  -s, --sql int32     sql option, 0 is transaction, 1 is table lock, 2 is auto commit, 3 is multi value insert, 4 is tsv
  -v, --verbose       show verbose message
```

--sqlで指定できるのは５パターン

- 0 単一トランザクションで、insert文を実行. insert into hospital () values() の繰り返し。
- 1 table lock してauto commit。
- 2 auto commit。
- 3 insertの複数value.  insert into hospital () values (),(),().. と、value 指定の繰り返し。
- 4 load data localで取る込む

### ５つ目の方法の解説

これが一番速いはず。基本的なアイデアは、geojsonから、tsvファイルを作成する。そして、`temporary table` に、`load data local`で`insert`後、目的のテーブルに、`insert into`でコピーする。概ね１分７秒程度で終わる。
そのために、下記のようなSQLを用意して、`${DATA_FILENAME}` の部分をgeojsonから生成する。もしかしたら、MySQL 8なら直接geojsonが読めるかもしれないが、調べていない。

```sql
drop temporary table if exists temp_hospital;

create temporary table temp_hospital (
  id bigint auto_increment primary key,
  name varchar(500),
  latpoint float,
  lngpoint float
);

load data local infile '${DATA_FILENAME}'
into table temp_hospital (@1, @2, @3)
set
  name = @1,
  latpoint = @2,
  lngpoint = @3;

insert into
  hospital (name, location)
select
  name,
  st_srid(point(lngpoint, latpoint), 4326)
from
  temp_hospital
where
  lngpoint is not null
  and latpoint is not null
  and name is not null;
```

※サーバーのストレージにコピーすることができれば、`load data`(local無し)が使えてそれが一番速いはず。だが残念ながら、Azure MySQLではその手は使えないので、確認は後にする。

## 実行結果

```sh
$ make benchmark
bin/geojson2sql -s 0 -j ../data/P04-20.geojson > ../data/P04-20-0.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-0.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-0.txt

real    2m30.266s
user    0m2.842s
sys     0m11.392s
bin/geojson2sql -s 1 -j ../data/P04-20.geojson > ../data/P04-20-1.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-1.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-1.txt

real    19m45.170s
user    0m4.547s
sys     0m15.002s
bin/geojson2sql -s 2 -j ../data/P04-20.geojson > ../data/P04-20-2.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-2.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-2.txt

real    21m30.051s
user    0m4.718s
sys     0m15.673s
bin/geojson2sql -s 3 -j ../data/P04-20.geojson > ../data/P04-20-3.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-3.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-3.txt

real    1m35.508s
user    0m0.124s
sys     0m0.118s
bin/geojson2sql -s 4 -j ../data/P04-20.geojson > ../data/P04-20-4.data
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time DATA_FILENAME=../data/P04-20-4.data <sql/loaddata.sql envsubst | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} \
--local-infile --database=geo

real    1m43.528s
user    0m0.002s
sys     0m0.020s
```

### データロードの結論

`-4 load data localで取る込む` が一番速いが、Azure MySQLでは使えない。そうなると、`-0 単一トランザクションで、insert文を実行.` か、`-3 insertの複数value.  insert into hospital () values (),(),().. と、value 指定の繰り返し。`が妥当そう。単一トランザクションだと、トランザクションログを消費がコントロールできないので、汎用的なコードとしは、`-3` にして、適当なサイズでトランザクションが分かれるように、valueの数を調整するのが良いだろう。


```sh
$ make benchmark
docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make clean build benchmark"

... snip ....

bin/geojson2sql -s 0 -j ../data/P04-20.geojson > ../data/P04-20-0.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-0.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-0.txt

real    1m47.623s
user    0m3.176s
sys     0m13.269s
bin/geojson2sql -s 1 -j ../data/P04-20.geojson > ../data/P04-20-1.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-1.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-1.txt

real    45m6.906s
user    0m6.778s
sys     0m16.964s
bin/geojson2sql -s 2 -j ../data/P04-20.geojson > ../data/P04-20-2.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-2.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-2.txt

real    43m59.870s
user    0m6.628s
sys     0m16.546s
bin/geojson2sql -s 3 -j ../data/P04-20.geojson > ../data/P04-20-3.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-3.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-3.txt

real    1m11.895s
user    0m0.247s
sys     0m0.091s
bin/geojson2sql -s 4 -j ../data/P04-20.geojson > ../data/P04-20-4.data
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time DATA_FILENAME=../data/P04-20-4.data <sql/loaddata.sql envsubst | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} \
--local-infile --database=geo

real    1m6.972s
user    0m0.014s
sys     0m0.000s
```


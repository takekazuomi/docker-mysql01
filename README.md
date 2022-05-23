# MySQL+Docker+go

MySQLを使った開発のdockerサンプルを作り始めたから長くなったのでメモを残す。幾つかの主題が混ざって分かりづらくなったので後で整理する。

やってることは、下記のような感じ

- 開発用のubuntuとmysqlのコンテナを上げる
- 国土数値情報/医療機関データをダウンロードしてきて、mysqlに入れる
- goのコードから、近くの医療機関を検索する

## Benchemark

JSON読んで、SQLを生成するまで、とSQLをMySQLに流してinsertするまでを比べると、圧倒的に後半が時間がかかるので、まずは後半を試してみる。

### MySQLへの投入

ここで試したのは４パターン

- 0 単一トランザクションで、insert文を実行. insert into hospital () values() の繰り返し。
- 1 table lock してauto commit。
- 2 auto commit。
- 3 insertの複数value.  insert into hospital () values (),(),().. と、value 指定の繰り返し。

```sh
$ make benchmark
docker compose -f docker-compose.yml exec dev /bin/bash -c "cd dataImport && make benchmark"
rm ../data/P04-20-0.sql ../data/P04-20-1.sql ../data/P04-20-2.sql ../data/P04-20-3.sql
rm ../data/P04-20-0.txt ../data/P04-20-1.txt ../data/P04-20-2.txt ../data/P04-20-3.txt
./bin/geojson2sql -s 0 -j ../data/P04-20.geojson > ../data/P04-20-0.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-0.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-0.txt

real    1m47.371s
user    0m2.902s
sys     0m13.037s
./bin/geojson2sql -s 1 -j ../data/P04-20.geojson > ../data/P04-20-1.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-1.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-1.txt

real    45m29.341s
user    0m7.345s
sys     0m16.694s
./bin/geojson2sql -s 2 -j ../data/P04-20.geojson > ../data/P04-20-2.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-2.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-2.txt

real    43m47.778s
user    0m7.280s
sys     0m16.113s
./bin/geojson2sql -s 3 -j ../data/P04-20.geojson > ../data/P04-20-3.sql
echo "truncate table hospital;" | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo
time cat ../data/P04-20-3.sql | mysql -h ${MYSQL_HOST} -u ${MYSQL_USER} --password=${MYSQL_PASSWORD} --database=geo | tee ../data/P04-20-3.txt

real    1m19.526s
user    0m0.167s
sys     0m0.167s
```

### SQLの生成

まずは、最初のコードをBenchmarkして、memprofile を見る。コード的には単純なので、まずはメモリから。

```sh
$ go test -benchmem -run='^$' -bench '^BenchmarkMain$' -memprofile tmp/memprofile.out -cpuprofile tmp/profile.out github.com/takekazuomi/docker-mysql01/import/cmd -v
goos: linux
goarch: amd64
pkg: github.com/takekazuomi/docker-mysql01/import/cmd
cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
BenchmarkMain
BenchmarkMain/P04-20.geojson
BenchmarkMain/P04-20.geojson-8                 1        1182726700 ns/op        557567760 B/op   2551299 allocs/op
BenchmarkMain/P04-20.geojson#01
BenchmarkMain/P04-20.geojson#01-8              1        1207487000 ns/op        557551016 B/op   2550972 allocs/op
BenchmarkMain/P04-20.geojson#02
BenchmarkMain/P04-20.geojson#02-8              1        1186791800 ns/op        557546168 B/op   2550968 allocs/op
BenchmarkMain/P04-20.geojson#03
BenchmarkMain/P04-20.geojson#03-8              1        1165761500 ns/op        573048088 B/op   2732283 allocs/op
```

```
$ go tool pprof tmp/memprofile.out

(pprof) top
Showing nodes accounting for 2128.31MB, 98.93% of 2151.36MB total
Dropped 38 nodes (cum <= 10.76MB)
Showing top 10 nodes out of 18
      flat  flat%   sum%        cum   cum%
 1500.29MB 69.74% 69.74%  1500.29MB 69.74%  io.ReadAll
  444.52MB 20.66% 90.40%   444.52MB 20.66%  reflect.unsafe_NewArray
  129.01MB  6.00% 96.40%   129.01MB  6.00%  encoding/json.(*decodeState).literalStore
      33MB  1.53% 97.93%       33MB  1.53%  github.com/takekazuomi/docker-mysql01/import/geojson.(*Features).printSQL
   19.50MB  0.91% 98.84%   464.02MB 21.57%  reflect.MakeSlice
       2MB 0.093% 98.93%    18.50MB  0.86%  github.com/takekazuomi/docker-mysql01/import/geojson.(*Features).printTsv
         0     0% 98.93%   593.02MB 27.57%  encoding/json.(*decodeState).array
         0     0% 98.93%   593.02MB 27.57%  encoding/json.(*decodeState).object
         0     0% 98.93%   593.02MB 27.57%  encoding/json.(*decodeState).unmarshal
         0     0% 98.93%   593.02MB 27.57%  encoding/json.(*decodeState).value
```
![ReadAll](./images/profile001.png)

ファイルサイズが、67MのJSONを４回読んでるが、それでio.ReadAllが1.5Gも使っている。reflect.unsafe_NewArray もガッツリメモリ使っているので、おそらく読んでる途中で。バッファーのリアロケーションを繰り返しているのだと予想される。死にそうだ。

### ReadFile版

ぐぐったら、`ioutil.ReadFile` が**まだまし**と書いてあったので、`ReadFile`にしてみた。どうやら、`ioutil.ReadAll`は、Readerから読む（＝サイズがわからない）。ReadFileはファイルから読む（= ファイルなら読む前にサイズがわかる）ということらしい。当たり前である。

```
$ go test -benchmem -run='^$' -bench '^BenchmarkMain$' -memprofile tmp/memprofile.out -cpuprofile tmp/profile.out github.com/takekazuomi/docker-mysql01/import/cmd -v
goos: linux
goarch: amd64
pkg: github.com/takekazuomi/docker-mysql01/import/cmd
cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
BenchmarkMain
BenchmarkMain/P04-20.geojson
BenchmarkMain/P04-20.geojson-8                 1        1182859700 ns/op        233744312 B/op   2551262 allocs/op
BenchmarkMain/P04-20.geojson#01
BenchmarkMain/P04-20.geojson#01-8              1        1141279700 ns/op        233720560 B/op   2550917 allocs/op
BenchmarkMain/P04-20.geojson#02
BenchmarkMain/P04-20.geojson#02-8              1        1119044800 ns/op        233717488 B/op   2550916 allocs/op
BenchmarkMain/P04-20.geojson#03
BenchmarkMain/P04-20.geojson#03-8              1        1115623000 ns/op        249203696 B/op   2732234 allocs/op
PASS
ok      github.com/takekazuomi/docker-mysql01/import/cmd        4.732s

$ go tool pprof tmp/memprofile.out
File: cmd.test
Type: alloc_space
Time: Apr 24, 2022 at 2:43pm (JST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 910.39MB, 99.38% of 916.07MB total
Dropped 29 nodes (cum <= 4.58MB)
Showing top 10 nodes out of 22
      flat  flat%   sum%        cum   cum%
  442.73MB 48.33% 48.33%   442.73MB 48.33%  reflect.unsafe_NewArray
  265.66MB 29.00% 77.33%   265.66MB 29.00%  os.ReadFile
  132.01MB 14.41% 91.74%   132.01MB 14.41%  encoding/json.(*decodeState).literalStore
   30.50MB  3.33% 95.07%    30.50MB  3.33%  github.com/takekazuomi/docker-mysql01/import/geojson.(*Features).printSQL
   14.50MB  1.58% 96.65%   457.23MB 49.91%  reflect.MakeSlice
   13.50MB  1.47% 98.13%    13.50MB  1.47%  strconv.FormatFloat
       8MB  0.87% 99.00%        8MB  0.87%  strings.(*Builder).grow (inline)
    3.50MB  0.38% 99.38%       25MB  2.73%  github.com/takekazuomi/docker-mysql01/import/geojson.(*Features).printTsv
         0     0% 99.38%   589.23MB 64.32%  encoding/json.(*decodeState).array
         0     0% 99.38%   589.23MB 64.32%  encoding/json.(*decodeState).object
(pprof) png
Generating report in profile001.png
```
![ReadFile](./images/profile002.png)

かかる時間は変わらないけど、メモリが半分になった。

## MySQLでload data

追加で、５つ目の方法を試す。これが一番速いはず。基本的なアイデアは、geojsonから、tsvファイルを作成する。そして、`temporary table` に、`load data local`で`insert`後、目的のテーブルに、`insert into`でコピーする。概ね１分７秒程度で終わる。
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

今までの全パターンを流すと下記のようになる。

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

## Query

近くを探すクエリーを書くとこんな感じになる。今回は、ストアド(neighborhood)にしてある。

簡単にクエリを解説する。このクエリーは２段階に分かれてる。まず、最初に指定された座標 (latpoint, lngpoint)を含む最小境界矩形(mbr)を作成し、中に入るものを抽出する(最小境界矩形内判定、bb)。その後、抽出結果(bb)から指定距離内に入る対象に絞る。
最小境界矩形内の判定は、指定座標から距離r離れた矩形の対角線をlinestring()で引き、それをmbrcontains()にわたすことで算出している。ここで、spatial index が使われる。
最初の抽出クエリ部分は、地球が（真球）と仮定して、地球の赤道周長から算出したものを１度の距離(@units)として使っている。unitの111のマジックナンバーは、40075.017/360から計算したものだ。

### 参考

- https://dev.mysql.com/blog-archive/why-boost-geometry-in-mysql/
  - MySQL 5.7からのGIS関連機能は、Boost.Geometry 由来だそうだ
- https://dev.mysql.com/doc/refman/8.0/ja/spatial-relation-functions-mbr.html
- https://ja.wikipedia.org/wiki/%E7%A9%BA%E9%96%93%E3%82%A4%E3%83%B3%E3%83%87%E3%83%83%E3%82%AF%E3%82%B9
  - ちなみにMySQLのspatial indexでは、R木が使われている
  - https://dev.mysql.com/blog-archive/innodb-spatial-indexes-in-5-7-4-lab-release/
- GeoHashを使う方法もある
  - https://en.wikipedia.org/wiki/Geohash
  - https://dev.mysql.com/doc/refman/5.7/en/spatial-geohash-functions.html

```sql
set @latpoint = 35.6884204226699;
set @lngpoint=139.72515649841105;
set @r = 1;
set @units = 111.0;

with bb (id, name, location) as (
  select
    id,
    name,
    location
  from
    hospital
  where
    mbrcontains(
      st_geomfromtext(
        concat(
          'linestring(',
          @latpoint - (@r / @units),
          ' ',
          @lngpoint - (@r / (@units * cos(radians(@latpoint)))),
          ',',
          @latpoint + (@r / @units),
          ' ',
          @lngpoint + (@r / (@units * cos(radians(@latpoint)))),
          ')'
        ),
        4326
      ),
      location
    )
),
target (id, name, location, distance) AS (
  select
    id,
    name,
    st_astext(location) location,
    st_distance_sphere(
      location,
      st_geomfromtext(
        concat('point(', @latpoint, ' ', @lngpoint, ')'),
        4326
      )
    ) st_distance_sphere
  from
    bb
)
select
  *
from
  target
where
  distance <= @r * 1000
order by
  distance asc;
```

### 実行計画

explain analyze

```text
-> Sort: st_distance_sphere(hospital.location,<cache>(st_geomfromtext(concat('point(',(@latpoint),' ',(@lngpoint),')'),4326)))  (cost=0.71 rows=1) (actual time=9.239..9.381 rows=178 loops=1)
    -> Filter: ((st_distance_sphere(hospital.location,<cache>(st_geomfromtext(concat('point(',(@latpoint),' ',(@lngpoint),')'),4326))) <= <cache>(((@r) * 1000))) and mbrcontains(<cache>(st_geomfromtext(concat('linestring(',((@latpoint) - ((@r) / (@units))),' ',((@lngpoint) - ((@r) / ((@units) * cos(radians((@latpoint)))))),',',((@latpoint) + ((@r) / (@units))),' ',((@lngpoint) + ((@r) / ((@units) * cos(radians((@latpoint)))))),')'),4326)),hospital.location))  (cost=0.71 rows=1) (actual time=0.465..7.192 rows=178 loops=1)
        -> Index range scan on hospital using idx_hospital_location over (location unprintable_geometry_value)  (cost=0.71 rows=1) (actual time=0.303..2.163 rows=221 loops=1)
```

explain

| id | select\_type | table | partitions | type | possible\_keys | key | key\_len | ref | rows | filtered | Extra |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | SIMPLE | hospital | NULL | range | idx\_hospital\_location | idx\_hospital\_location | 34 | NULL | 1 | 100 | Using where; Using filesort |

### Index無しで実行

本来は、Geometry型を使わない場合と比較したいのだが、少し手間がかかるので簡易的にindexを削除して比較。当然フルスキャンになる。`(cost=17027.75 rows=167790)` の部分に注目。全行が対象になる。`actual time=1385.490` も大きく上がる。

explain analyze

```
-> Sort: st_distance_sphere(hospital.location,<cache>(st_srid(point((@lngpoint),(@latpoint)),4326)))  (cost=17027.75 rows=167790) (actual time=1385.490..1385.604 rows=178 loops=1)
    -> Filter: ((st_distance_sphere(hospital.location,<cache>(st_srid(point((@lngpoint),(@latpoint)),4326))) <= <cache>(((@r) * 1000))) and mbrcontains(<cache>(st_geomfromtext(concat('linestring(',((@latpoint) - ((@r) / (@units))),' ',((@lngpoint) - ((@r) / ((@units) * cos(radians((@latpoint)))))),',',((@latpoint) + ((@r) / (@units))),' ',((@lngpoint) + ((@r) / ((@units) * cos(radians((@latpoint)))))),')'),4326)),hospital.location))  (cost=17027.75 rows=167790) (actual time=320.947..1384.123 rows=178 loops=1)
        -> Table scan on hospital  (cost=17027.75 rows=167790) (actual time=0.028..225.674 rows=181316 loops=1)
```

explain

| id | select\_type | table | partitions | type | possible\_keys | key | key\_len | ref | rows | filtered | Extra |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | SIMPLE | hospital | NULL | ALL | NULL | NULL | NULL | NULL | 167790 | 100 | Using where; Using filesort |

## json stream 版

全然はやくならない、メモリーも使いすぎ。これはだめそう。

## TODO

- DBへの繋ぎ方
  - https://hub.docker.com/_/mysql
- [ ] go/sql版
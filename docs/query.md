# Query

MySQLのspatial indexを使って、近傍検索をする

近くを探すクエリーを書くとこんな感じになる。今回は、ストアド(neighborhood)にしてある。

簡単にクエリを解説する。このクエリーは２段階に分かれてる。まず、最初に指定された座標 (latpoint, lngpoint)を含む最小境界矩形(mbr)を作成し、中に入るものを抽出する(最小境界矩形内判定、bb)。その後、抽出結果(bb)から指定距離内に入る対象に絞る。
最小境界矩形内の判定は、指定座標から距離r離れた矩形の対角線をlinestring()で引き、それをmbrcontains()にわたすことで算出している。ここで、spatial index が使われる。
最初の抽出クエリ部分は、地球が（真球）と仮定して、地球の赤道周長から算出したものを１度の距離(@units)として使っている。unitの111のマジックナンバーは、40075.017/360から計算したものだ。

## 参考

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

## 実行計画

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

## Index無しで実行

本来は、Geometry型を使わない場合と比較したいのだが、少し手間がかかるので簡易的にindexを削除して比較。当然フルスキャンになる。`(cost=17027.75 rows=167790)` の部分に注目。全行が対象になる。`actual time=1385.490` も大きく上がる。

explain analyze

```sql
-> Sort: st_distance_sphere(hospital.location,<cache>(st_srid(point((@lngpoint),(@latpoint)),4326)))  (cost=17027.75 rows=167790) (actual time=1385.490..1385.604 rows=178 loops=1)
    -> Filter: ((st_distance_sphere(hospital.location,<cache>(st_srid(point((@lngpoint),(@latpoint)),4326))) <= <cache>(((@r) * 1000))) and mbrcontains(<cache>(st_geomfromtext(concat('linestring(',((@latpoint) - ((@r) / (@units))),' ',((@lngpoint) - ((@r) / ((@units) * cos(radians((@latpoint)))))),',',((@latpoint) + ((@r) / (@units))),' ',((@lngpoint) + ((@r) / ((@units) * cos(radians((@latpoint)))))),')'),4326)),hospital.location))  (cost=17027.75 rows=167790) (actual time=320.947..1384.123 rows=178 loops=1)
        -> Table scan on hospital  (cost=17027.75 rows=167790) (actual time=0.028..225.674 rows=181316 loops=1)
```

explain

| id | select\_type | table | partitions | type | possible\_keys | key | key\_len | ref | rows | filtered | Extra |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | SIMPLE | hospital | NULL | ALL | NULL | NULL | NULL | NULL | 167790 | 100 | Using where; Using filesort |

## 結論

spatial indexは便利、使った方が良い。だた、古いMySQLだと怪しいのでバージョンには注意。


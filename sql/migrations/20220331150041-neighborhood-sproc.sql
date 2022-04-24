-- +migrate Up

-- +migrate StatementBegin

-- クエリーは２段階に分かれています。まず、最初に(bb)、指定された座標
-- (latpoint, lngpoint)を含む最小境界矩形を作成し、中に入るものを抽出し(bb)。
-- その後、抽出結果から指定距離内に入る対象にしぼります。
-- 最小境界矩形は、指定座標から距離r離れた矩形の対角線をlinestring()で引き。
-- それをmbrcontains()にわたすことで算出しています。
-- 111のマジックナンバーは、地球の赤道周長	40075.017/360からの概算
-- https://dev.mysql.com/doc/refman/8.0/ja/spatial-relation-functions-mbr.html

create procedure neighborhood(
  -- 緯度
  in latpoint float,
  -- 経度
  in lngpoint float,
  -- 距離(km)
  in r float
) begin

-- 111 statute km per degree
set @units = 111.0;

-- explain analyze
-- bounding box query
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
          latpoint - (r / @units),
          ' ',
          lngpoint - (r / (@units * cos(radians(latpoint)))),
          ',',
          latpoint + (r / @units),
          ' ',
          lngpoint + (r / (@units * cos(radians(latpoint)))),
          ')'
        ),
        4326
      ),
      location
    )
),
-- exact neighborhood
target (id, name, location, distance) AS (
  select
    id,
    name,
    st_astext(location) location,
    st_distance_sphere(
      location,
      st_geomfromtext(
        concat('point(', latpoint, ' ', lngpoint, ')'),
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
  distance <= r * 1000
order by
  distance asc;
-- limit 100
end;
-- +migrate StatementEnd

-- +migrate Down
drop procedure if exists neighborhood;

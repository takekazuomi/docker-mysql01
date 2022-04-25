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

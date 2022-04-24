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

-- truncate table hospital;

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
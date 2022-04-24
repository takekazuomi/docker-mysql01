select
    id,
    st_astext(location)
from
    hospital

select
    id,
    st_astext(location) location,
    st_distance_sphere(
        location,
        st_geomfromtext(
            'point(35.78644425007628, 139.6275956501499)',
            4326
        )
    ) d
from
    hospital


WITH cte (id, name, location, distance) AS
         (
             select id,
                    name,
                    st_astext(location) location,
                    st_distance_sphere(
                            location,
                            st_geomfromtext(
                                    'point(35.78644425007628 139.6275956501499)',
                                    4326
                                )
                        ) st_distance_sphere
             from hospital
         )
select * from cte where distance < 1000 order by  distance asc limit 1

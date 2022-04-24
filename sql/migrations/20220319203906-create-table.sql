
-- +migrate Up
create table hospital (
    id bigint auto_increment,
    name varchar(500) not null,
    location point not null srid 4326,
    constraint hospital_pk primary key (id)
) comment '病院';

create spatial index idx_hospital_location ON hospital(location);

-- +migrate Down
drop table hospital;

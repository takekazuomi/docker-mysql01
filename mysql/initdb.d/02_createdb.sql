create schema geo;
create user 'geouser'@'%' identified with mysql_native_password BY 'mysql';
grant all privileges on geo.* to 'geouser'@'%';
flush privileges;

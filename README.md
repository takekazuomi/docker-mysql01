https://hub.docker.com/_/mysql


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
$

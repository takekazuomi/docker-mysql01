
https://hub.docker.com/_/mysql

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

```
$ go test -benchmem -run='^$' -bench '^BenchmarkMain$' github.com/takekazuomi/docker-mysql01/import/cmd -v
goos: linux
goarch: amd64
pkg: github.com/takekazuomi/docker-mysql01/import/cmd
cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
BenchmarkMain
BenchmarkMain/P04-20.geojson
BenchmarkMain/P04-20.geojson-8          1000000000               0.0000749 ns/op               0 B/op          0 allocs/op
BenchmarkMain/P04-20.geojson#01
BenchmarkMain/P04-20.geojson#01-8       1000000000               0.0000829 ns/op               0 B/op          0 allocs/op
BenchmarkMain/P04-20.geojson#02
BenchmarkMain/P04-20.geojson#02-8              1        1175299800 ns/op        557521072 B/op   2550968 allocs/op
BenchmarkMain/P04-20.geojson#03
BenchmarkMain/P04-20.geojson#03-8              1        1152912700 ns/op        557521112 B/op   2550967 allocs/op
BenchmarkMain/P04-20.geojson#04
BenchmarkMain/P04-20.geojson#04-8              1        1160774400 ns/op        557518904 B/op   2550959 allocs/op
BenchmarkMain/P04-20.geojson#05
BenchmarkMain/P04-20.geojson#05-8              1        1142539700 ns/op        573023280 B/op   2732279 allocs/op
PASS
ok      github.com/takekazuomi/docker-mysql01/import/cmd        4.657s
```

## TODO

- [ ] LOAD DATA版
- [ ] go/sql版
SHELL=bash

GEOJSON_BIN	= bin/geojson2sql
SQL_OPTS 	= ../data/P04-20-0.sql ../data/P04-20-1.sql ../data/P04-20-2.sql ../data/P04-20-3.sql
BENCHMARK_OPTS 	= ../data/P04-20-0.txt ../data/P04-20-1.txt ../data/P04-20-2.txt ../data/P04-20-3.txt

build:
	go build -o $(GEOJSON_BIN) cmd/main.go

tidy:
	go mod tidy

$(SQL_OPTS): ../data/P04-20-%.sql: ../data/P04-20.geojson
	$(GEOJSON_BIN) -s $* -j $< > $@

$(BENCHMARK_OPTS): ../data/P04-20-%.txt: ../data/P04-20-%.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo | tee $@

sql: $(SQL_OPTS) $(GEOJSON_BIN)

import0: ../data/P04-20-0.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo

import3: ../data/P04-20-3.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo

import: import3

benchmark: clean $(BENCHMARK_OPTS)

clean:
	-rm $(GEOJSON_BIN)
	-rm $(SQL_OPTS)
	-rm $(BENCHMARK_OPTS)

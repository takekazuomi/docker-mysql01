SHELL=bash

build: bin/geojson2sql

bin/geojson2sql:
	go build -o bin/geojson2sql cmd/main.go

tidy:
	go mod tidy

SQL_OPTS = ../data/P04-20-0.sql ../data/P04-20-1.sql ../data/P04-20-2.sql ../data/P04-20-3.sql

$(SQL_OPTS): ../data/P04-20-%.sql: ../data/P04-20.geojson
	./bin/geojson2sql -s $* -j $< > $@

BENCHMARK_OPTS = ../data/P04-20-0.txt ../data/P04-20-1.txt ../data/P04-20-2.txt ../data/P04-20-3.txt

$(BENCHMARK_OPTS): ../data/P04-20-%.txt: ../data/P04-20-%.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo | tee $@

sql: $(SQL_OPTS)

import0: ../data/P04-20-0.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo

import3: ../data/P04-20-3.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo

import: import3

benchmark: clean $(BENCHMARK_OPTS)

clean:
	-rm $(SQL_OPTS)
	-rm $(BENCHMARK_OPTS)

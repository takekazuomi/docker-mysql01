SHELL=bash

GEOJSON_BIN	= bin/geojson2sql
SQL_OPTS 	= ../data/P04-20-0.sql ../data/P04-20-1.sql ../data/P04-20-2.sql ../data/P04-20-3.sql
DATA_OPTS 	= ../data/P04-20-4.data ../data/P04-20-5.data
BENCHMARK_OPTS 	= ../data/P04-20-0.txt ../data/P04-20-1.txt ../data/P04-20-2.txt ../data/P04-20-3.txt

build:
	go build -o $(GEOJSON_BIN) cmd/main.go

tidy:
	go mod tidy

$(SQL_OPTS): ../data/P04-20-%.sql: ../data/P04-20.geojson
	$(GEOJSON_BIN) -s $* -j $< > $@

$(DATA_OPTS): ../data/P04-20-%.data: ../data/P04-20.geojson
	$(GEOJSON_BIN) -s $* -j $< > $@

$(BENCHMARK_OPTS): ../data/P04-20-%.txt: ../data/P04-20-%.sql
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time cat $< | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo | tee $@

sql: $(SQL_OPTS) $(DATA_OPTS) $(GEOJSON_BIN)

import4: ../data/P04-20-4.data
	echo "truncate table hospital;" | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} --database=geo
	time DATA_FILENAME=$< <sql/loaddata.sql envsubst | mysql -h $${MYSQL_HOST} -u $${MYSQL_USER} --password=$${MYSQL_PASSWORD} \
	--local-infile --database=geo

import: $(DATA_OPTS)

benchmark: $(BENCHMARK_OPTS) $(DATA_OPTS)

clean:
	-rm $(GEOJSON_BIN)
	-rm $(SQL_OPTS)
	-rm $(DATA_OPTS)
	-rm $(BENCHMARK_OPTS)

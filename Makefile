export UID_GID=$(shell id -u):$(shell id -g)

.PHONY: query

help:		## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

bootstrap:	## bootstrap database
bootstrap: build up sql-migrate-up import

benchmark:	## Run benchmark
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make clean build benchmark"

query:		## Run query
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd query && go run cmd/main.go"

up: 		## Up
	docker compose -f docker-compose.yml up --force-recreate -d --wait

down: 		## Down
	docker compose -f docker-compose.yml down

build: 		## Build
	docker compose -f docker-compose.yml build

logs:		## Show logs
	docker compose -f docker-compose.yml logs -f

login-dev:	## login dev container
	docker compose -f docker-compose.yml exec dev /bin/bash

login-db:	## login db container
	docker compose -f docker-compose.yml exec db /bin/bash

clean: down
	rm -rf mysql/data/*

mysql-client:	## connect db from mysql client
	docker compose -f docker-compose.yml exec db /bin/bash -c "LANG=C.UTF-8 mysql -q -u root -p\$${MYSQL_ROOT_PASSWORD} -D geo"

data/P04-20.geojson:
	./import/p04download.sh

sql-migrate-up:
	docker compose -f docker-compose.yml exec dev /bin/bash -c "sql-migrate up; sql-migrate status"

sql-migrate-new:
	docker compose -f docker-compose.yml exec dev /bin/bash -c "sql-migrate new tempolaryname"

import/bin/geojson2sql:
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make build"

import: data/P04-20.geojson import/bin/geojson2sql sql-migrate-up
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make import"

sql: data/P04-20.geojson import/bin/geojson2sql
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make sql"




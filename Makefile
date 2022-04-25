export UID_GID=$(shell id -u):$(shell id -g)

help:		## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

up: 		## Up
	docker compose -f docker-compose.yml up --force-recreate -d

down: 		## Down
	docker compose -f docker-compose.yml down

build: 		## Build
	docker compose -f docker-compose.yml build

logs:		## Show logs
	docker compose -f docker-compose.yml logs -f

login-dev:
	docker compose -f docker-compose.yml exec dev /bin/bash

login-db:
	docker compose -f docker-compose.yml exec db /bin/bash

clean: down
	rm -rf mysql/data/*

mysql-client:
	docker compose -f docker-compose.yml exec db /bin/bash -c "mysql -u root -p -D db"

data/P04-20.geojson:
	./datagen/p04download.sh

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

benchmark:
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd import && make clean build benchmark"

.PHONY: query
query:
	docker compose -f docker-compose.yml exec dev /bin/bash -c "cd query && go run cmd/main.go"




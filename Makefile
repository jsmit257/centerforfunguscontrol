SHELL := /bin/bash

.PHONY: unit
unit:
	go test -cover ./internal/...

.PHONY: postgres
postgres:
	docker-compose up -d postgres

.PHONY: inspect
inspect:
	docker-compose exec -it postgres psql -Upostgres huautla

.PHONY: run-local
run-local:
	(go run ./ingress/http/... >log.json 2>&1 & k=$!; tail -f log.json | jq -a .; kill $k)

.PHONY: run-docker
run-docker:
	docker-compose up --build --remove-orphans -d run-docker

run-web:
	docker-compose down -t5 --remove-orphans cffc-web
	docker-compose up --build --remove-orphans -d cffc-web

.PHONY: tests
tests: down unit
	sudo rm -v ./testalbum/*
	docker-compose up --build --remove-orphans system-test
	docker tag jsmit257/cffc:latest jsmit257/cffc:lkg

.PHONY: public
public: tests
	docker-compose up -d cffc-web

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: down
down:
	docker-compose down -t5 --remove-orphans

.PHONY: deploy
deploy: # no hard dependency on `tests/public/etc` for mow
	docker-compose build run-docker
	docker tag jsmit257/cffc:latest jsmit257/cffc:lkg

.PHONY: push
push: # just docker, not git
	docker push jsmit257/cffc:lkg

.PHONY: push-all
push-all: push
	docker push jsmit257/huautla:lkg
	docker push jsmit257/us-db-mysql-test:lkg
	docker push jsmit257/us-srv-mysql:lkg
	docker push jsmit257/cffc-web:lkg
	
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
	docker-compose up --build --remove-orphans cffc-web

.PHONY: tests
tests: docker-down unit
	docker-compose up --build --remove-orphans system-test
	docker tag jsmit257/cffc:latest jsmit257/cffc:lkg

.PHONY: public
public: tests
	docker-compose up -d cffc-web

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

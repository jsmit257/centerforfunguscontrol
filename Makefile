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

.PHONY: tests
tests: docker-down
# tests: docker-down unit
	docker-compose up --build --remove-orphans system-test
	docker tag cffc:latest jsmit257/cffc:lkg

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

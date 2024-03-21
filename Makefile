SHELL := /bin/bash

.PHONY: unit
unit:
	go test -cover ./internal/...

.PHONY: postgres
postgres:
	docker-compose up -d postgres

.PHONY: run-local
run-local:
	go run ./ingress/http/... 2>&1 | tee log.json | jq .

.PHONY: run-docker
run-docker:
	docker-compose up --build --remove-orphans -d run-docker

.PHONY: system-test
tests: unit
	docker-compose up --remove-orphans system-test
	docker tag cffc:latest jsmit257/cffc:lkg

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

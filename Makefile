SHELL := /bin/bash

.PHONY: unit
unit:
	go test -cover ./. ./internal/...

.PHONY: postgres
postgres:
	docker-compose up -d postgres
	sleep 2s

.PHONY: install
install: postgres
	docker-compose up --build --force-recreate install
	make docker-down

.PHONY: system-test
system-test: postgres
	# docker-compose up --build --force-recreate install-system-test
	# -POSTGRES_SSLMODE="${POSTGRES_SSLMODE}" docker-compose up --build --force-recreate system-test
	-docker-compose up --build --force-recreate system-test
	make docker-down

.PHONY: run-local
run-local:
go run ./ingress/http/... 2>&1 | tee log.json | jq .

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

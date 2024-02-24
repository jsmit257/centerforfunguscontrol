SHELL := /bin/bash

.PHONY: build
build:
	go build -o system-test -v ./...

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
	docker-compose up --build --force-recreate install-system-test
	# -POSTGRES_SSLMODE="${POSTGRES_SSLMODE}" docker-compose up --build --force-recreate system-test
	-docker-compose up --build --force-recreate system-test
	make docker-down

# .PHONY: package-serve-mysql
# package-serve-mysql: compile-serve-mysql
# 	docker-compose up --build --force-recreate package-serve-mysql

# .PHONY: test-serve-mysql
# test-serve-mysql: docker-down build
# 	docker-compose up --build --force-recreate test-serve-mysql

# .PHONY: system-test
# system-test: docker-down package-serve-mysql
# 	docker-compose up --build --force-recreate schema
# 	docker-compose up serve-mysql &
# 	sleep 2s
# 	-cd ./tests/system; go test -v ./user/...
# 	-curl localhost:3000/metrics

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

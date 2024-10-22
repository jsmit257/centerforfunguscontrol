SHELL := /bin/bash

.PHONY: build
build:
	go build -o system-test -v ./...

.PHONY: unit
unit:
	go test -cover ./. ./types/... ./internal/...

.PHONY: tag-dockerfile
tag-dockerfile:
	docker tag huautla:latest huautla:lkg
	
.PHONY: postgres
postgres:
	docker-compose up --build --remove-orphans -d postgres

.PHONY: inspect
inspect:
	docker-compose exec -it postgres psql -Upostgres huautla

.PHONY: install-system-test
install-system-test: postgres
	sleep 2s
	docker-compose exec \
		-ePOSTGRES_HOST=huautla \
		-ePOSTGRES_PORT=5432 \
		-ePOSTGRES_USER=postgres \
		-ePOSTGRES_PASSWORD=root \
		postgres \
		/bin/sh -c "cd /huautla && ./bin/install-system-test.sh"
	docker tag huautla:latest jsmit257/huautla:lkg
	# make docker-down

.PHONY: system-test
system-test: docker-down unit install-system-test
	docker-compose up system-test #>out 2>&1
	#cat out
	# docker push jsmit257/huautla:lkg
	# make docker-down

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

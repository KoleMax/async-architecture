ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
else
	include .env.dist
	export $(shell sed 's/=.*//' .env.dist)
endif

######################################## LOCAL BINARIES ########################################

LOCAL_BIN := $(CURDIR)/bin

GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG := 1.41.0

GOOSE_BIN := $(LOCAL_BIN)/goose
GOOSE_TAG := 2.6.0
MIGRATIONS_DIR := $(CURDIR)/migrations
AUTH_MIGRATIONS_DIR := $(CURDIR)/migrations-accounts

ENVSUBST_BIN := $(LOCAL_BIN)/envsubst
ENVSUBST_TAG := 1.2.0

GOIMPORTS_BIN := $(LOCAL_BIN)/goimports
GOIMPORTS_TAG := 0.1.8

SWAG_BIN := $(LOCAL_BIN)/swag
SWAG_TAG := 1.7.8

install-bin-deps: .install-lint .install-goose .install-envsubst .install-goimports .install-swag

.install-lint:
	$(info Installing golangci-lint v$(GOLANGCI_TAG))
	tmp_dir=$$(mktemp -d) && \
	cd $$tmp_dir && \
	go mod init tmp && \
	go get -d github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG) && \
	go build -o $(GOLANGCI_BIN) github.com/golangci/golangci-lint/cmd/golangci-lint && \
	rm -rf $$tmp_dir

.install-envsubst:
	$(info Installing envsubst v$(ENVSUBST_TAG))
	tmp_dir=$$(mktemp -d) && \
	cd $$tmp_dir && \
	go mod init tmp && \
	go get -d github.com/a8m/envsubst/cmd/envsubst@v$(ENVSUBST_TAG) && \
	go build -o $(ENVSUBST_BIN) github.com/a8m/envsubst/cmd/envsubst && \
	rm -rf $$tmp_dir

.install-goose:
	$(info Installing goose v$(GOLANGCI_TAG) migration tool)
	tmp_dir=$$(mktemp -d) && \
	cd $$tmp_dir && \
	go mod init tmp && \
	go get -d github.com/pressly/goose/cmd/goose@v$(GOOSE_TAG) && \
	go build -o $(GOOSE_BIN) github.com/pressly/goose/cmd/goose && \
	rm -rf $$tmp_dir

.install-goimports:
	$(info Installing goimports v$(GOIMPORTS_TAG))
	tmp_dir=$$(mktemp -d) && \
	cd $$tmp_dir && \
	go mod init tmp && \
	go get -d golang.org/x/tools/cmd/goimports@v$(GOIMPORTS_TAG) && \
	go build -o $(GOIMPORTS_BIN) golang.org/x/tools/cmd/goimports && \
	rm -rf $$tmp_dir

.install-swag:
	$(info Installing swag v$(SWAG_TAG))
	tmp_dir=$$(mktemp -d) && \
	cd $$tmp_dir && \
	go mod init tmp && \
	go get -d github.com/swaggo/swag/cmd/swag@v$(SWAG_TAG) && \
	go build -o $(SWAG_BIN) github.com/swaggo/swag/cmd/swag && \
	rm -rf $$tmp_dir

######################################## MIGRATIONS ########################################

# Get migration name for create-migrations command
ifeq (create-migration,$(firstword $(MAKECMDGOALS)))
  MIGRATION_NAME := $(word 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(MIGRATION_NAME):;@:)
endif

create-migration:
	$(info Create new migration file)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) create $(MIGRATION_NAME) sql

upgrade-db:
	$(info Upgrade db migrations version)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable" up

downgrade-db:
	$(info Rollback db migrations)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable" down

######################################## MIGRATIONS ########################################

######################################## GEN ########################################

generate:
	$(SWAG_BIN) init -g ./cmd/api/main.go

######################################## GEN ########################################

######################################## LINTERS ########################################

lint: lint-swag
	CGO_ENABLED=1 $(GOLANGCI_BIN) --timeout 240s run --config=./.golangci.yml ./...

lint-swag: .build-SFotaUtils
	CGO_ENABLED=1 $(SWAG_BIN) fmt -g ./cmd/api/main.go

lint-fix: .build-SFotaUtils
	CGO_ENABLED=1 $(GOLANGCI_BIN) run --config=./.golangci.yml --fix ./...

######################################## LINTERS ########################################

build:
	$(info Building executables)
	CGO_ENABLED=1 go build -o $(LOCAL_BIN) ./cmd/api

build-auth:
	$(info Building executables)
	CGO_ENABLED=1 go build -o $(LOCAL_BIN) ./cmd/auth

generate-config:
	$(ENVSUBST_BIN) < config.template.yml > bin/config.yml

######################################## TESTING ########################################

.TEST_POSTGRES_DB_NAME=test
.TEST_POSTGRES_HOST=localhost
.TEST_POSTGRES_PORT=15432
.TEST_POSTGRES_USER=test
.TEST_POSTGRES_PASSWORD=test
.TEST_POSTGRES_CONTAINER_NAME=api-db

.TEST_AUTH_POSTGRES_DB_NAME=auth-test
.TEST_AUTH_POSTGRES_HOST=localhost
.TEST_AUTH_POSTGRES_PORT=15433
.TEST_AUTH_POSTGRES_USER=auth-test
.TEST_AUTH_POSTGRES_PASSWORD=auth-test
.TEST_AUTH_POSTGRES_CONTAINER_NAME=auth-db

test-api: .test-api-environment-up
	CGO_ENABLED=1 go test -v \
	./internal/... -test-db-dsn=postgres://${.TEST_POSTGRES_USER}:${.TEST_POSTGRES_PASSWORD}@${.TEST_POSTGRES_HOST}:${.TEST_POSTGRES_PORT}/${.TEST_POSTGRES_DB_NAME}?sslmode=disable \
		-test-storage-endpoint=${.TEST_STORAGE_HOST}:${.TEST_STORAGE_PORT} \
		-test-storage-bucket=${.TEST_STORAGE_BUCKET} \
		-test-storage-access-key=${.TEST_STORAGE_ACCESS_KEY} \
		-test-storage-secret-key=${.TEST_STORAGE_SECRET_KEY}

.test-api-environment-up:
ifeq ("$(shell docker ps -f 'name=$(.TEST_POSTGRES_CONTAINER_NAME)' -f 'status=running' -q)","")
	docker run -d --name $(.TEST_POSTGRES_CONTAINER_NAME) \
	-e POSTGRES_DB=${.TEST_POSTGRES_DB_NAME} \
	-e POSTGRES_USER=${.TEST_POSTGRES_USER} \
	-e POSTGRES_PASSWORD=${.TEST_POSTGRES_PASSWORD} \
	-e POSTGRES_HOST_AUTH_METHOD=trust \
	-p 127.0.0.1:${.TEST_POSTGRES_PORT}:5432/tcp \
	postgres:12.2-alpine
	sleep 10
endif
ifeq ("$(shell docker ps -f 'name=$(.TEST_AUTH_POSTGRES_CONTAINER_NAME)' -f 'status=running' -q)","")
	docker run -d --name $(.TEST_AUTH_POSTGRES_CONTAINER_NAME) \
	-e POSTGRES_DB=${.TEST_AUTH_POSTGRES_DB_NAME} \
	-e POSTGRES_USER=${.TEST_AUTH_POSTGRES_USER} \
	-e POSTGRES_PASSWORD=${.TEST_AUTH_POSTGRES_PASSWORD} \
	-e POSTGRES_HOST_AUTH_METHOD=trust \
	-p 127.0.0.1:${.TEST_AUTH_POSTGRES_PORT}:5432/tcp \
	postgres:12.2-alpine
	sleep 10
endif
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${.TEST_POSTGRES_HOST} port=${.TEST_POSTGRES_PORT} user=${.TEST_POSTGRES_USER} dbname=${.TEST_POSTGRES_DB_NAME} password=${.TEST_POSTGRES_PASSWORD} sslmode=disable" reset || true
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${.TEST_POSTGRES_HOST} port=${.TEST_POSTGRES_PORT} user=${.TEST_POSTGRES_USER} dbname=${.TEST_POSTGRES_DB_NAME} password=${.TEST_POSTGRES_PASSWORD} sslmode=disable" up
	$(GOOSE_BIN) -dir $(AUTH_MIGRATIONS_DIR) postgres "host=${.TEST_AUTH_POSTGRES_HOST} port=${.TEST_AUTH_POSTGRES_PORT} user=${.TEST_AUTH_POSTGRES_USER} dbname=${.TEST_AUTH_POSTGRES_DB_NAME} password=${.TEST_AUTH_POSTGRES_PASSWORD} sslmode=disable" reset || true
	$(GOOSE_BIN) -dir $(AUTH_MIGRATIONS_DIR) postgres "host=${.TEST_AUTH_POSTGRES_HOST} port=${.TEST_AUTH_POSTGRES_PORT} user=${.TEST_AUTH_POSTGRES_USER} dbname=${.TEST_AUTH_POSTGRES_DB_NAME} password=${.TEST_AUTH_POSTGRES_PASSWORD} sslmode=disable" up

test-api-environment-down:
	docker stop ${.TEST_POSTGRES_CONTAINER_NAME} | true
	docker rm ${.TEST_POSTGRES_CONTAINER_NAME} | true
	docker stop ${.TEST_AUTH_POSTGRES_CONTAINER_NAME} | true
	docker rm ${.TEST_AUTH_POSTGRES_CONTAINER_NAME} | true

test-api-environment-up-force: test-api-environment-down .test-api-environment-up

include Makefile.ci.mk

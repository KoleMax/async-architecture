.CI_TEST_POSTGRES_DB_NAME=${POSTGRES_DB_NAME}
.CI_TEST_POSTGRES_HOST=${POSTGRES_HOST}
.CI_TEST_POSTGRES_PORT=${POSTGRES_PORT}
.CI_TEST_POSTGRES_USER=${POSTGRES_USER}
.CI_TEST_POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

.CI_TEST_STORAGE_ENDPOINT=${CI_STORAGE_ENDPOINT}
.CI_TEST_STORAGE_ACCESS_KEY=${CI_STORAGE_ACCESS_KEY}
.CI_TEST_STORAGE_SECRET_KEY=${CI_STORAGE_SECRET_KEY}
.CI_TEST_STORAGE_BUCKET=${CI_STORAGE_BUCKET}

ci-migrate-db: .install-goose
	@$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} dbname=${POSTGRES_DB_NAME} password=${POSTGRES_PASSWORD} sslmode=disable" up

ci-lint: .install-lint .build-SFotaUtils
	CGO_ENABLED=1 $(GOLANGCI_BIN) --timeout 240s run --config=./.golangci.yml ./...

ci-test: .install-goose .build-SFotaUtils
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${.CI_TEST_POSTGRES_HOST} port=${.CI_TEST_POSTGRES_PORT} user=${.CI_TEST_POSTGRES_USER} dbname=${.CI_TEST_POSTGRES_DB_NAME} password=${.CI_TEST_POSTGRES_PASSWORD} sslmode=disable" reset || true
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "host=${.CI_TEST_POSTGRES_HOST} port=${.CI_TEST_POSTGRES_PORT} user=${.CI_TEST_POSTGRES_USER} dbname=${.CI_TEST_POSTGRES_DB_NAME} password=${.CI_TEST_POSTGRES_PASSWORD} sslmode=disable" up
	CGO_ENABLED=1 go test -v \
	./internal/... -test-db-dsn=postgres://${.CI_TEST_POSTGRES_USER}:${.CI_TEST_POSTGRES_PASSWORD}@${.CI_TEST_POSTGRES_HOST}:${.CI_TEST_POSTGRES_PORT}/${.CI_TEST_POSTGRES_DB_NAME}?sslmode=disable \
		-test-storage-endpoint=${.CI_TEST_STORAGE_ENDPOINT} \
		-test-storage-bucket=${.CI_TEST_STORAGE_BUCKET} \
		-test-storage-access-key=${.CI_TEST_STORAGE_ACCESS_KEY} \
		-test-storage-secret-key=${.CI_TEST_STORAGE_SECRET_KEY}